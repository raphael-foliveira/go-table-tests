package security

import "golang.org/x/crypto/bcrypt"

type HashingService struct{}

func NewHashingService() *HashingService {
	return &HashingService{}
}

func (s *HashingService) Hash(str string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *HashingService) Compare(givenPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(givenPassword))
	return err != nil
}
