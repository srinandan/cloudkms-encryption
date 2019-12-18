package cloudkms

import (
	"context"
	"encoding/base64"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

//EncryptSymmetric will encrypt the input plaintext with the specified symmetric key.
func EncryptSymmetric(name string, plaintext []byte) (string, error) {
	ctx := context.Background()

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return "", fmt.Errorf("kms.NewKeyManagementClient: %v", err)
	}

	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:      name,
		Plaintext: plaintext,
	}
	// Call the API.
	resp, err := client.Encrypt(ctx, req)
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

	ctx := context.Background()
	
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("kms.NewKeyManagementClient: %v", err)
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       name,
		Ciphertext: cipherText,
	}
	// Call the API.
	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %v", err)
	}
	return resp.Plaintext, nil
}
