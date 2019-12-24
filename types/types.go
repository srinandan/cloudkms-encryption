package types

import (
	kms "cloud.google.com/go/kms/apiv1"
	"context"
	"log"
)

//ErrorMessage hold the return value when there is an error
type ErrorMessage struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

//EncryptResponse
type EncryptResponse struct {
	Base64EncodedCipherText string `json:"base64_cipher_text,omitempty"`
}

//DecryptResponse
type DecryptResponse struct {
	ClearText string `json:"clear_text,omitempty"`
}

//log levels, default is error
var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

//Name
var Name string

//Client connection to KMS
var Client *kms.KeyManagementClient

//Ctx for client connection
var Ctx context.Context
