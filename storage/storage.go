package storage

import (
	"bankServerGO/utils"
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface { //will help in migrating to any database, refer to server.go line 34
	CreateAccount(*utils.Account) error
	DeletAccount(int) error
	UpdateAccount(*utils.Account) error
	GetAccountByID(int) (*utils.Account, error)
}

type PostgressStore struct {
	db *sql.DB
}

func NewPostgressConnection() (*PostgressStore, error) {
	connStr := "user=postgres dbname=postgres password=test123 sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgressStore{
		db: db,
	}, nil
}

func (s * PostgressStore) CreateAccount(account *utils.Account) error {
	return nil
}

func (s * PostgressStore) DeletAccount(id int) error {
	return nil
}

func (s * PostgressStore) UpdateAccount(account *utils.Account) error {
	return nil
}

func (s * PostgressStore) GetAccountByID(id int) (*utils.Account, error) {
	return nil, nil
}