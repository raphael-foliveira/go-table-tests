package service_test

import (
	"errors"
	"testing"

	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
	"github.com/raphael-foliveira/go-table-tests/internal/core/service"
	"github.com/raphael-foliveira/go-table-tests/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUp(t *testing.T) (*mocks.MockHasher, *mocks.MockUserRepository, *service.UserService) {
	hasherMock := mocks.NewMockHasher(t)
	userRepositoryMock := mocks.NewMockUserRepository(t)
	userService := service.NewUserService(userRepositoryMock, hasherMock)
	return hasherMock, userRepositoryMock, userService
}

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		expectedData          *service.LoginResponse
		mockFindByEmailResult *domain.User
		mockFindByEmailError  error
		userEmail             string
		userPassword          string
		name                  string
		expectError           bool
		hasherMockResult      bool
	}{
		{
			name:         "Successful login",
			userEmail:    "test@user.com",
			userPassword: "unhashedPassword",
			expectedData: &service.LoginResponse{
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
			expectedData: &service.LoginResponse{
				Username: "testuser",
				Email:    "test@user.com",
			},
			hasherMockResult: false,
			expectError:      true,
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
			expectError:           true,
			mockFindByEmailResult: nil,
			mockFindByEmailError:  errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasherMock, userRepositoryMock, userService := setUp(t)

			userRepositoryMock.EXPECT().
				FindByEmail(tt.userEmail).
				Return(tt.mockFindByEmailResult, tt.mockFindByEmailError)

			if tt.mockFindByEmailError == nil {
				hasherMock.EXPECT().
					Compare(tt.userPassword, tt.mockFindByEmailResult.Password.Value).
					Return(tt.hasherMockResult)
			}

			result, err := userService.Login(tt.userEmail, tt.userPassword)
			assert.Equal(t, tt.expectError, err != nil)

			if !tt.expectError {
				assert.Equal(t, tt.expectedData.Email, result.Email)
				assert.Equal(t, tt.expectedData.Username, result.Username)
			}
		})
	}
}

func TestUserService_Signup(t *testing.T) {
	tests := []struct {
		findByEmailError     error
		findByUsernameError  error
		createError          error
		findByEmailReturn    *domain.User
		findByUsernameReturn *domain.User
		name                 string
		payloadEmail         string
		payloadUsername      string
		payloadPassword      string
		expectError          bool
		skipFindByEmail      bool
		skipCreate           bool
	}{
		{
			name:            "Successful signup",
			payloadEmail:    "test@test.com",
			payloadUsername: "testusername",
			payloadPassword: "valid_password",
		},
		{
			name: "Email already taken",
			findByEmailReturn: &domain.User{
				Email: &domain.Email{
					Value: "taken@email.com",
				},
			},
			payloadEmail:    "taken@email.com",
			payloadUsername: "validusername",
			payloadPassword: "validpassword",
			expectError:     true,
			skipFindByEmail: true,
			skipCreate:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, userRepositoryMock, userService := setUp(t)

			userRepositoryMock.EXPECT().
				FindByEmail(tt.payloadEmail).
				Return(tt.findByEmailReturn, tt.findByEmailError)

			if !tt.skipFindByEmail {
				userRepositoryMock.EXPECT().
					FindByUsername(tt.payloadUsername).
					Return(tt.findByUsernameReturn, tt.findByUsernameError)
			}

			if !tt.skipCreate {
				userRepositoryMock.EXPECT().
					Create(mock.Anything).Return(tt.createError)
			}

			response, err := userService.Signup(&service.SignupPayload{
				Username: tt.payloadUsername,
				Email:    tt.payloadEmail,
				Password: tt.payloadPassword,
			})

			assert.Equal(t, tt.expectError, err != nil)

			if !tt.expectError {
				assert.Equal(t, tt.payloadEmail, response.Email)
				assert.Equal(t, tt.payloadUsername, response.Username)
				return
			}
			assert.Nil(t, response)
		})
	}
}
