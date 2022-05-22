package model

import "time"

// UserStatModel User data statistics table
type UserStatModel struct {
	ID            uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	UserID        uint64    `gorm:"column:user_id;not null" json:"user_id" binding:"required"`
	FollowCount   int       `gorm:"column:follow_count" json:"follow_count"`
	FollowerCount int       `gorm:"column:follower_count" json:"follower_count"`
	Status        int       `gorm:"column:status" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"-"`
}

// TableName Table Name
func (u *UserStatModel) TableName() string {
	return "user_stat"
}
