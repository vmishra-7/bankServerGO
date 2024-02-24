package utils

import (
	"time"

	"github.com/google/uuid"
)

type TranserRequest struct {
	AccountNumber uuid.UUID `json:"accountNumber"`
	Amount        int       `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    uuid.UUID `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	uuID, _ := uuid.NewUUID()
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Number:    uuID,
		CreatedAt: time.Now().UTC(),
	}
}
