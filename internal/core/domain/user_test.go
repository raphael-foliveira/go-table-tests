package domain_test

import (
	"testing"

	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestUser_Email(t *testing.T) {
	t.Run("Valid Email", func(t *testing.T) {
		validEmail := &domain.Email{
			Value: "valid@email.com",
		}

		assert.True(t, validEmail.IsValid())
	})

	t.Run("Invalid Email", func(t *testing.T) {
		invalidEmail := &domain.Email{
			Value: "invalid_email",
		}

		assert.False(t, invalidEmail.IsValid())
	})
}

func TestUser_Password(t *testing.T) {
	t.Run("Valid Password", func(t *testing.T) {
		validPassword := &domain.Password{
			Value:    "valid_password",
			IsHashed: false,
		}

		assert.True(t, validPassword.IsValid())
	})

	t.Run("Invalid Password", func(t *testing.T) {
		invalidPassword := &domain.Password{
			Value:    "invp",
			IsHashed: false,
		}

		assert.False(t, invalidPassword.IsValid())
	})

	t.Run("Already hashed password", func(t *testing.T) {
		alreadyHashedPassword := &domain.Password{
			Value:    "hashed_password_string",
			IsHashed: true,
		}

		assert.False(t, alreadyHashedPassword.IsValid())
	})
}

func TestUser_User(t *testing.T) {
	t.Run("Valid User", func(t *testing.T) {
		validUser := &domain.User{
			Username: "valid_user",
			Email: &domain.Email{
				Value: "valid@email.com",
			},
			Password: &domain.Password{
				Value:    "valid_password_string",
				IsHashed: false,
			},
		}

		assert.True(t, validUser.IsValid())
	})

	t.Run("Invalid User", func(t *testing.T) {
		validUser := &domain.User{
			Username: "valid_user",
			Email: &domain.Email{
				Value: "valid_email.com",
			},
			Password: &domain.Password{
				Value:    "valid_password_string",
				IsHashed: false,
			},
		}

		assert.False(t, validUser.IsValid())
	})
}
