package main

import (
	server "bankServerGO/server"
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
	router := mux.NewRouter()
	s := server.NewAPIServer(":8080", store)

	router.HandleFunc("/account", server.MakeHTTPHandleFunc(s.HandleAccount))

	log.Println("Starting up the server at port:", s.ListenAddr)
	http.ListenAndServe(s.ListenAddr, router)
}
