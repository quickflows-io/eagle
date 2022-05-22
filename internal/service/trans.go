package service

import (
	"github.com/go-eagle/eagle/internal/model"
)

// TransferUserInput Convert input fields
type TransferUserInput struct {
	CurUser  *model.UserBaseModel
	User     *model.UserBaseModel
	UserStat *model.UserStatModel
	IsFollow int `json:"is_follow"`
	IsFans   int `json:"is_fans"`
}

// TransferUser Assemble data and output
// The user structure exposed to the outside world should be converted through this structure
func TransferUser(input *TransferUserInput) *model.UserInfo {
	if input.User == nil {
		return &model.UserInfo{}
	}

	return &model.UserInfo{
		ID:         input.User.ID,
		Username:   input.User.Username,
		Avatar:     input.User.Avatar, // todo: 转为url
		Sex:        input.User.Sex,
		UserFollow: transferUserFollow(input),
	}
}

// transferUserFollow Convert users to follow related fields
func transferUserFollow(input *TransferUserInput) *model.UserFollow {
	followCount := 0
	if input.UserStat != nil {
		followCount = input.UserStat.FollowCount
	}
	followerCount := 0
	if input.UserStat != nil {
		followerCount = input.UserStat.FollowerCount
	}

	return &model.UserFollow{
		FollowNum: followCount,
		FansNum:   followerCount,
		IsFollow:  input.IsFollow,
		IsFans:    input.IsFans,
	}
}
