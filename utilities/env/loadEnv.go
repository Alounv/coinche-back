package env

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) {
	if path == "" {
		path = ".env"
	}

	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Error 1 loading .env file")
	}
}
