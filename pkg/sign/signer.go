package sign

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/go-eagle/eagle/pkg/utils"
)

// CryptoFunc Signature encryption function
type CryptoFunc func(secretKey string, args string) []byte

// Signer define
type Signer struct {
	*DefaultKeyName

	body       url.Values // Signature parameter body
	bodyPrefix string     // parameter body prefix
	bodySuffix string     // parameter body suffix
	splitChar  string     // prefix, suffix separator

	secretKey  string // signing key
	cryptoFunc CryptoFunc
}

// NewSigner Instantiate Signer
func NewSigner(cryptoFunc CryptoFunc) *Signer {
	return &Signer{
		DefaultKeyName: newDefaultKeyName(),
		body:           make(url.Values),
		bodyPrefix:     "",
		bodySuffix:     "",
		splitChar:      "",
		cryptoFunc:     cryptoFunc,
	}
}

// NewSignerMd5 md5 encryption algorithm
func NewSignerMd5() *Signer {
	return NewSigner(Md5Sign)
}

// NewSignerHmac hmac encryption algorithm
func NewSignerHmac() *Signer {
	return NewSigner(HmacSign)
}

// NewSignerAes aes symmetric encryption algorithm
func NewSignerAes() *Signer {
	return NewSigner(AesSign)
}

// SetBody Sets the entire parameter body object.
func (s *Signer) SetBody(body url.Values) {
	for k, v := range body {
		s.body[k] = v
	}
}

// GetBody Return Body content
func (s *Signer) GetBody() url.Values {
	return s.body
}

// AddBody Add signature body fields and values
func (s *Signer) AddBody(key string, value string) *Signer {
	return s.AddBodies(key, []string{value})
}

// AddBodies add value to body
func (s *Signer) AddBodies(key string, value []string) *Signer {
	s.body[key] = value
	return s
}

// SetTimeStamp set timestamp parameter
func (s *Signer) SetTimeStamp(ts int64) *Signer {
	return s.AddBody(s.Timestamp, strconv.FormatInt(ts, 10))
}

// GetTimeStamp Get TimeStamp
func (s *Signer) GetTimeStamp() string {
	return s.body.Get(s.Timestamp)
}

// SetNonceStr set random string parameter
func (s *Signer) SetNonceStr(nonce string) *Signer {
	return s.AddBody(s.NonceStr, nonce)
}

// GetNonceStr Returns the NonceStr string
func (s *Signer) GetNonceStr() string {
	return s.body.Get(s.NonceStr)
}

// SetAppID Set AppId parameter
func (s *Signer) SetAppID(appID string) *Signer {
	return s.AddBody(s.AppID, appID)
}

// GetAppID get app id
func (s *Signer) GetAppID() string {
	return s.body.Get(s.AppID)
}

// RandNonceStr Automatically generate 16-bit random string parameters
func (s *Signer) RandNonceStr() *Signer {
	return s.SetNonceStr(utils.RandomStr(16))
}

// SetSignBodyPrefix Set the prefix string for the signature string
func (s *Signer) SetSignBodyPrefix(prefix string) *Signer {
	s.bodyPrefix = prefix
	return s
}

// SetSignBodySuffix Set the suffix string of the signature string
func (s *Signer) SetSignBodySuffix(suffix string) *Signer {
	s.bodySuffix = suffix
	return s
}

// SetSplitChar Set the separator between prefix, suffix and signature body. Default is empty string
func (s *Signer) SetSplitChar(split string) *Signer {
	s.splitChar = split
	return s
}

// SetAppSecret Set the signing key
func (s *Signer) SetAppSecret(appSecret string) *Signer {
	s.secretKey = appSecret
	return s
}

// SetAppSecretWrapBody Concatenate the AppSecret string at the head and tail of the signature parameter body.
func (s *Signer) SetAppSecretWrapBody(appSecret string) *Signer {
	s.SetSignBodyPrefix(appSecret)
	s.SetSignBodySuffix(appSecret)
	return s.SetAppSecret(appSecret)
}

// GetSignBodyString Get the raw string used for signing
func (s *Signer) GetSignBodyString() string {
	return s.MakeRawBodyString()
}

// MakeRawBodyString Get the raw string used for signing
func (s *Signer) MakeRawBodyString() string {
	return s.bodyPrefix + s.splitChar + s.getSortedBodyString() + s.splitChar + s.bodySuffix
}

// GetSignedQuery Get query string with signed parameters
func (s *Signer) GetSignedQuery() string {
	return s.MakeSignedQuery()
}

// MakeSignedQuery get string with signature parameter
func (s *Signer) MakeSignedQuery() string {
	body := s.getSortedBodyString()
	sign := s.GetSignature()
	return body + "&" + s.Sign + "=" + sign
}

// GetSignature get signature
func (s *Signer) GetSignature() string {
	return s.MakeSign()
}

// MakeSign Generate signature
func (s *Signer) MakeSign() string {
	sign := fmt.Sprintf("%x", s.cryptoFunc(s.secretKey, s.GetSignBodyString()))
	return sign
}

func (s *Signer) getSortedBodyString() string {
	return SortKVPairs(s.body)
}

// SortKVPairs Concatenate the key-value pairs of Map into strings in lexicographical order
func SortKVPairs(m url.Values) string {
	size := len(m)
	if size == 0 {
		return ""
	}
	keys := make([]string, size)
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	pairs := make([]string, size)
	for i, key := range keys {
		pairs[i] = key + "=" + strings.Join(m[key], ",")
	}
	return strings.Join(pairs, "&")
}
