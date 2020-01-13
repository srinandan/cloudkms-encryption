// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudkms

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"	
	"encoding/base64"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	types "github.com/srinandan/cloudkms-encryption/types"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

//kmsClient contains a client connection to cloud KMS
var kmsClient *kms.KeyManagementClient

//publicKey
var publicKey *kmspb.PublicKey

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

//EncryptRSA will encrypt using a public key
func EncryptRSA(name string, plaintext []byte) (b64CipherText string, err error) {
	// name: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
	// plaintext := []byte("Sample message")

	if publicKey == nil {
		// Retrieve the public key from KMS.
		publicKey, err = kmsClient.GetPublicKey(types.Ctx, &kmspb.GetPublicKeyRequest{Name: name})
		if err != nil {
				return "", fmt.Errorf("GetPublicKey: %v", err)
		}
	}

	// Parse the key.
	block, _ := pem.Decode([]byte(publicKey.Pem))
	abstractKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
			return "", fmt.Errorf("x509.ParsePKIXPublicKey: %+v", err)
	}

	rsaKey, ok := abstractKey.(*rsa.PublicKey)
	if !ok {
			return "", fmt.Errorf("key %q is not RSA", name)
	}
	
	// Encrypt data using the RSA public key.
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, plaintext, nil)
	if err != nil {
			return "", fmt.Errorf("rsa.EncryptOAEP: %v", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

//DecryptRSA will decrypt using a private key
func DecryptRSA(name string, b64CipherText []byte) ([]byte, error) {
	//base64 encode the cipher
	cipherText, err := base64.StdEncoding.DecodeString(string(b64CipherText))
	if err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	// Build the request.
	req := &kmspb.AsymmetricDecryptRequest{
		Name:       name,
		Ciphertext: cipherText,
	}
	
	// Call the API.
	resp, err := kmsClient.AsymmetricDecrypt(types.Ctx, req)
	if err != nil {
			return nil, fmt.Errorf("asymmetricDecrypt: %v", err)
	}
	
	return resp.Plaintext, nil
}