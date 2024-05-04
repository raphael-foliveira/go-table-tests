package service_test

import (
	"errors"
	"testing"

	"github.com/raphael-foliveira/login-service-hexagonal/internal/core/domain"
	"github.com/raphael-foliveira/login-service-hexagonal/internal/core/service"
	"github.com/raphael-foliveira/login-service-hexagonal/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

var userStub = &domain.User{
	Username: "Test User",
	Email: &domain.Email{
		Value: "test@user.com",
	},
	Password: &domain.Password{
		Value:    "hashedTestPassword",
		IsHashed: true,
	},
}

func setUp(t *testing.T) (*mocks.MockHasher, *mocks.MockUserRepository, *service.LoginService) {
	hasherMock := mocks.NewMockHasher(t)
	userRepositoryMock := mocks.NewMockUserRepository(t)
	service := service.New(userRepositoryMock, hasherMock)
	return hasherMock, userRepositoryMock, service
}

func TestLoginService_Login(t *testing.T) {
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
			hasherMock, userRepositoryMock, service := setUp(t)

			userRepositoryMock.
				On("FindByEmail", tt.userEmail).
				Return(tt.mockFindByEmailResult, tt.mockFindByEmailError)

			if tt.mockFindByEmailError == nil {
				hasherMock.
					On("Compare", tt.userPassword, tt.mockFindByEmailResult.Password.Value).
					Return(tt.hasherMockResult)
			}

			result, err := service.Login(tt.userEmail, tt.userPassword)
			assert.Equal(t, tt.expectError, err != nil)

			if !tt.expectError {
				assert.Equal(t, tt.expectedData.Email, result.Email)
				assert.Equal(t, tt.expectedData.Username, result.Username)
			}
		})
	}
}
