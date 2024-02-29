package storage

import (
	"bankServerGO/utils"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage interface { //will help in migrating to any database, refer to server.go line 34
	CreateAccount(*utils.Account) error
	DeletAccount(int) error
	UpdateAccount(*utils.Account) error
	GetAccountByID(int) (*utils.Account, error)
	GetAccountByNumber(uuid.UUID) (*utils.Account, error)
	GetAccounts() ([]*utils.Account, error)
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

func (s *PostgressStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgressStore) CreateAccountTable() error {
	query := `create table if not exists account(
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number UUID,
		balance serial,
		created_at timestamp,
		hashed_password varchar(72)
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgressStore) CreateAccount(account *utils.Account) error {
	query := `Insert into account
	(first_name, last_name, number, balance, created_at, hashed_password)
	values ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(query, account.FirstName,
		account.LastName, account.Number,
		account.Balance, account.CreatedAt, account.HashedPassword)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgressStore) DeletAccount(id int) error {
	_, err := s.db.Exec("delete from account where id = $1", id)
	if err != nil {
		return fmt.Errorf("account %d not found", id)
	}
	return nil
}

func (s *PostgressStore) UpdateAccount(account *utils.Account) error {
	return nil
}

func (s *PostgressStore) GetAccountByNumber(num uuid.UUID) (*utils.Account, error) {
	account := new(utils.Account)
	resp := s.db.QueryRow("select * from account where number = $1", num)
	err := resp.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.HashedPassword)
	
	if err != nil {
		return nil, fmt.Errorf("account %s not found", num)
	}
	return account, nil
}

func (s *PostgressStore) GetAccountByID(id int) (*utils.Account, error) {
	account := new(utils.Account)
	resp := s.db.QueryRow("select * from account where id = $1", id)
	err := resp.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.HashedPassword)
	
	if err != nil {
		return nil, fmt.Errorf("account %d not found", id)
	}
	return account, nil
}

func (s *PostgressStore) GetAccounts() ([]*utils.Account, error) {
	resp, err := s.db.Query("Select * from account")
	if err != nil {
		return nil, err
	}
	accounts := []*utils.Account{}
	for resp.Next() {
		account := new(utils.Account)
		err := resp.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
			&account.HashedPassword)

		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
