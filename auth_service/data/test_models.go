package data

import (
	"database/sql"
)

type PostgresTestRepository struct {
	Conn *sql.DB
}

func NewPostgresTestRepository(db *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{
		Conn: db,
	}
}

func (p *PostgresTestRepository) GetAll() ([]*User, error) {
	users := []*User{}
	return users, nil
}

func (p *PostgresTestRepository) GetByEmail(email string) (*User, error) {
	user := User{}
	return &user, nil
}

func (p *PostgresTestRepository) GetOne(id int) (*User, error) {
	user := User{}
	return &user, nil
}

func (p *PostgresTestRepository) Update(u User) error {
	return nil
}

func (p *PostgresTestRepository) DeleteByID(id int) error {
	return nil
}

func (p *PostgresTestRepository) Insert(user User) (int, error) {
	newID := 1
	return newID, nil
}

func (p *PostgresTestRepository) ResetPassword(password string, u User) error {
	return nil
}

func (p *PostgresTestRepository) PasswordMatches(plainText string, u User) (bool, error) {
	return true, nil
}
