package service

import (
	"errors"
	"fmt"

	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
)

type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
}

type Hasher interface {
	Compare(givenPassword, hashedPassword string) bool
}

type UserService struct {
	userRepository UserRepository
	hasher         Hasher
}

func NewUserService(repository UserRepository, hasher Hasher) *UserService {
	return &UserService{
		userRepository: repository,
		hasher:         hasher,
	}
}

type LoginResponse struct {
	Username string
	Email    string
}

func (s *UserService) Login(email, password string) (*LoginResponse, error) {
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

	return &LoginResponse{
		Username: foundUser.Username,
		Email:    foundUser.Email.Value,
	}, nil
}

var ErrInvalidCredentials = errors.New("invalid credentials")

type SignupPayload struct {
	Email    string
	Password string
	Username string
}

func (p *SignupPayload) ToDomainUser() *domain.User {
	return &domain.User{
		Username: p.Username,
		Email: &domain.Email{
			Value: p.Email,
		},
		Password: &domain.Password{
			Value: p.Password,
		},
	}
}

type SignupResponse struct {
	Username string
	Email    string
	ID       uint
}

func NewSignupResponse(user *domain.User) *SignupResponse {
	return &SignupResponse{
		Username: user.Username,
		Email:    user.Email.Value,
		ID:       user.ID,
	}
}

func (s *UserService) Signup(payload *SignupPayload) (*SignupResponse, error) {
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

var ErrInvalidUserPayload = errors.New("invalid user payload")

func (s *UserService) checkIfUserAlreadyExists(username, email string) error {
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
)
