package cloudkms

import (
	"encoding/base64"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	types "github.com/srinandan/cloudkms-encryption/types"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

//kmsClient contains a client connection to cloud KMS
var kmsClient *kms.KeyManagementClient

//InitKms initializes a connection to KMS
func Init() (err error) {
	kmsClient, err = kms.NewKeyManagementClient(types.Ctx)
	if err != nil {
		return err
	}
	types.Info.Println("Cloud KMS initialized successfully")

	return nil
}

//CloseKMS closes the client connection when shutting down the server
func Close() {
	if kmsClient != nil {
		_ = kmsClient.Close()
	}

	types.Info.Println("Cloud KMS closed successfully")
}

//EncryptSymmetric will encrypt the input plaintext with the specified symmetric key.
func EncryptSymmetric(name string, plaintext []byte) (string, error) {
	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:      name,
		Plaintext: plaintext,
	}

	// Call the API.
	resp, err := kmsClient.Encrypt(types.Ctx, req)
	if err != nil {
		return "", fmt.Errorf("encrypt error: %v", err)
	}

	//base64 encode the cipher
	b64CipherText := base64.StdEncoding.EncodeToString(resp.Ciphertext)

	return b64CipherText, nil
}

//DecryptSymmetric will decrypt the input ciphertext bytes using the specified symmetric key.
func DecryptSymmetric(name string, b64CipherText []byte) ([]byte, error) {
	//base64 encode the cipher
	cipherText, err := base64.StdEncoding.DecodeString(string(b64CipherText))
	if err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       name,
		Ciphertext: cipherText,
	}
	// Call the API.
	resp, err := kmsClient.Decrypt(types.Ctx, req)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %v", err)
	}

	return resp.Plaintext, nil
}
