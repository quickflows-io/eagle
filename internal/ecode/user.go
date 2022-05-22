package ecode

import "github.com/go-eagle/eagle/pkg/errcode"

//nolint: golint
var (
	// user errors
	ErrUserNotFound          = errcode.NewError(20101, "The user was not found.")
	ErrPasswordIncorrect     = errcode.NewError(20102, "Incorrect username or password")
	ErrAreaCodeEmpty         = errcode.NewError(20103, "Mobile phone area code cannot be empty")
	ErrPhoneEmpty            = errcode.NewError(20104, "Mobile number cannot be empty")
	ErrGenVCode              = errcode.NewError(20105, "Generate verification code error")
	ErrSendSMS               = errcode.NewError(20106, "Error sending SMS")
	ErrSendSMSTooMany        = errcode.NewError(20107, "Exceeded today's limit, please try again tomorrow")
	ErrVerifyCode            = errcode.NewError(20108, "Verification code error")
	ErrEmailOrPassword       = errcode.NewError(20109, "Mail or password is incorrect")
	ErrTwicePasswordNotMatch = errcode.NewError(20110, "Inconsistent password entered twice")
	ErrRegisterFailed        = errcode.NewError(20111, "registration failed")

	ErrCannotFollowSelf = errcode.NewError(20201, "can't focus on oneself")
)
