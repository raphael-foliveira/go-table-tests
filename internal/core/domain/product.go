package domain

type Manufacturer int

const (
	Yamaha Manufacturer = iota
)

type Product struct {
	Model        string
	Manufacturer Manufacturer
	ID           uint
	Price        float32
}
