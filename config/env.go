package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("failed to load .env", err)
	}
}

func Get(key string) string {

	return os.Getenv(key)
}
