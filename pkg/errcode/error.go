package errcode

import (
	"fmt"
	"net/http"
)

// Error Returns a structure of error codes and messages
// nolint: govet
type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:"details"`
}

var errorCodes = map[int]struct{}{}

// NewError create a error
func NewError(code int, msg string) *Error {
	if _, ok := errorCodes[code]; ok {
		panic(fmt.Sprintf("code %d is exsit, please change one", code))
	}
	errorCodes[code] = struct{}{}
	return &Error{code: code, msg: msg}
}

// Error return a error string
func (e Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code(), e.Msg())
}

// Code return error code
func (e *Error) Code() int {
	return e.code
}

// Msg return error msg
func (e *Error) Msg() string {
	return e.msg
}

// Msgf format error string
func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

// Details return more error details
func (e *Error) Details() []string {
	return e.details
}

// WithDetails return err with detail
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	newError.details = append(newError.details, details...)

	return &newError
}

// ToHTTPStatusCode convert custom error code to http status code and avoid return unknown status code.
func ToHTTPStatusCode(code int) int {
	switch code {
	case Success.Code():
		return http.StatusOK
	case ErrInternalServer.Code():
		return http.StatusInternalServerError
	case ErrInvalidParam.Code():
		return http.StatusBadRequest
	case ErrToken.Code():
		fallthrough
	case ErrInvalidToken.Code():
		fallthrough
	case ErrTokenTimeout.Code():
		return http.StatusUnauthorized
	case ErrTooManyRequests.Code():
		return http.StatusTooManyRequests
	case ErrServiceUnavailable.Code():
		return http.StatusServiceUnavailable
	}

	return http.StatusInternalServerError
}

// Err represents an error
type Err struct {
	Code    int
	Message string
	Err     error
}

// Error return error string
func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

// DecodeErr Decode the error, return the error code and error message
func DecodeErr(err error) (int, string) {
	if err == nil {
		return Success.code, Success.msg
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Error:
		return typed.code, typed.msg
	default:
	}

	return ErrInternalServer.Code(), err.Error()
}
