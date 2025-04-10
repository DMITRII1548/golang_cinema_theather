package database

import (
	"api/online-cinema-theather/internal/models"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	log.Println("Running migrations")

	err := db.AutoMigrate(
		&models.Video{},
		&models.Movie{},
		&models.Admin{},
	)

	if err != nil {
		log.Fatal("Migrate error: ", err)
	}

	log.Println("Migrations runned successfully")
}