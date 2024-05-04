package domain

type Sale struct {
	User     *User
	Products []*Product
	ID       uint
}
