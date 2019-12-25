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
