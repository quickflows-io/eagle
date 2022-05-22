package sign

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

// RsaSign Asymmetric encryption
func RsaSign(secretKey, body string) []byte {
	ret, _ := PublicEncrypt(body, secretKey)
	return []byte(ret)
}

// PublicEncrypt public key encryption
func PublicEncrypt(encryptStr string, path string) (string, error) {
	// open a file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	// read file content
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, _ = file.Read(buf)

	// pem decoding
	block, _ := pem.Decode(buf)

	// x509 decoding
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// type assertion
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	//Encrypt plaintext
	encryptedStr, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(encryptStr))
	if err != nil {
		return "", err
	}

	//return ciphertext
	return base64.URLEncoding.EncodeToString(encryptedStr), nil
}

// PrivateDecrypt private key decryption
func PrivateDecrypt(decryptStr string, path string) (string, error) {
	// open a file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	// get file content
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, _ = file.Read(buf)

	// pem decoding
	block, _ := pem.Decode(buf)

	// X509 decoding
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	decryptBytes, err := base64.URLEncoding.DecodeString(decryptStr)

	//decrypt the ciphertext
	decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decryptBytes)

	//return plaintext
	return string(decrypted), nil
}
