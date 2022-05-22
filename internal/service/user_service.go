package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-eagle/eagle/internal/repository"

	"github.com/pkg/errors"

	pb "github.com/go-eagle/eagle/api/grpc/user/v1"
	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/auth"
	"github.com/go-eagle/eagle/pkg/log"
)

// UserService define interface func
type UserService interface {
	Register(ctx context.Context, username, email, password string) error
	EmailLogin(ctx context.Context, email, password string) (tokenStr string, err error)
	PhoneLogin(ctx context.Context, phone int64, verifyCode int) (tokenStr string, err error)
	// LoginByPhone(ctx context.Context, req *pb.PhoneLoginRequest) (reply *pb.PhoneLoginReply, err error)
	GetUserByID(ctx context.Context, id uint64) (*model.UserBaseModel, error)
	GetUserInfoByID(ctx context.Context, id uint64) (*model.UserInfo, error)
	GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error)
	GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error)
	UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error
	BatchGetUsers(ctx context.Context, userID uint64, userIDs []uint64) ([]*model.UserInfo, error)
}

type userService struct {
	repo repository.Repository
}

var _ UserService = (*userService)(nil)

func newUsers(svc *service) *userService {
	return &userService{repo: svc.repo}
}

// Register registered user
func (s *userService) Register(ctx context.Context, username, email, password string) error {
	pwd, err := auth.HashAndSalt(password)
	if err != nil {
		return errors.Wrapf(err, "encrypt password err")
	}

	u := model.UserBaseModel{
		Username:  username,
		Password:  pwd,
		Email:     email,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	isExist, err := s.repo.UserIsExist(&u)
	if err != nil {
		return errors.Wrapf(err, "create user")
	}
	if isExist {
		return errors.New("user already exists")
	}
	_, err = s.repo.CreateUser(ctx, &u)
	if err != nil {
		return errors.Wrapf(err, "create user")
	}
	return nil
}

// EmailLogin Email Login
func (s *userService) EmailLogin(ctx context.Context, email, password string) (tokenStr string, err error) {
	u, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.Wrapf(err, "get user info err by email")
	}

	// ComparePasswords the login password with the user password.
	if !auth.ComparePasswords(u.Password, password) {
		return "", errors.New("invalid password")
	}

	// issue signature Sign the json web token.
	payload := map[string]interface{}{"user_id": u.ID, "username": u.Username}
	tokenStr, err = app.Sign(ctx, payload, app.Conf.JwtSecret, 86400)
	if err != nil {
		return "", errors.Wrapf(err, "gen token sign err")
	}

	return tokenStr, nil
}

// LoginByPhone phone login, grpc wrapper
func (s *userService) LoginByPhone(ctx context.Context, req *pb.PhoneLoginRequest) (reply *pb.PhoneLoginReply, err error) {
	tokenStr, err := s.PhoneLogin(ctx, req.Phone, int(req.VerifyCode))
	if err != nil {
		log.Warnf("[service.user] phone login err: %v, params: %v", err, req)
	}
	reply = &pb.PhoneLoginReply{
		Ret: tokenStr,
		Err: "",
	}
	return
}

// PhoneLogin Email Login
func (s *userService) PhoneLogin(ctx context.Context, phone int64, verifyCode int) (tokenStr string, err error) {
	// If it is a registered user, obtain user information through mobile phone number
	u, err := s.GetUserByPhone(ctx, phone)
	if err != nil {
		return "", errors.Wrapf(err, "[login] get u info err")
	}

	// Otherwise, create a new user information and get the user information
	if u.ID == 0 {
		u := model.UserBaseModel{
			Phone:    phone,
			Username: strconv.Itoa(int(phone)),
		}
		u.ID, err = s.repo.CreateUser(ctx, &u)
		if err != nil {
			return "", errors.Wrapf(err, "[login] create user err")
		}
	}

	// Sign the signature Sign the json web token.
	payload := map[string]interface{}{"user_id": u.ID, "username": u.Username}
	tokenStr, err = app.Sign(ctx, payload, app.Conf.JwtSecret, 86400)
	if err != nil {
		return "", errors.Wrapf(err, "[login] gen token sign err")
	}
	return tokenStr, nil
}

