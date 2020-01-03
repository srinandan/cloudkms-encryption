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
	cryptoKey := os.Getenv("CRYPTO_KEY")

	if (cryptoKey == "") || (keyRing == "") || (region == "") || (projectID == "") {
		return false
	}

	types.Parent = "projects/" + projectID

	types.KMSName = types.Parent + "/locations/" + region + "/keyRings/" +
		keyRing + "/cryptoKeys/" + cryptoKey

	types.Info.Printf("Initialized parameters with PROJECT_ID=%s, REGION=%s, KEY_RING=%s and CRYPTO_KEY=%s\n",
		projectID, region, keyRing, cryptoKey)

	return true
}

//Initialize logging, context, sec mgr and kms
func Initialize() {
	//init logging
	initLog()
	//init params
	if !initParams() {
		types.Error.Fatalln("PROJECT_ID, REGION, KEY_RING and CRYPTO_KEY are mandatory params")
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
