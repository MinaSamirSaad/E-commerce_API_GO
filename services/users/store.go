package users

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MinaSamirSaad/ecommerce/services/shared"
)

type UserStore interface {
	GetUserByEmail(email string) (*shared.User, error)
	CreateUser(u *shared.User) (*shared.User, error)
	GetUserByID(id int) (*shared.User, error)
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*shared.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	u := &shared.User{}
	for rows.Next() {
		u, err = ScanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func ScanRowIntoUser(rows *sql.Rows) (*shared.User, error) {
	u := &shared.User{}
	err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) CreateUser(u *shared.User) (*shared.User, error) {
	result, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		return nil, err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = int(lastID)

	u.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	return u, nil
}
func (s *Store) GetUserByID(id int) (*shared.User, error) {
	return nil, nil
}
