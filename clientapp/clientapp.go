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

package clientapp

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	cloudkms "github.com/srinandan/cloudkms-encryption/cloudkms"
	secmgr "github.com/srinandan/cloudkms-encryption/secmgr"
	types "github.com/srinandan/cloudkms-encryption/types"
)

//initLog function initializes the logger objects
func initLog() {
	var infoHandle = ioutil.Discard

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	if debug {
		infoHandle = os.Stdout
	}

	errorHandle := os.Stdout

	types.Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	types.Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

//initParams initializes parameters for cloud kms
func initParams() bool {
	projectID := os.Getenv("PROJECT_ID")

	region := os.Getenv("REGION")
	keyRing := os.Getenv("KEY_RING")
	symCryptoKey := os.Getenv("SYM_CRYPTO_KEY")
	asymCryptoKey := os.Getenv("ASYM_CRYPTO_KEY")

	if (region == "") || (keyRing == "") || (symCryptoKey == "") || (projectID == "") || (asymCryptoKey == "") {
		return false
	}

	types.Parent = "projects/" + projectID

	types.SymmetricKMSName = types.Parent + "/locations/" + region + "/keyRings/" +
		keyRing + "/cryptoKeys/" + symCryptoKey

//Resource name 'projects/nandanks-151422/locations/us-west1/keyRings/test/cryptoKeys/asymmetric-key' 
//does not match pattern 'projects/([^/]+)/locations/([a-zA-Z0-9_-]{1,63})/keyRings/([a-zA-Z0-9_-]{1,63})/cryptoKeys/([a-zA-Z0-9_-]{1,63})/cryptoKeyVersions/([a-zA-Z0-9_-]{1,63})'."}
	types.AsymmetricKMSName = types.Parent + "/locations/" + region + "/keyRings/" +
		keyRing + "/cryptoKeys/" + asymCryptoKey + "/cryptoKeyVersions/" + "1"

	types.Info.Printf("Initialized parameters with PROJECT_ID=%s, REGION=%s, KEY_RING=%s, SYM_CRYPTO_KEY=%s and ASYM_CRYPTO_KEY=%s\n",
		projectID, region, keyRing, symCryptoKey, asymCryptoKey)

	return true
}

//Initialize logging, context, sec mgr and kms
func Initialize() {
	//init logging
	initLog()
	//init params
	if !initParams() {
		types.Error.Fatalln("PROJECT_ID, REGION, KEY_RING, SYM_CRYPTO_KEY and ASYM_CRYPTO_KEY are mandatory params")
	}
	//init ctx
	types.Ctx = context.Background()
	//init cloud kms
	if err := cloudkms.Init(); err != nil {
		types.Error.Fatalln("error connecting to KMS ", err)
	}
	//init sec manager
	if err := secmgr.Init(); err != nil {
		types.Error.Fatalln("error connecting to Secret Manager ", err)
	}
}

//Close client connections
func Close() {
	cloudkms.Close()
	secmgr.Close()
}
