package main

import (
	"net/http"

	"github.com/desteves/fle/api"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)

	// foo handles encrypted docs
	router.HandleFunc("/foo", api.CreateEncryptedFoobarHandler).Methods("POST")
	router.HandleFunc("/foo/{id}", api.ReadEncryptedFoobarHandler).Methods("GET")

	// bar handles unencrypted docs and won't show the values of encrypted fields/docs
	router.HandleFunc("/bar", api.CreateFoobarHandler).Methods("POST")
	router.HandleFunc("/bar/{id}", api.ReadFoobarHandler).Methods("GET")

	if err := http.ListenAndServe(":8888", router); err != nil {
		panic(err)
	}
}
