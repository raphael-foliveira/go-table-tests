package service

import (
	"errors"

	"github.com/raphael-foliveira/login-service-hexagonal/internal/core/domain"
)

type UserRepository interface {
	FindByEmail(string) (*domain.User, error)
}

type Hasher interface {
	Compare(givenPassword, hashedPassword string) bool
}

type LoginService struct {
	userRepository UserRepository
	hasher         Hasher
}

func New(repository UserRepository, hasher Hasher) *LoginService {
	return &LoginService{
		userRepository: repository,
		hasher:         hasher,
	}
}

type LoginResponse struct {
	Username string
	Email    string
}

func (ls *LoginService) Login(email, password string) (*LoginResponse, error) {
	foundUser, err := ls.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	ok := ls.hasher.Compare(password, foundUser.Password.Value)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	return &LoginResponse{
		Username: foundUser.Username,
		Email:    foundUser.Email.Value,
	}, nil
}

var ErrInvalidCredentials = errors.New("invalid credentials")
