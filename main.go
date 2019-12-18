package main

import (
	"context"
	"github.com/gorilla/mux"
	apis "github.com/srinandan/cloudkms-encryption/apis"
	types "github.com/srinandan/cloudkms-encryption/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

//Init function initializes the logger objects
func Init() {
	var infoHandle = ioutil.Discard

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	if debug {
		infoHandle = os.Stdout
	}

	warningHandle := os.Stdout
	errorHandle := os.Stdout

	types.Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	types.Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	types.Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func InitParams() bool {
	projectID := os.Getenv("PROJECT_ID")
	region := os.Getenv("REGION")
	keyRing := os.Getenv("KEY_RING")
	cryptoKey := os.Getenv("CRYPTO_KEY")

	if (cryptoKey == "") || (keyRing == "") || (region == "") || (projectID == "") {
		return false
	}

	types.Name = "projects/" + projectID + "/locations/" + region + "/keyRings/" +
		keyRing + "/cryptoKeys/" + cryptoKey

	return true
}

func main() {
	var wait time.Duration
	//init logging
	Init()
	//init params
	if !InitParams() {
		types.Error.Fatalln("PROJECT_ID, REGION, KEY_RING and CRYPTO_KEY are mandatory params")
	}

	r := mux.NewRouter()
	r.HandleFunc("/encrypt", apis.EncryptionHandler).
		Methods("POST")
	r.HandleFunc("/decrypt", apis.DecryptionHandler).
		Methods("POST")

		//the following code is from gorilla mux samples
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
