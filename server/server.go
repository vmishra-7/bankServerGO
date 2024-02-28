package server

import (
	"bankServerGO/storage"
	"bankServerGO/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (s *APIServer) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method now allowed: %+s", r.Method)
	}

	loginReq := new(utils.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(loginReq); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusAccepted, loginReq)
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
	accReq := new(utils.CreateAccountRequest)
	err := json.NewDecoder(r.Body).Decode(accReq)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	account, err := utils.NewAccount(accReq.FirstName, accReq.LastName, accReq.Password)
	if err != nil {
		return err
	}
	if err := s.Store.CreateAccount(account); err != nil {
		return err
	}
	
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

func (s *APIServer) HandleTransferRequest(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(utils.TranserRequest)
	err := json.NewDecoder(r.Body).Decode(transferReq)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
}

func WithJWT(handlerFunc http.HandlerFunc, s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Into JWT middleware!")

		tokenString := r.Header.Get("Authorization")
		token, err := validateJWT(tokenString)

		if err != nil {
			WriteJSON(w, http.StatusForbidden, apiError{Error: "permission denied"})
			return
		}

		if !token.Valid {
			WriteJSON(w, http.StatusForbidden, apiError{Error: "permission denied"})
			return
		}

		userID := mux.Vars(r)["id"]
		id, err := strconv.Atoi(userID)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, apiError{Error: "permission denied"})
			return
		}
		account, err := s.GetAccountByID(id)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, apiError{Error: "permission denied"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		accNum, _ := uuid.Parse(claims["accountNumber"].(string))
		if account.Number != accNum {
			WriteJSON(w, http.StatusForbidden, apiError{Error: "permission denied"})
			return
		}

		handlerFunc(w, r)
	}
}

func createJWT(account *utils.Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     time.Now().Add(time.Hour * 24),
		"accountNumber": account.Number,
	}

	secret := os.Getenv("jsonAPISecretKEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv("jsonAPISecretKEY")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secretKey), nil
	})
}
