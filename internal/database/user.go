package database

import (
	"fmt"
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Email     string    `json:"email"`
	AccountId string    `json:"acccount_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *service) FindUserById(userId string) (*User, error) {
	query := `SELECT id, name, avatar, email, account_id, created_at, updated_at FROM users where id = $1`

	row := s.db.QueryRow(query, userId)

	user := &User{}

	err := row.Scan(&user.Id, &user.Name, &user.Avatar, &user.Email, &user.AccountId, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("User not found")
	}

	return user, nil
}

func (s *service) FindUserByAccountId(accountId string) (*User, error) {
	query := `SELECT id, name, avatar, email, account_id, created_at, updated_at FROM users where account_id = $1`

	row := s.db.QueryRow(query, accountId)

	user := &User{}

	err := row.Scan(&user.Id, &user.Name, &user.Avatar, &user.Email, &user.AccountId, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("User not found")
	}

	return user, nil
}

func (s *service) FindUserByEmail(email string) (*User, error) {
	query := `SELECT id, name, avatar, email, account_id, created_at, updated_at FROM users where email = $1`

	row := s.db.QueryRow(query, email)

	user := &User{}

	err := row.Scan(&user.Id, &user.Name, &user.Avatar, &user.Email, &user.AccountId, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("User not found")
	}

	return user, nil
}

func (s *service) FindAllUsers(page int64, perPage int64) ([]User, error) {
	offset := page * perPage
	query := `SELECT id, name, avatar, email, account_id, created_at, updated_at FROM users offset $1 limit $2`

	rows, err := s.db.Query(query, perPage, offset)

	if err != nil {
		return nil, fmt.Errorf("Internal Error", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Avatar, &user.Email, &user.AccountId, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *service) CreateUser(user *User) (string, error) {
	query := `INSERT INTO users (id, name, avatar, email, account_id, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.Exec(query, user.Id, user.Name, user.Avatar, user.Email, user.AccountId, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return "", fmt.Errorf("could not create user: %v", err)
	}

	return user.Id, nil
}

func (s *service) UpdateUser(user *User) (*User, error) {
	query := `UPDATE accounts SET name = $1, avatar = $2, email = $3
			, updated_at = $4 WHERE id = $5`

	_, err := s.db.Exec(query, user.Name, user.Avatar, user.Email, user.UpdatedAt, user.Id)
	if err != nil {
		return nil, fmt.Errorf("could not update account: %v", err)
	}

	return user, nil
}

func (s *service) DeleteUser(userId string) (string, error) {
	query := `DELETE FROM users where id = $1`

	_, err := s.db.Exec(query, userId)

	if err != nil {
		return "", fmt.Errorf("Fail to delete user %v", userId)
	}

	return userId, nil
}
