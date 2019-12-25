package types

import (
	"context"
	"log"
)

//ErrorMessage hold the return value when there is an error
type ErrorMessage struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

//Response structure used by all methods
type Response struct {
	Payload string `json:"payload,omitempty"`
}

//log levels, default is error
var (
	//Info is used for debug logs
	Info *log.Logger
	//Error is used to log errors
	Error *log.Logger
)

//KMSName stores the url in the formatproject/{project-id}/secrets/{secret}
var KMSName string

//Parent stores the url in the format project/{project-id}
var Parent string

//Ctx for client connection
var Ctx context.Context
