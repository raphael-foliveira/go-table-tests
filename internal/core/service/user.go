package service

import (
	"errors"
	"fmt"

	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
	"github.com/raphael-foliveira/go-table-tests/internal/core/ports"
)

type UsersService struct {
	userRepository ports.UsersRepository
	hasher         ports.Hasher
}

func NewUsersService(repository ports.UsersRepository, hasher ports.Hasher) *UsersService {
	return &UsersService{
		userRepository: repository,
		hasher:         hasher,
	}
}

func (s *UsersService) Login(email, password string) (*domain.LoginResponse, error) {
	foundUser, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if foundUser == nil {
		return nil, ErrInvalidCredentials
	}

	ok := s.hasher.Compare(password, foundUser.Password.Value)
	if !ok {
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

func (s *UsersService) Signup(payload *domain.SignupPayload) (*domain.SignupResponse, error) {
	err := s.checkIfUserAlreadyExists(payload.Username, payload.Email)
	if err != nil {
		return nil, err
	}

	userToCreate := payload.ToDomainUser()
	if err := userToCreate.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidUserPayload, err)
	}

	err = s.userRepository.Create(userToCreate)
	if err != nil {
		return nil, err
	}

	return NewSignupResponse(userToCreate), nil
}

func (s *UsersService) checkIfUserAlreadyExists(username, email string) error {
	foundUser, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return err
	}
	if foundUser != nil {
		return ErrEmailAlreadyTaken
	}

	foundUser, err = s.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if foundUser != nil {
		return ErrUsernameAlreadyTaken
	}

	return nil
}

var (
	ErrEmailAlreadyTaken    = errors.New("email already taken")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrInvalidUserPayload   = errors.New("invalid user payload")
)
