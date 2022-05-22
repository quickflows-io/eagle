This package mainly provides signature generation and signature verification for API requests.
It mainly includes the following two parts:

## Generate signature

The generated signature needs to meet the following points:

- Mutability: the signature must be different each time
- Timeliness: Timeliness of each request, invalid after expiration
- Uniqueness: Each signature is unique
- Integrity: Ability to verify incoming data to prevent tampering

If signature verification is added, several more parameters need to be passed:

- app_id represents the App Id, which is used to identify the caller
- timestamp represents the timestamp, which is used to verify the timeliness of the interface
- sign represents a signature encrypted string, which is used to verify the integrity of the data and prevent data tampering

Supports two signature generation algorithms：

- MD5: provided by NewSingerMd5()
- Sha1 + Hmac: provided by NewSingerHmac()

If the above two are not satisfied, you can customize the signature algorithm and use `NewSigner(FUNC)` Specify the implementation that implements the signature generation algorithm.




### Usage

```go
signer := NewSignerMd5()

// Set signature basic parameters
signer.SetAppId("94857djfi49484")
signer.SetTimeStamp(1594294833)
signer.SetNonceStr("xiKdApRhbuxVckJa")

// Set other parameters involved in signing
signer.AddBody("plate_number", "golang")

// AppSecretKey，Before and after wrapping the signature body string
signer.SetAppSecretWrapBody("x90449dfde34d")

fmt.Println("Generate signature string：" + signer.GetUnsignedString())
fmt.Println("output URL string：" + signer.GetSignedQuery())
```

### Result

## check signature

The sign.Verifier tool class is used to verify the format and timestamp of the signature parameters. It is used together with Signer to verify the signature information of API requests on the server side.

### Usage

```go
requestUri := "/restful/api/numbers?app_id=9d8a121ce581499d&nonce_str=ibuaiVcKdpRxkhJA&plate_number=豫A66666" +
		"&time_stamp=1532585241&sign=072defd1a251dc58e4d1799e17ffe7a4"

	// Step 1: Create a Verifier verification class
	verifier := NewVerifier()

	// Assume that the verification parameters are read from the RequestUri
	if err := verifier.ParseQuery(requestUri); nil != err {
		t.Fatal(err)
	}

	// Or use verifier.ParseValues(Values) to parse.

	// Step 2: (Optional) Check whether the necessary parameters for signature verification are included
	if err := verifier.MustHasOtherFields("plate_number"); nil != err {
		t.Fatal(err)
	}

	// Step 3: Check if the timestamp has timed out.

	// Timestamp timeout: 5 minutes
	verifier.SetTimeout(time.Minute * 5)
	if err := verifier.CheckTimeStamp(); nil != err {
		t.Fatal(err)
	}

	// Step 4: Create a Signer to reproduce the client's signature information
	signer := NewSignerMd5()

	// Step 5: Read all request parameters from Verifier
	signer.SetBody(verifier.GetBodyWithoutSign())

	// Step 6: Read the SecretKey corresponding to the AppID from the database
	// appId := verifier.GetAppId()
	secretKey := "123abc456"

	// Use the same WrapBody method
	signer.SetAppSecretWrapBody(secretKey)

	// The server generates a signature based on the client parameters
	sign := signer.GetSignature()

    //Finally, compare the signature information generated by the server to see if it is consistent with the signature provided by the client.
	if verifier.MustString("sign") != sign {
		t.Fatal("Verification failed")
	}
```