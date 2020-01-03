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

package secmgr

import (
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	types "github.com/srinandan/cloudkms-encryption/types"
	secretpb "google.golang.org/genproto/googleapis/cloud/secrets/v1beta1"
)

//secClient contains a client connection to Secret Manager
var secClient *secretmanager.Client

//Init initializes a connection to KMS
func Init() (err error) {
	secClient, err = secretmanager.NewClient(types.Ctx)
	if err != nil {
		return err
	}
	types.Info.Println("SecurityManager initialized successfully")
	return nil
}

//Close closes the client connection when shutting down the server
func Close() {
	if secClient != nil {
		_ = secClient.Close()
	}
	types.Info.Println("SecurityManager closed successfully")
}

//RetrieveSecret from Secret Manager
func RetrieveSecret(name string) ([]byte, error) {
	// Build the request.
	req := &secretpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	resp, err := secClient.AccessSecretVersion(types.Ctx, req)
	if err != nil {
		return nil, fmt.Errorf("access error: %v", err)
	}

	return resp.Payload.Data, nil
}

//CreateSecret version in Secret Manager
func CreateSecret(parent string, secretId string) (string, error) {
	// Build the request.
	req := &secretpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretId,
		Secret: &secretpb.Secret{
			Replication: &secretpb.Replication{
				Replication: &secretpb.Replication_Automatic_{
					Automatic: &secretpb.Replication_Automatic{},
				},
			},
		},
	}

	// Call the API.
	secResp, err := secClient.CreateSecret(types.Ctx, req)
	if err != nil {
		return "", err
	}

	return secResp.Name, nil
}

//AddSecret into Secret Manager
func AddSecret(parent string, payload string) (string, error) {
	// Build the request.
	req := &secretpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretpb.SecretPayload{
			Data: []byte(payload),
		},
	}

	// Call the API.
	secVerResp, err := secClient.AddSecretVersion(types.Ctx, req)
	if err != nil {
		return "", err
	}

	return secVerResp.Name, nil
}
