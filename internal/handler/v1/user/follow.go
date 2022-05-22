package user

import (
	"github.com/go-eagle/eagle/internal/ecode"
	"github.com/go-eagle/eagle/internal/service"

	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/log"
)

// Follow focus on
// @Summary Follow user by user id
// @Description Get an user by user id
// @Tags user
// @Accept  json
// @Produce  json
// @Param user_id body string true "userid"
// @Success 200 {object} model.UserInfo "User Info"
// @Router /users/follow [post]
func Follow(c *gin.Context) {
	var req FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("follow bind param err: %v", err)
		response.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	// Get the user by the `user_id` from the database.
	_, err := service.Svc.Users().GetUserByID(c, req.UserID)
	if err != nil {
		response.Error(c, ecode.ErrUserNotFound.WithDetails(err.Error()))
		return
	}

	userID := service.GetUserID(c)
	// can't focus on oneself
	if userID == req.UserID {
		response.Error(c, ecode.ErrCannotFollowSelf)
		return
	}

	// Check if you have followed
	isFollowed := service.Svc.Relations().IsFollowing(c.Request.Context(), userID, req.UserID)
	if isFollowed {
		response.Error(c, errcode.Success)
		return
	}

	// Check if you have followed
	err = service.Svc.Relations().Follow(c.Request.Context(), userID, req.UserID)
	if err != nil {
		log.Warnf("[follow] add user follow err: %v", err)
		response.Error(c, errcode.ErrInternalServer.WithDetails(err.Error()))
		return
	}

	response.Success(c, nil)
}
