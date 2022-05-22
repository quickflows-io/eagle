package sign

import (
	"fmt"
	"net/url"
	"testing"
)

func TestSignMd5(t *testing.T) {
	signer := NewSignerMd5()
	signer.SetAppID("123456")
	signer.SetTimeStamp(1594458195)
	signer.SetNonceStr("supertempstr")
	signer.AddBody("city", "beijing")
	signer.SetAppSecretWrapBody("20200711")

	fmt.Println("Generate pre-signature string：" + signer.GetSignBodyString())
	fmt.Println("generate sign：" + signer.GetSignature())
	fmt.Println("output URL string：" + signer.GetSignedQuery())
	if "app_id=123456&city=beijing&nonce_str=supertempstr&timestamp=1594458195&sign=af603c2375aa2265e970737600555d7f" != signer.GetSignedQuery() {
		t.Fatal("Md5校验失败")
	}
}

func TestSigner_AddBody(t *testing.T) {
	body := make(url.Values)
	body["username"] = []string{"1024casts"}
	body["tags"] = []string{"github", "gopher"}

	signer := NewSignerHmac()
	signer.SetAppSecret("20200711")
	signer.SetTimeStamp(1594458195)
	signer.SetAppID("112233")
	signer.SetNonceStr("supertempstr")
	for k, v := range body {
		signer.AddBodies(k, v)
	}

	body.Add(KeyNameTimeStamp, "1594458195")
	body.Add(KeyNameAppID, "eagle")
	body.Add(KeyNameNonceStr, "eagle_nonce")

	fmt.Println("Generate signature string：" + signer.GetSignBodyString())
	fmt.Println("output URL string：" + signer.GetSignedQuery())

	verifier := NewVerifier()
	verifier.ParseValues(body)

	resigner := NewSignerHmac()
	resigner.SetAppSecret("eagle_key")
	resigner.SetBody(verifier.GetBodyWithoutSign())

	fmt.Println("Regenerate the signature string：" + resigner.GetSignBodyString())
	fmt.Println("Reprint URL string：" + resigner.GetSignedQuery())
}
