package middleware

import (
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/app"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/sign"
)

// SignMd5Middleware md5 signature verification middleware
func SignMd5Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sn, err := verifySign(c)
		response := app.NewResponse()
		if err != nil {
			response.Error(c, errcode.ErrInternalServer)
			c.Abort()
			return
		}

		if sn != nil {
			response.Error(c, errcode.ErrSignParam)
			c.Abort()
			return
		}

		c.Next()
	}
}

// verifySign Verify signature
func verifySign(c *gin.Context) (map[string]string, error) {
	requestURI := c.Request.RequestURI
	// Create a Verify validator
	verifier := sign.NewVerifier()
	sn := verifier.GetSign()

	// Assume that the verification parameters are read from the RequestUri
	if err := verifier.ParseQuery(requestURI); nil != err {
		return nil, err
	}

	// Check if the timestamp has timed out.
	if err := verifier.CheckTimeStamp(); nil != err {
		return nil, errors.Errorf("%s error", sign.KeyNameTimeStamp)
	}

	// Verify signature
	localSign := genSign()
	if sn == "" || sn != localSign {
		return nil, errors.New(fmt.Sprintf("%s error", sign.KeyNameSign))
	}

	return nil, nil
}

// genSign Generate signature
func genSign() string {
	// todo: 读取配置
	signer := sign.NewSignerMd5()
	signer.SetAppID("123456")
	signer.SetTimeStamp(time.Now().Unix())
	signer.SetNonceStr("supertempstr")
	signer.SetAppSecretWrapBody("20200711")

	return signer.GetSignature()
}
