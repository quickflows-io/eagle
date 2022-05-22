package user

import (
	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/internal/ecode"
	"github.com/go-eagle/eagle/internal/service"
	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/log"
)

// Register
// @Summary register
// @Description User registration
// @Tags user
// @Produce  json
// @Param req body RegisterRequest true "request parameters"
// @Success 200 {object} model.UserInfo "User Info"
// @Router /Register [post]
func Register(c *gin.Context) {
	// Binding the data with the u struct.
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("register bind param err: %v", err)
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	log.Infof("register req: %#v", req)
	// check param
	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Warnf("params is empty: %v", req)
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	// User information twice password is correct
	if req.Password != req.ConfirmPassword {
		log.Warnf("twice password is not same")
		response.Error(c, ecode.ErrTwicePasswordNotMatch)
		return
	}

	err := service.Svc.Users().Register(c, req.Username, req.Email, req.Password)
	if err != nil {
		log.Warnf("register err: %v", err)
		response.Error(c, ecode.ErrRegisterFailed.WithDetails(err.Error()))
		return
	}

	response.Success(c, nil)
}
