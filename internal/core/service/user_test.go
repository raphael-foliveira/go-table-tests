package service_test

import (
	"testing"

	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
	"github.com/raphael-foliveira/go-table-tests/internal/core/service"
	"github.com/raphael-foliveira/go-table-tests/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUp(t *testing.T) (*mocks.MockHasher, *mocks.MockUsersRepository, *service.UsersService) {
	hasherMock := mocks.NewMockHasher(t)
	userRepositoryMock := mocks.NewMockUsersRepository(t)
	userService := service.NewUserService(userRepositoryMock, hasherMock)
	return hasherMock, userRepositoryMock, userService
}

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		mockFindByEmailError  error
		expectedError         error
		expectedData          *domain.LoginResponse
		mockFindByEmailResult *domain.User
		userEmail             string
		userPassword          string
		name                  string
		hasherMockResult      bool
	}{
		{
			name:         "Successful login",
			userEmail:    "test@user.com",
			userPassword: "unhashedPassword",
			expectedData: &domain.LoginResponse{
				Username: "testuser",
				Email:    "test@user.com",
			},
			hasherMockResult: true,
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
		{
			name:         "Valid email with invalid password",
			userEmail:    "test@user.com",
			userPassword: "unhashedPassword",
			expectedData: &domain.LoginResponse{
				Username: "testuser",
				Email:    "test@user.com",
			},
			hasherMockResult: false,
			expectedError:    service.ErrInvalidCredentials,
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
		{
			name:                  "invalid email",
			userEmail:             "invalid@email.com",
			userPassword:          "unhashedPassword",
			expectedData:          nil,
			expectedError:         service.ErrInvalidCredentials,
			mockFindByEmailResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasherMock, userRepositoryMock, userService := setUp(t)

			userRepositoryMock.EXPECT().
				FindByEmail(tt.userEmail).
				Return(tt.mockFindByEmailResult, tt.mockFindByEmailError)

			if tt.mockFindByEmailResult != nil {
				hasherMock.EXPECT().
					Compare(tt.userPassword, tt.mockFindByEmailResult.Password.Value).
					Return(tt.hasherMockResult)
			}

			result, err := userService.Login(tt.userEmail, tt.userPassword)
			assert.Equal(t, tt.expectedError != nil, err != nil)

			if tt.expectedError == nil {
				assert.Equal(t, tt.expectedData.Email, result.Email)
				assert.Equal(t, tt.expectedData.Username, result.Username)
				return
			}
			assert.ErrorIs(t, tt.expectedError, err)
		})
	}
}

func TestUserService_Signup(t *testing.T) {
	tests := []struct {
		findByEmailError     error
		findByUsernameError  error
		createError          error
		expectedError        error
		findByEmailReturn    *domain.User
		findByUsernameReturn *domain.User
		name                 string
		payloadEmail         string
		payloadUsername      string
		payloadPassword      string
	}{
		{
			name:            "Successful signup",
			payloadEmail:    "test@test.com",
			payloadUsername: "testusername",
			payloadPassword: "valid_password",
		},
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
		{
			name:            "invalid payload",
			payloadEmail:    "invalid_email",
			payloadPassword: "invp",
			expectedError:   service.ErrInvalidUserPayload,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, userRepositoryMock, userService := setUp(t)

			userRepositoryMock.EXPECT().
				FindByEmail(tt.payloadEmail).
				Return(tt.findByEmailReturn, tt.findByEmailError)

			if tt.findByEmailReturn == nil {
				userRepositoryMock.EXPECT().
					FindByUsername(tt.payloadUsername).
					Return(tt.findByUsernameReturn, tt.findByUsernameError)

				if tt.findByUsernameReturn == nil && tt.expectedError != service.ErrInvalidUserPayload {
					userRepositoryMock.EXPECT().
						Create(mock.Anything).Return(tt.createError)
				}
			}

			response, err := userService.Signup(&domain.SignupPayload{
				Username: tt.payloadUsername,
				Email:    tt.payloadEmail,
				Password: tt.payloadPassword,
			})

			assert.Equal(t, tt.expectedError != nil, err != nil)

			if tt.expectedError == nil {
				assert.Equal(t, tt.payloadEmail, response.Email)
				assert.Equal(t, tt.payloadUsername, response.Username)
				return
			}
			assert.Nil(t, response)
		})
	}
}
