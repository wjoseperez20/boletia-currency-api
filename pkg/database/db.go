package database

import (
	"fmt"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var database *gorm.DB
	var err error

	db_hostname := os.Getenv("POSTGRES_HOST")
	db_name := os.Getenv("POSTGRES_DB")
	db_user := os.Getenv("POSTGRES_USER")
	db_pass := os.Getenv("POSTGRES_PASSWORD")
	db_port := os.Getenv("POSTGRES_PORT")

	dbURl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_pass, db_hostname, db_port, db_name)

	for i := 1; i <= 3; i++ {
		database, err = gorm.Open(postgres.Open(dbURl), &gorm.Config{})
		if err == nil {
			break
		} else {
			log.Printf("Attempt %d: Failed to initialize database. Retrying...", i)
			time.Sleep(3 * time.Second)
		}
	}

	err = database.AutoMigrate(&models.RequestHistory{})
	if err != nil {
		log.Printf("Failed to migrate request history table: %v", err)
		return
	}

	err = database.AutoMigrate(&models.Currency{})
	if err != nil {
		log.Printf("Failed to migrate currency table: %v", err)
		return
	}

	err = database.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("Failed to migrate user table: %v", err)
		return
	}

	DB = database
}
