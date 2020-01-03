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

package apis

import (
	"encoding/json"

	"github.com/gorilla/mux"

	cloudkms "github.com/srinandan/cloudkms-encryption/cloudkms"
	secmgr "github.com/srinandan/cloudkms-encryption/secmgr"
	types "github.com/srinandan/cloudkms-encryption/types"

	"io/ioutil"
	"net/http"
)

var errorMessage = types.ErrorMessage{StatusCode: http.StatusInternalServerError}

func errorHandler(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)

	errorMessage.Message = err.Error()

	if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
		types.Error.Println(err)
	}
}

func responseHandler(w http.ResponseWriter, response types.Response) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		types.Error.Println(err)
	}
}

//HealthHandler handles kubernetes healthchecks
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

//EncryptionHandler handles POST /encrypt
func EncryptionHandler(w http.ResponseWriter, r *http.Request) {
	cipherResponse := types.Response{}

	//read the body
	clearText, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		errorHandler(w, err)
		return
	}

	//encrypt the payload
	b64CipherText, err := cloudkms.EncryptSymmetric(types.KMSName, clearText)

	if err != nil {
		errorHandler(w, err)
		return
	}

	cipherResponse.Payload = b64CipherText
	responseHandler(w, cipherResponse)
}

//DecryptionHandler handles POST /encrypt
func DecryptionHandler(w http.ResponseWriter, r *http.Request) {
	clearResponse := types.Response{}

	//read the body
	b64CipherText, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		errorHandler(w, err)
		return
	}

	//decrypt the payload
	clearText, err := cloudkms.DecryptSymmetric(types.KMSName, b64CipherText)

	if err != nil {
		errorHandler(w, err)
		return
	}

	clearResponse.Payload = string(clearText)
	responseHandler(w, clearResponse)
}

//RetrieveSecretHandler retrieves a secret
func RetrieveSecretHandler(w http.ResponseWriter, r *http.Request) {
	encrypted := false
	//read path variables
	vars := mux.Vars(r)
	//read query params
	queries := r.URL.Query()
	//check if the value is encrypted
	encryptedParam := queries.Get("encrypted")
	if encryptedParam == "true" {
		encrypted = true
	}

	secretName := types.Parent + "/secrets/" + vars["secretName"] +
		"/versions/" + vars["version"]

	types.Info.Println("Retrieving seret ", secretName)

	secretBytes, err := secmgr.RetrieveSecret(secretName)

	if err != nil {
		errorHandler(w, err)
		return
	}

	secretResponse := types.Response{}

	if encrypted {
		clearText, err := cloudkms.DecryptSymmetric(types.KMSName, secretBytes)
		if err != nil {
			errorHandler(w, err)
			return
		}
		secretResponse.Payload = string(clearText)
	} else {
		secretResponse.Payload = string(secretBytes)
	}

	responseHandler(w, secretResponse)
}

//CreateSecretHandler retrieves a secret
func CreateSecretHandler(w http.ResponseWriter, r *http.Request) {
	type SecretRequest struct {
		SecretId string `json:"secretId,omitempty"`
	}
	//read the body
	secretRequestBytes, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		errorHandler(w, err)
		return
	}

	secretRequest := SecretRequest{}

	err = json.Unmarshal(secretRequestBytes, &secretRequest)
	if err != nil {
		errorHandler(w, err)
		return
	}

	types.Info.Println("Creating secret ", secretRequest.SecretId)

	secretName, err := secmgr.CreateSecret(types.Parent, secretRequest.SecretId)
	if err != nil {
		errorHandler(w, err)
		return
	}
	secretResponse := types.Response{}

	secretResponse.Payload = secretName
	responseHandler(w, secretResponse)
}

//StoreSecretHandler encryptes and stores a secret
func StoreSecretHandler(w http.ResponseWriter, r *http.Request) {
	type StoreSecretRequest struct {
		SecretId  string `json:"secretId,omitempty"`
		Payload   string `json:"payload,omitempty"`
		Encrypted bool   `json:"encrypted,omitempty"`
	}

	//read the body
	storeSecretRequestBytes, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		errorHandler(w, err)
		return
	}

	storeSecretRequest := StoreSecretRequest{}

	if err = json.Unmarshal(storeSecretRequestBytes, &storeSecretRequest); err != nil {
		errorHandler(w, err)
		return
	}

	parent := types.Parent + "/secrets/" + storeSecretRequest.SecretId
	payload := storeSecretRequest.Payload

	types.Info.Printf("Store seret %s, encrypted = %t", parent, storeSecretRequest.Encrypted)

	if storeSecretRequest.Encrypted {
		//encrypt the payload
		payload, err = cloudkms.EncryptSymmetric(types.KMSName, []byte(storeSecretRequest.Payload))
		if err != nil {
			errorHandler(w, err)
			return
		}
	}

	secretVersion, err := secmgr.AddSecret(parent, payload)
	if err != nil {
		errorHandler(w, err)
		return
	}

	storeSecretResponse := types.Response{}
	storeSecretResponse.Payload = secretVersion
	responseHandler(w, storeSecretResponse)
}
