package server

import (
	"bankServerGO/storage"
	"bankServerGO/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type apiError struct {
	Error string `json:"error"`
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
	Store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	a := APIServer{
		ListenAddr: listenAddr,
		Store:      store,
	}
	return &a
}

func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.HandleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.HandleCreateAccount(w, r)
	}
	return fmt.Errorf("method now allowed: %+s", r.Method)
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.Store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) HandleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		reqID := mux.Vars(r)["id"]
		id, err := strconv.Atoi(reqID)

		if err != nil {
			return fmt.Errorf("invalid id %s", reqID)
		}

		account, err := s.Store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}
	if r.Method == "DELETE" {
		return s.HandleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccRequest := new(utils.CreateAccountRequest)
	err := json.NewDecoder(r.Body).Decode(createAccRequest)
	if err != nil {
		return err
	}
	account := utils.NewAccount(createAccRequest.FirstName, createAccRequest.LastName)
	s.Store.CreateAccount(account)
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	reqID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(reqID)
	if err != nil {
		return fmt.Errorf("invalid id %s", reqID)
	}
	resp := s.Store.DeletAccount(id)
	if resp != nil {
		return resp
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}
