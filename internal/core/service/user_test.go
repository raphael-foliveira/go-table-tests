package service_test

import (
	"testing"

	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
	"github.com/raphael-foliveira/go-table-tests/internal/core/service"
	"github.com/raphael-foliveira/go-table-tests/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUp(t *testing.T) (*mocks.MockHasher, *mocks.MockUsersRepository, *service.Users) {
	hasherMock := mocks.NewMockHasher(t)
	userRepositoryMock := mocks.NewMockUsersRepository(t)
	userService := service.NewUsersService(userRepositoryMock, hasherMock)
	return hasherMock, userRepositoryMock, userService
}

func TestUserServiceLogin_Success(t *testing.T) {
	tests := []struct {
		expectedData          *domain.LoginResponse
		mockFindByEmailResult *domain.User
		userEmail             string
		userPassword          string
		name                  string
	}{
		{
			name:         "Successful login",
			userEmail:    "test@user.com",
			userPassword: "unhashedPassword",
			expectedData: &domain.LoginResponse{
				Username: "testuser",
				Email:    "test@user.com",
			},
			mockFindByEmailResult: &domain.User{
				Email: &domain.Email{
					Value: "test@user.com",
				},
				Password: &domain.Password{
					IsHashed: true,
					Value:    "hashedPassword",
				},
				Username: "testuser",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasherMock, userRepositoryMock, userService := setUp(t)

			userRepositoryMock.EXPECT().
				FindByEmail(tt.userEmail).
				Return(tt.mockFindByEmailResult, nil)

			hasherMock.EXPECT().
				Compare(tt.userPassword, tt.mockFindByEmailResult.Password.Value).
				Return(true)

			result, err := userService.Login(tt.userEmail, tt.userPassword)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedData.Email, result.Email)
			assert.Equal(t, tt.expectedData.Username, result.Username)
		})
	}
}

func TestUserServiceLogin_InvalidPassword(t *testing.T) {
	tests := []struct {
		expectedData          *domain.LoginResponse
		mockFindByEmailResult *domain.User
		userEmail             string
		name                  string
		hasherMockResult      bool
	}{
		{
			name:      "Valid email with invalid password",
			userEmail: "test@user.com",
			expectedData: &domain.LoginResponse{
				Username: "testuser",
				Email:    "test@user.com",
			},
			hasherMockResult: false,
			mockFindByEmailResult: &domain.User{
				Email: &domain.Email{
					Value: "test@user.com",
				},
				Password: &domain.Password{
					IsHashed: true,
					Value:    "hashedPassword",
				},
				Username: "testuser",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasherMock, userRepositoryMock, userService := setUp(t)

			userRepositoryMock.EXPECT().
				FindByEmail(tt.userEmail).
				Return(tt.mockFindByEmailResult, nil)

			hasherMock.EXPECT().
				Compare(mock.Anything, mock.Anything).
				Return(tt.hasherMockResult)

			_, err := userService.Login(tt.userEmail, "invalidPassword")
			assert.ErrorIs(t, service.ErrInvalidCredentials, err)
		})
	}
}

func TestUserServiceLogin_InvalidEmail(t *testing.T) {
	tests := []struct {
		userEmail    string
		userPassword string
		name         string
	}{
		{
			name:         "invalid email",
			userEmail:    "invalid@email.com",
			userPassword: "unhashedPassword",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, userRepositoryMock, userService := setUp(t)

			userRepositoryMock.EXPECT().
				FindByEmail(tt.userEmail).
				Return(nil, nil)

			_, err := userService.Login(tt.userEmail, tt.userPassword)
			assert.ErrorIs(t, service.ErrInvalidCredentials, err)
		})
	}
}

func TestUserServiceSignup_InvalidPayload(t *testing.T) {
	tests := []struct {
		expectedError   error
		name            string
		payloadEmail    string
		payloadUsername string
		payloadPassword string
	}{
		{
			name:            "invalid password",
			payloadEmail:    "valid@email.com",
			payloadPassword: "invp",
			expectedError:   domain.ErrPasswordTooShort,
		},
		{
			name:            "invalid email",
			payloadEmail:    "invalid_email.com",
			payloadPassword: "valid_password",
			expectedError:   domain.ErrEmailInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, userService := setUp(t)

			response, err := userService.Signup(&domain.SignupPayload{
				Username: tt.payloadUsername,
				Email:    tt.payloadEmail,
				Password: tt.payloadPassword,
			})

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Nil(t, response)
		})
	}
}

func TestUserServiceSignup_Success(t *testing.T) {
	tests := []struct {
		name            string
		payloadEmail    string
		payloadUsername string
		payloadPassword string
	}{
		{
			name:            "Successful signup",
			payloadEmail:    "test@test.com",
			payloadUsername: "testusername",
			payloadPassword: "valid_password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasherMock, userRepositoryMock, userService := setUp(t)

			userRepositoryMock.EXPECT().
				FindByEmail(tt.payloadEmail).
				Return(nil, nil)

			userRepositoryMock.EXPECT().
				FindByUsername(tt.payloadUsername).
				Return(nil, nil)

			hasherMock.EXPECT().Hash(tt.payloadPassword).Return("hashedPassword", nil)
			userRepositoryMock.EXPECT().
				Create(mock.Anything).Return(nil)

			response, err := userService.Signup(&domain.SignupPayload{
				Username: tt.payloadUsername,
				Email:    tt.payloadEmail,
				Password: tt.payloadPassword,
			})

			assert.NoError(t, err)
			assert.Equal(t, tt.payloadEmail, response.Email)
			assert.Equal(t, tt.payloadUsername, response.Username)
		})
	}
}

func TestUserServiceSignup_Conflict(t *testing.T) {
	tests := []struct {
		expectedError        error
		findByEmailReturn    *domain.User
		findByUsernameReturn *domain.User
		name                 string
		payloadEmail         string
		payloadUsername      string
		payloadPassword      string
	}{
		{
			name:              "Email already taken",
			findByEmailReturn: &domain.User{},
			payloadEmail:      "taken@email.com",
			payloadUsername:   "validusername",
			payloadPassword:   "validpassword",
			expectedError:     service.ErrEmailAlreadyTaken,
		},
		{
			name:                 "Username already taken",
			findByUsernameReturn: &domain.User{},
			payloadEmail:         "valid@email.com",
			payloadUsername:      "takenusername",
			payloadPassword:      "validpassword",
			expectedError:        service.ErrUsernameAlreadyTaken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, userRepositoryMock, userService := setUp(t)

			userRepositoryMock.EXPECT().
				FindByEmail(tt.payloadEmail).
				Return(tt.findByEmailReturn, nil)

			if tt.findByEmailReturn == nil {
				userRepositoryMock.EXPECT().
					FindByUsername(tt.payloadUsername).
					Return(tt.findByUsernameReturn, nil)
			}

			response, err := userService.Signup(&domain.SignupPayload{
				Username: tt.payloadUsername,
				Email:    tt.payloadEmail,
				Password: tt.payloadPassword,
			})

			assert.ErrorIs(t, tt.expectedError, err)

			assert.Nil(t, response)
		})
	}
}
