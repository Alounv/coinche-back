package utilities

import (
	"github.com/joho/godotenv"
)

func LoadEnv(path string) {
	if path == "" {
		path = ".env"
	}

	err := godotenv.Load(path)
	if err != nil {
		panic("Error loading .env file")
	}
}
