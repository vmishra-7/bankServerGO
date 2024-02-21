package server

import (
	"bankServerGO/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type apiError struct {
	Error string
}

func MakeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	ListenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	a := APIServer{
		ListenAddr: listenAddr,
	}
	return &a
}

func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	method := r.Method
	if method == "GET" {
		return s.HandleGetAccount(w, r)
	}
	if method == "POST" {
		return s.HandleCreateAccount(w, r)
	}
	if method == "DELETE" {
		return s.HandleDeleteAccount(w, r)
	}
	return fmt.Errorf("method now allowed: %+s", method)
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account := utils.NewAccount("Vaibhav", "Mishra")
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
