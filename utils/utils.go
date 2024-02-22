package utils

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    uuid.UUID `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time  `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	uuID, _ := uuid.NewUUID()
	return &Account{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    uuID,
		CreatedAt: time.Now().UTC(),
	}
}
