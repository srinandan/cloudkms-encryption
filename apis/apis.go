package apis

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	cloudkms "github.com/srinandan/cloudkms-encryption/cloudkms"
	types "github.com/srinandan/cloudkms-encryption/types"
)

//EncryptionHandler handles POST /encrypt
func EncryptionHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := types.ErrorMessage{}
	errorMessage.StatusCode = http.StatusInternalServerError

	cipherResponse := types.EncryptResponse{}

	//read the body
	clearText, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		errorMessage.Message = err.Error()

		if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
			types.Error.Println(err)
		}
		return
	}

	//encrypt the payload
	b64CipherText, err := cloudkms.EncryptSymmetric(types.Name, clearText)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		errorMessage.Message = err.Error()
		if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
			types.Error.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	cipherResponse.Base64EncodedCipherText = b64CipherText

	if err := json.NewEncoder(w).Encode(cipherResponse); err != nil {
		types.Error.Println(err)
	}
}

//DecryptionHandler handles POST /encrypt
func DecryptionHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := types.ErrorMessage{}
	errorMessage.StatusCode = http.StatusInternalServerError

	clearResponse := types.DecryptResponse{}

	//read the body
	b64CipherText, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		errorMessage.Message = err.Error()
		if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
			types.Error.Println(err)
		}
		return
	}

	//decrypt the payload
	clearText, err := cloudkms.DecryptSymmetric(types.Name, b64CipherText)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		errorMessage.Message = err.Error()
		if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
			types.Error.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	clearResponse.ClearText = string(clearText)

	if err := json.NewEncoder(w).Encode(clearResponse); err != nil {
		types.Error.Println(err)
	}
}
