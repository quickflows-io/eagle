package sign

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-eagle/eagle/pkg/utils"
)

// Verifier define struct
type Verifier struct {
	*DefaultKeyName
	body url.Values

	timeout time.Duration // Signature expiration time
}

// NewVerifier Instantiate Verifier
func NewVerifier() *Verifier {
	return &Verifier{
		DefaultKeyName: newDefaultKeyName(),
		body:           make(url.Values),
		timeout:        time.Minute * 5,
	}
}

// ParseQuery Parse argument string into argument list
func (v *Verifier) ParseQuery(requestURI string) error {
	requestQuery := ""
	idx := strings.Index(requestURI, "?")
	if idx > 0 {
		requestQuery = requestURI[idx+1:]
	}
	query, err := url.ParseQuery(requestQuery)
	if nil != err {
		return err
	}
	v.ParseValues(query)
	return nil
}

// ParseValues Parse the values' parameter list into a parameter Map. If the parameters are multi-valued, join them into a string with commas.
func (v *Verifier) ParseValues(values url.Values) {
	for key, value := range values {
		v.body[key] = value
	}
}

// SetTimeout Set the signature verification expiration time
func (v *Verifier) SetTimeout(timeout time.Duration) *Verifier {
	v.timeout = timeout
	return v
}

// MustString get string value
func (v *Verifier) MustString(key string) string {
	ss := v.MustStrings(key)
	if len(ss) == 0 {
		return ""
	}
	return ss[0]
}

// MustStrings get array of string values
func (v *Verifier) MustStrings(key string) []string {
	return v.body[key]
}

// MustInt64 Get Int64 value
func (v *Verifier) MustInt64(key string) int64 {
	n, _ := utils.StringToInt64(v.MustString(key))
	return n
}

// MustHasKeys Get Int64 value
func (v *Verifier) MustHasKeys(keys ...string) error {
	for _, key := range keys {
		if _, hit := v.body[key]; !hit {
			return fmt.Errorf("KEY_MISSED:<%s>", key)
		}
	}
	return nil
}

// MustHasOtherKeys Must contain the specified field parameters other than the specific [timestamp, nonce_str, sign, app_id] etc.
func (v *Verifier) MustHasOtherKeys(keys ...string) error {
	fields := []string{v.Timestamp, v.NonceStr, v.Sign, v.AppID}
	if len(keys) > 0 {
		fields = append(fields, keys...)
	}
	return v.MustHasKeys(fields...)
}

// CheckTimeStamp Check timestamp validity
func (v *Verifier) CheckTimeStamp() error {
	timestamp := v.GetTimestamp()
	thatTime := time.Unix(timestamp, 0)
	if timestamp > time.Now().Unix() || time.Since(thatTime) > v.timeout {
		return fmt.Errorf("TIMESTAMP_TIMEOUT:<%d>", timestamp)
	}
	return nil
}

// GetAppID get app id
func (v *Verifier) GetAppID() string {
	return v.MustString(v.AppID)
}

// GetNonceStr get random string
func (v *Verifier) GetNonceStr() string {
	return v.MustString(v.NonceStr)
}

// GetSign get signature
func (v *Verifier) GetSign() string {
	return v.MustString(v.Sign)
}

// GetTimestamp get timestamp
func (v *Verifier) GetTimestamp() int64 {
	return v.MustInt64(v.Timestamp)
}

// GetBodyWithoutSign 获取所有参数体。其中不包含sign 字段
func (v *Verifier) GetBodyWithoutSign() url.Values {
	out := make(url.Values)
	for k, val := range v.body {
		if k != v.Sign {
			out[k] = val
		}
	}
	return out
}

// GetBody get body
func (v *Verifier) GetBody() url.Values {
	out := make(url.Values)
	for k, val := range v.body {
		out[k] = val
	}
	return out
}
