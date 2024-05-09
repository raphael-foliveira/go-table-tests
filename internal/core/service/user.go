package service

import (
	"errors"
	"fmt"

	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
	"github.com/raphael-foliveira/go-table-tests/internal/core/ports"
)

type Users struct {
	userRepository ports.UsersRepository
	hasher         ports.Hasher
}

func NewUsersService(repository ports.UsersRepository, hasher ports.Hasher) *Users {
	return &Users{
		userRepository: repository,
		hasher:         hasher,
	}
}

func (s *Users) Login(email, password string) (*domain.LoginResponse, error) {
	foundUser, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	ok := s.hasher.Compare(password, foundUser.Password.Value)
	if foundUser == nil || !ok {
		return nil, ErrInvalidCredentials
	}

	return &domain.LoginResponse{
		Username: foundUser.Username,
		Email:    foundUser.Email.Value,
	}, nil
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func NewSignupResponse(user *domain.User) *domain.SignupResponse {
	return &domain.SignupResponse{
		Username: user.Username,
		Email:    user.Email.Value,
		ID:       user.ID,
	}
}

func (s *Users) Signup(payload *domain.SignupPayload) (*domain.SignupResponse, error) {
	userToCreate := payload.ToDomainUser()
	if err := userToCreate.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidUserPayload, err)
	}

	err := s.checkIfUserAlreadyExists(userToCreate)
	if err != nil {
		return nil, err
	}

	err = s.hashPassword(userToCreate.Password)
	if err != nil {
		return nil, err
	}

	err = s.userRepository.Create(userToCreate)
	if err != nil {
		return nil, err
	}

	return NewSignupResponse(userToCreate), nil
}

func (s *Users) checkIfUserAlreadyExists(user *domain.User) error {
	foundUser, err := s.userRepository.FindByEmail(user.Email.Value)
	if err != nil {
		return err
	}
	if foundUser != nil {
		return ErrEmailAlreadyTaken
	}

	foundUser, err = s.userRepository.FindByUsername(user.Username)
	if err != nil {
		return err
	}
	if foundUser != nil {
		return ErrUsernameAlreadyTaken
	}

	return nil
}

func (s *Users) hashPassword(password *domain.Password) error {
	if password.IsHashed {
		return ErrPasswordAlreadyHashed
	}
	hashedPassword, err := s.hasher.Hash(password.Value)
	if err != nil {
		return err
	}
	password.Value = hashedPassword
	password.IsHashed = true
	return nil
}

var (
	ErrEmailAlreadyTaken     = errors.New("email already taken")
	ErrUsernameAlreadyTaken  = errors.New("username already taken")
	ErrInvalidUserPayload    = errors.New("invalid user payload")
	ErrPasswordAlreadyHashed = errors.New("cannot hash password that is already hashed")
)
