package utilities

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) {
	if path == "" {
		path = ".env"
	}

	err := godotenv.Load(path)
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}
