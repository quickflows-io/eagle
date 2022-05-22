package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/internal/ecode"
	"github.com/go-eagle/eagle/internal/service"
	"github.com/go-eagle/eagle/pkg/log"
)

// VCode get verification code
// @Summary Get verification code based on mobile phone number
// @Description Get an user by username
// @Tags user
// @Accept  json
// @Produce  json
// @Param area_code query string true "area code，such as 86"
// @Param phone query string true "phone number"
// @Success 200 {object} app.Response
// @Router /vcode [get]
func VCode(c *gin.Context) {
	// Verify that the area code and mobile number are empty
	if c.Query("area_code") == "" {
		log.Warn("vcode area code is empty")
		response.Error(c, ecode.ErrAreaCodeEmpty)
		return
	}

	phone := c.Query("phone")
	if phone == "" {
		log.Warn("vcode phone is empty")
		response.Error(c, ecode.ErrPhoneEmpty)
		return
	}

	// TODO: Frequency control to prevent attacks

	// Generate SMS verification code
	verifyCode, err := service.Svc.VCode().GenLoginVCode(phone)
	if err != nil {
		log.Warnf("gen login verify code err, %v", errors.WithStack(err))
		response.Error(c, ecode.ErrGenVCode)
		return
	}

	// send messages
	err = service.Svc.SMS().SendSMS(phone, verifyCode)
	if err != nil {
		log.Warnf("send phone sms err, %v", errors.WithStack(err))
		response.Error(c, ecode.ErrSendSMS)
		return
	}

	response.Success(c, nil)
}
