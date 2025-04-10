package main

import (
	"api/online-cinema-theather/internal/config"
	"api/online-cinema-theather/internal/database"
	"api/online-cinema-theather/internal/models"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

const defaultCost = 13

func main() {
	config.InitAppConfig()
	config.InitDatabaseConfig()

	database.Connect()

	migrateAdminTable()

	login, password := scanAdmin()

	password, err := hashPassword(password)

	if err != nil {
		fmt.Print(err.Error())
	}

	createAdmin(login, password)
}

func migrateAdminTable() {
	log.Println("Running migrations")

	err := database.DB.AutoMigrate(
		&models.Admin{},
	)

	if err != nil {
		log.Fatal("Migrate error: ", err)
	}

	log.Println("Migrations runned successfully")
}

func hashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
    return string(hash), err
}


func scanAdmin() (string, string) {
	var login string
	var password string

	fmt.Print("Input admin's login:")
	fmt.Scan(&login)

	fmt.Print("Input admin's password: ")
	fmt.Scan(&password)

	return login, password
}

func createAdmin(login string, passwordHash string) {
	admin := models.Admin{Login: login, Password: passwordHash}

	if err := database.DB.Create(&admin).Error; err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Created Admin successfully")
}