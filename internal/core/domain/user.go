package domain

import (
	"errors"
	"strings"
)

type Password struct {
	Value    string
	IsHashed bool
}

func (p *Password) Validate() error {
	if p.IsHashed {
		return ErrPasswordAlreadyHashed
	}
	if len(p.Value) < 6 {
		return ErrPasswordTooShort
	}
	return nil
}

var (
	ErrPasswordTooShort      = errors.New("password is too short")
	ErrPasswordAlreadyHashed = errors.New("password already hashed")
)

type Email struct {
	Value string
}

func (e *Email) Validate() error {
	if !strings.Contains(e.Value, "@") {
		return ErrEmailInvalid
	}
	return nil
}

var ErrEmailInvalid = errors.New("email is not valid")

type User struct {
	Password *Password
	Email    *Email
	Username string
	ID       uint
}

func (u *User) Validate() error {
	if err := u.Email.Validate(); err != nil {
		return err
	}
	if err := u.Password.Validate(); err != nil {
		return err
	}
	return nil
}
