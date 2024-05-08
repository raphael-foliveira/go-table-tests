package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func NewConfig(filePath string) (*Config, error) {
	err := godotenv.Load(filePath)
	if err != nil {
		return nil, err
	}

	err = validateConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}, nil
}

func validateConfig() error {
	requiredVariables := []string{
		"DATABASE_URL",
	}

	for _, variable := range requiredVariables {
		if os.Getenv(variable) == "" {
			return fmt.Errorf("%w: %s", ErrRequiredVariableNotSet, variable)
		}
	}
	return nil
}

var ErrRequiredVariableNotSet = errors.New("required variable not set")
