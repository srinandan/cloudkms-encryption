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
