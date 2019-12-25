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
