package domain_test

import (
	"testing"

	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestUser_Email(t *testing.T) {
	tests := []struct {
		name        string
		email       *domain.Email
		expectError bool
	}{
		{
			name: "Valid Email",
			email: &domain.Email{
				Value: "valid@email.com",
			},
			expectError: false,
		},
		{
			name: "Invalid Email",
			email: &domain.Email{
				Value: "invalid_email",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.email.Validate()
			assert.Equal(t, tt.expectError, err != nil)
		})
	}
}

func TestUser_Password(t *testing.T) {
	tests := []struct {
		name        string
		password    *domain.Password
		expectError bool
	}{
		{
			name: "Valid Password",
			password: &domain.Password{
				Value:    "valid_password",
				IsHashed: false,
			},
			expectError: false,
		},
		{
			name: "Invalid Password",
			password: &domain.Password{
				Value:    "invp",
				IsHashed: false,
			},
			expectError: true,
		},
		{
			name: "Password already hashed",
			password: &domain.Password{
				Value:    "valid_password",
				IsHashed: true,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.password.Validate()
			assert.Equal(t, tt.expectError, err != nil)
		})
	}
}

func TestUser_User(t *testing.T) {
	tests := []struct {
		name        string
		user        *domain.User
		expectError bool
	}{
		{
			name: "Valid User",
			user: &domain.User{
				Username: "valid_user",
				Email: &domain.Email{
					Value: "valid@email.com",
				},
				Password: &domain.Password{
					Value:    "valid_password_string",
					IsHashed: false,
				},
			},
			expectError: false,
		},
		{
			name: "Invalid User",
			user: &domain.User{
				Username: "valid_user",
				Email: &domain.Email{
					Value: "invalid_email.com",
				},
				Password: &domain.Password{
					Value:    "valid_password_string",
					IsHashed: false,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		err := tt.user.Validate()
		assert.Equal(t, tt.expectError, err != nil)
	}
}
