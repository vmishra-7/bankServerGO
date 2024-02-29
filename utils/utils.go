package utils

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Number   uuid.UUID `json:"number"`
	Passowrd string    `json:"password"`
}

type LoginResponse struct {
	Number uuid.UUID `json:"number"`
	Token  string    `json:"token"`
}
type TranserRequest struct {
	ToAccount uuid.UUID `json:"toAccount"`
	Amount    int       `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type Account struct {
	ID             int       `json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Number         uuid.UUID `json:"number"`
	Balance        int64     `json:"balance"`
	CreatedAt      time.Time `json:"createdAt"`
	HashedPassword string    `json:"-"`
}

func NewAccount(firstName, lastName, passowrd string) (*Account, error) {
	hsdPswd, err := bcrypt.GenerateFromPassword([]byte(passowrd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	uuID, _ := uuid.NewUUID()
	return &Account{
		FirstName:      firstName,
		LastName:       lastName,
		Number:         uuID,
		CreatedAt:      time.Now().UTC(),
		HashedPassword: string(hsdPswd),
	}, nil
}

func (a *Account) ValidatePassword(pswd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.HashedPassword), []byte(pswd)) == nil
}
