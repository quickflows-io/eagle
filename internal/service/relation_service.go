package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/internal/repository"
	"github.com/go-eagle/eagle/pkg/log"
)

const (
	// FollowStatusNormal Attention Status - Normal
	FollowStatusNormal int = 1 //normal
	// FollowStatusDelete Follow Status - Delete
	FollowStatusDelete = 0 // delete
)

// RelationService .
type RelationService interface {
	Follow(ctx context.Context, userID uint64, followedUID uint64) error
	Unfollow(ctx context.Context, userID uint64, followedUID uint64) error
	IsFollowing(ctx context.Context, userID uint64, followedUID uint64) bool
	GetFollowingUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFollowModel, error)
	GetFollowerUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFansModel, error)
}

type relationService struct {
	repo repository.Repository
}

var _ RelationService = (*relationService)(nil)

func newRelations(svc *service) *relationService {
	return &relationService{repo: svc.repo}
}

// IsFollowing Are you following a user
func (s *relationService) IsFollowing(ctx context.Context, userID uint64, followedUID uint64) bool {
	userFollowModel := &model.UserFollowModel{}
	result := model.GetDB().
		Where("user_id=? AND followed_uid=? ", userID, followedUID).
		Find(userFollowModel)

	if err := result.Error; err != nil {
		log.Warnf("[user_service] get user follow err, %v", err)
		return false
	}

	if userFollowModel.ID > 0 && userFollowModel.Status == FollowStatusNormal {
		return true
	}

	return false
}

// Follow Follow target users
func (s *relationService) Follow(ctx context.Context, userID uint64, followedUID uint64) error {
	db := model.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// add to watchlist
	err := s.repo.CreateUserFollow(ctx, tx, userID, followedUID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert into user follow err")
	}

	// add to fan list
	err = s.repo.CreateUserFans(ctx, tx, followedUID, userID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert into user fans err")
	}

	// Add followers
	err = s.repo.IncrFollowCount(ctx, tx, userID, 1)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow count err")
	}

	// Add followers
	err = s.repo.IncrFollowerCount(ctx, tx, followedUID, 1)
	if err != nil {
		return errors.Wrap(err, "update user fans count err")
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "tx commit err")
	}

	return nil
}

// Unfollow unfollow user
func (s *relationService) Unfollow(ctx context.Context, userID uint64, followedUID uint64) error {
	db := model.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// remove follow
	err := s.repo.UpdateUserFollowStatus(ctx, tx, userID, followedUID, FollowStatusDelete)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow err")
	}

	// delete followers
	err = s.repo.UpdateUserFansStatus(ctx, tx, followedUID, userID, FollowStatusDelete)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow err")
	}

	// reduce the number of followers
	err = s.repo.IncrFollowCount(ctx, tx, userID, -1)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user follow count err")
	}

	// reduce the number of followers
	err = s.repo.IncrFollowerCount(ctx, tx, followedUID, -1)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update user fans count err")
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "tx commit err")
	}

	return nil
}

// GetFollowingUserList Get a list of users you are following
func (s *relationService) GetFollowingUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFollowModel, error) {
	if lastID == 0 {
		lastID = MaxID
	}
	userFollowList, err := s.repo.GetFollowingUserList(ctx, userID, lastID, limit)
	if err != nil {
		return nil, err
	}

	return userFollowList, nil
}

// GetFollowerUserList Get a list of fans users
func (s *relationService) GetFollowerUserList(ctx context.Context, userID uint64, lastID uint64, limit int) ([]*model.UserFansModel, error) {
	if lastID == 0 {
		lastID = MaxID
	}
	userFollowerList, err := s.repo.GetFollowerUserList(ctx, userID, lastID, limit)
	if err != nil {
		return nil, err
	}

	return userFollowerList, nil
}