// UpdateUser update user info
func (s *userService) UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error {
	err := s.repo.UpdateUser(ctx, id, userMap)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByID Get a single user information
func (s *userService) GetUserByID(ctx context.Context, id uint64) (*model.UserBaseModel, error) {
	return s.repo.GetUser(ctx, id)
}

// GetUserStatByID  Get a single user information
func (s *userService) GetUserStatByID(ctx context.Context, id uint64) (*model.UserStatModel, error) {
	return s.repo.GetUserStatByID(ctx, id)
}

// GetUserInfoByID Get assembled user data
func (s *userService) GetUserInfoByID(ctx context.Context, id uint64) (*model.UserInfo, error) {
	userInfos, err := s.BatchGetUsers(ctx, id, []uint64{id})
	if err != nil {
		return nil, err
	}
	return userInfos[0], nil
}

// BatchGetUsers Get user information in batches
// 1. Handling following and being followed states
// 2. Get followers and followers data
func (s *userService) BatchGetUsers(ctx context.Context, userID uint64, userIDs []uint64) ([]*model.UserInfo, error) {
	infos := make([]*model.UserInfo, 0)
	// Get current user information
	curUser, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "[user_service] get one user err")
	}

	// Get user information in batches
	users, err := s.repo.GetUsersByIds(ctx, userIDs)
	if err != nil {
		return nil, errors.Wrap(err, "[user_service] batch get user err")
	}

	wg := sync.WaitGroup{}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// Get your own follow status on a watchlist
	userFollowMap, err := s.repo.GetFollowByUIds(ctx, userID, userIDs)
	if err != nil {
		errChan <- err
	}

	// Get your follow status on the watchlist
	userFansMap, err := s.repo.GetFansByUIds(ctx, userID, userIDs)
	if err != nil {
		errChan <- err
	}

	// Get user statistics
	userStatMap, err := s.repo.GetUserStatByIDs(ctx, userIDs)
	if err != nil {
		errChan <- err
	}

	var m sync.Map

	// parallel processing
	for _, u := range users {
		wg.Add(1)
		go func(u *model.UserBaseModel) {
			defer wg.Done()

			isFollow := 0
			_, ok := userFollowMap[u.ID]
			if ok {
				isFollow = 1
			}

			isFollowed := 0
			_, ok = userFansMap[u.ID]
			if ok {
				isFollowed = 1
			}

			userStatMap, ok := userStatMap[u.ID]
			if !ok {
				userStatMap = nil
			}

			transInput := &TransferUserInput{
				CurUser:  curUser,
				User:     u,
				UserStat: userStatMap,
				IsFollow: isFollow,
				IsFans:   isFollowed,
			}
			userInfo := TransferUser(transInput)
			if err != nil {
				errChan <- err
				return
			}
			m.Store(u.ID, &userInfo)
		}(u)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		log.Warnf("[user_service] batch get user err chan: %v", err)
		return nil, err
	case <-time.After(3 * time.Second):
		return nil, fmt.Errorf("list users timeout after 3 seconds")
	}

	// Guaranteed order
	for _, u := range users {
		info, _ := m.Load(u.ID)
		infos = append(infos, info.(*model.UserInfo))
	}

	return infos, nil
}

func (s *userService) GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error) {
	userModel, err := s.repo.GetUserByPhone(ctx, phone)
	if err != nil {
		return userModel, errors.Wrapf(err, "get user info err from db by phone: %d", phone)
	}

	return userModel, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error) {
	userModel, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return userModel, errors.Wrapf(err, "get user info err from db by email: %s", email)
	}

	return userModel, nil
}
