package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	apis "github.com/srinandan/cloudkms-encryption/apis"
	clientapp "github.com/srinandan/cloudkms-encryption/clientapp"
	types "github.com/srinandan/cloudkms-encryption/types"
)

func main() {
	var wait time.Duration

	const address = "0.0.0.0:8080"

	//initialize
	clientapp.Initialize()

	r := mux.NewRouter()
	r.HandleFunc("/healthz", apis.HealthHandler).
		Methods("GET")
	r.HandleFunc("/encrypt", apis.EncryptionHandler).
		Methods("POST")
	r.HandleFunc("/decrypt", apis.DecryptionHandler).
		Methods("POST")
	//registering this handler twice since the query param is optional
	r.HandleFunc("/secrets/{secretName}/{version}", apis.RetrieveSecretHandler).
		Methods("GET").
		Queries("encrypted", "{encrypted}")
	r.HandleFunc("/secrets/{secretName}/{version}", apis.RetrieveSecretHandler).
		Methods("GET")

	r.HandleFunc("/secrets", apis.CreateSecretHandler).
		Methods("POST")
	r.HandleFunc("/storesecrets", apis.StoreSecretHandler).
		Methods("POST")

	types.Info.Println("Starting server - ", address)

	//the following code is from gorilla mux samples
	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			types.Error.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	var cancel context.CancelFunc
	types.Ctx, cancel = context.WithTimeout(context.Background(), wait)

	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(types.Ctx)
	//close connection
	clientapp.Close()

	types.Info.Println("Shutting down")
	os.Exit(0)
}
