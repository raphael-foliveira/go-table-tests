package ports

import (
	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
)

type UsersRepository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
}

type Hasher interface {
	Compare(givenPassword, hashedPassword string) bool
	Hash(password string) string
}

type UsersService interface {
	Login(email, password string) (*domain.LoginResponse, error)
	Signup()
}
