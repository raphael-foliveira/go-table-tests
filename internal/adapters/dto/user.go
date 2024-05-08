package dto

import "github.com/raphael-foliveira/go-table-tests/internal/core/domain"

type SignupPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	ID       uint   `json:"id"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type User struct {
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
	ID       uint   `db:"id"`
}

func (u *User) ToDomainUser() *domain.User {
	return &domain.User{
		Username: u.Username,
		Email: &domain.Email{
			Value: u.Email,
		},
		Password: &domain.Password{
			Value:    u.Password,
			IsHashed: true,
		},
	}
}
