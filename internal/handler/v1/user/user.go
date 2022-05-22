package user

import (
	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/app"
)

var response = app.NewResponse()

// CreateRequest Create user request
type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateResponse Create User Response
type CreateResponse struct {
	Username string `json:"username"`
}

// RegisterRequest register
type RegisterRequest struct {
	Username        string `json:"username" form:"username"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}

// LoginCredentials Default login method - email
type LoginCredentials struct {
	Email    string `json:"email" form:"email" binding:"required" `
	Password string `json:"password" form:"password" binding:"required" `
}

// PhoneLoginCredentials login by phone
type PhoneLoginCredentials struct {
	Phone      int64 `json:"phone" form:"phone" binding:"required" example:"13010002000"`
	VerifyCode int   `json:"verify_code" form:"verify_code" binding:"required" example:"120110"`
}

// UpdateRequest update request
type UpdateRequest struct {
	Avatar string `json:"avatar"`
	Sex    int    `json:"sex"`
}

// FollowRequest follow request
type FollowRequest struct {
	UserID uint64 `json:"user_id"`
}

// ListResponse generic list resp
type ListResponse struct {
	TotalCount uint64      `json:"total_count"`
	HasMore    int         `json:"has_more"`
	PageKey    string      `json:"page_key"`
	PageValue  int         `json:"page_value"`
	Items      interface{} `json:"items"`
}

// SwaggerListResponse Documentation
type SwaggerListResponse struct {
	TotalCount uint64           `json:"totalCount"`
	UserList   []model.UserInfo `json:"userList"`
}
