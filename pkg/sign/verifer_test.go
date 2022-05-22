package sign

import (
	"fmt"
	"testing"
)

func TestVerify_ParseQuery(t *testing.T) {
	requestURI := "/restful/api/numbers?app_id=9d8a121ce581499d&nonce_str=tempstring&city=beijing" +
		"&timestamp=1532585241&sign=0f5b8c97920bc95f1a8b893f41b42d9e"

	// Step 1: Create a Verify verification class
	verifier := NewVerifier()

	// Assume that the verification parameters are read from the RequestUri
	if err := verifier.ParseQuery(requestURI); nil != err {
		t.Fatal(err)
	}

	// Step 2: (Optional) Check whether the necessary parameters for signature verification are included
	if err := verifier.MustHasOtherKeys("city"); nil != err {
		t.Fatal(err)
	}

	// Step 3: Check if the timestamp has timed out.
	//if err := verifier.CheckTimeStamp(); nil != err {
	//	t.Fatal(err)
	//}

	// The fourth step, create a Sign to reproduce the client's signature information:
	signer := NewSignerMd5()

	// Step 5: Read all request parameters from Verify
	signer.SetBody(verifier.GetBodyWithoutSign())

	// Step 6: Read the SecretKey corresponding to the AppID from the database
	// appId := verifier.MustString("app_id")
	secretKey := "d93047a4d6fe6111"

	// Use the same WrapBody method
	signer.SetAppSecretWrapBody(secretKey)

	// generate
	sign := signer.GetSignature()
	t.Log("sign", sign)

	// Verify that the generated one is the same as the one passed over
	if verifier.MustString("sign") != sign {
		t.Fatal("Verification failed")
	}

	fmt.Println(sign)
}
