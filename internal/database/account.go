package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Account represents the account model.
type Account struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *service) CreateAccount(account *Account) error {
	query := `INSERT INTO accounts (id, username, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(query, account.Id, account.Username, account.Password, account.Role, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return fmt.Errorf("could not create account: %v", err)
	}

	return nil
}

func (s *service) FindAccountById(accountId string) (*Account, error) {
	query := `SELECT id, username, password, role, created_at, updated_at FROM accounts WHERE id = $1`

	account := &Account{}

	row := s.db.QueryRow(query, accountId)
	err := row.Scan(&account.Id, &account.Username, &account.Password, &account.Role, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with ID %s not found", accountId)
		}
		return nil, fmt.Errorf("error finding account: %v", err)
	}

	return account, nil
}

func (s *service) FindAccountByUsername(username string) (*Account, error) {
	query := `SELECT id, username, password, role, created_at, updated_at FROM accounts WHERE username = $1`

	account := &Account{}

	row := s.db.QueryRow(query, username)
	err := row.Scan(&account.Id, &account.Username, &account.Password, &account.Role, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with username %s not found", username)
		}
		return nil, fmt.Errorf("error finding account: %v", err)
	}

	return account, nil
}

func (s *service) FindAllAccounts(page int64, perPage int64) ([]Account, error) {
	offset := (page - 1) * perPage
	query := `SELECT id, username, password, role, created_at, updated_at FROM accounts LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(query, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts: %v", err)
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		account := Account{}
		err := rows.Scan(&account.Id, &account.Username, &account.Password, &account.Role, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning account: %v", err)
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over accounts: %v", err)
	}

	return accounts, nil
}

func (s *service) UpdateAccount(account *Account) error {
	query := `UPDATE accounts SET username = $1, password = $2, role = $3, updated_at = $4 WHERE id = $5`

	_, err := s.db.Exec(query, account.Username, account.Password, account.Role, account.UpdatedAt, account.Id)
	if err != nil {
		return fmt.Errorf("could not update account: %v", err)
	}

	return nil
}

func (s *service) DeleteAccount(accountId string) error {
	query := `DELETE FROM accounts WHERE id = $1`

	_, err := s.db.Exec(query, accountId)
	if err != nil {
		return fmt.Errorf("could not delete account: %v", err)
	}

	return nil
}
