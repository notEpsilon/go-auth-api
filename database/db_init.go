package database

import (
	"go-auth/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// globaly accessible database variable as a `*gorm.DB`
var DB *gorm.DB

// initializes and migrates a postgresql database using the `DB_DSN`
// environment variable and stops server execution on failure using `os.Exit(1)`
func MustInit() {
	var err error

	DB, err = gorm.Open(postgres.Open(os.Getenv("DB_DSN")))
	if err != nil {
		log.Fatalln(err)
	}

	if err = DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalln(err)
	}
}
