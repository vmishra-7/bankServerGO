package main

import (
	server "bankServerGO/server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	s := server.NewAPIServer(":8080")

	router.HandleFunc("/account", server.MakeHTTPHandleFunc(s.HandleAccount))

	log.Println("Starting up the server at port:", s.ListenAddr)
	http.ListenAndServe(s.ListenAddr, router)
}
