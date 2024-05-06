package postgresRepository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/raphael-foliveira/go-table-tests/internal/adapters/dto"
	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
)

type Users struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) *Users {
	return &Users{
		db: db,
	}
}

func (r *Users) FindByEmail(email string) (*domain.User, error) {
	var user dto.User
	err := r.queryWithContext(&user, "SELECT * FROM users WHERE email = $1", email)
	return user.ToDomainUser(), err
}

func (r *Users) FindByUsername(username string) (*domain.User, error) {
	var user dto.User
	err := r.queryWithContext(&user, "SELECT * FROM users WHERE username = $1", username)
	return user.ToDomainUser(), err
}

func (r *Users) Create(user *domain.User) error {
	return r.queryWithContext(user,
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING *",
		user.Username,
		user.Email,
		user.Password,
	)
}

func (r *Users) queryWithContext(target interface{}, query string, args ...interface{}) error {
	tx, err := r.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(query)
	if err != nil {
		return err
	}

	err = stmt.Get(target, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
