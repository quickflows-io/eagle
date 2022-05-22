package user

import (
	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/internal/ecode"
	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/log"
)

// Login Email Login
// @Summary User login interface
// @Description User login interface
// @Tags user
// @Produce  json
// @Param req body LoginCredentials true "request parameters"
// @Success 200 {object} model.UserInfo "User Info"
// @Router /login [post]
func Login(c *gin.Context) {
	// Binding the data with the u struct.
	var req LoginCredentials
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		log.Warnf("app.BindAndValid errs: %v", errs)
		response.Error(c, errcode.ErrInvalidParam.WithDetails(errs.Errors()...))
		return
	}

	log.Infof("login req %#v", req)

	t, err := service.Svc.Users().EmailLogin(c, req.Email, req.Password)
	if err != nil {
		log.Warnf("email login err: %v", err)
		response.Error(c, ecode.ErrEmailOrPassword)
		return
	}

	response.Success(c, model.Token{Token: t})
}

// PhoneLogin Mobile phone login interface
// @Summary User login interface
// @Description Mobile login only
// @Tags user
// @Produce  json
// @Param req body PhoneLoginCredentials true "phone"
// @Success 200 {object} model.UserInfo "User Info"
// @Router /users/login [post]
func PhoneLogin(c *gin.Context) {
	log.Info("Phone Login function called.")

	// Binding the data with the u struct.
	var req PhoneLoginCredentials
	if err := c.Bind(&req); err != nil {
		log.Warnf("phone login bind param err: %v", err)
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	log.Infof("req %#v", req)
	// check param
	if req.Phone == 0 || req.VerifyCode == 0 {
		log.Warn("phone login bind param is empty")
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	// verification code
	if !service.Svc.VCode().CheckLoginVCode(req.Phone, req.VerifyCode) {
		response.Error(c, ecode.ErrVerifyCode)
		return
	}

	// Log in
	t, err := service.Svc.Users().PhoneLogin(c, req.Phone, req.VerifyCode)
	if err != nil {
		response.Error(c, ecode.ErrVerifyCode)
		return
	}

	response.Success(c, model.Token{Token: t})
}
