package main

import (
	srv "bankServerGO/server"
	"bankServerGO/storage"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	store, err := storage.NewPostgressConnection()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	s := srv.NewAPIServer(":8080", store)

	router.HandleFunc("/account", srv.MakeHTTPHandleFunc(s.HandleAccount))
	router.HandleFunc("/account/{id}", srv.WithJWT(srv.MakeHTTPHandleFunc(s.HandleGetAccountByID), s.Store))
	router.HandleFunc("/transfer", srv.MakeHTTPHandleFunc(s.HandleTransferRequest))

	log.Println("Starting up the server at port:", s.ListenAddr)
	http.ListenAndServe(s.ListenAddr, router)
}
