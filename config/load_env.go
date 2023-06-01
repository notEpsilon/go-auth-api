package config

import (
	"log"

	"github.com/joho/godotenv"
)

// loads environment variables from a `.env` file
// and stops server execution on failure using `os.Exit(1)`
func MustLoadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
}
