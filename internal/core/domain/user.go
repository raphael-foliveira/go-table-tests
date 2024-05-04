package domain

import "strings"

type Password struct {
	Value    string
	IsHashed bool
}

func (p *Password) IsValid() bool {
	if p.IsHashed {
		return false
	}
	return len(p.Value) > 5
}

type Email struct {
	Value string
}

func (e *Email) IsValid() bool {
	return strings.Contains(e.Value, "@")
}

type User struct {
	Password *Password
	Email    *Email
	Username string
}

func (u *User) IsValid() bool {
	return u.Email.IsValid() && u.Password.IsValid()
}
