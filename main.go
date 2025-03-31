package main

import (
	"api/online-cinema-theather/internal/config"
	"api/online-cinema-theather/internal/database"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Server started")   

	config.InitAppConfig()
	config.InitDatabaseConfig() 

	database.Connect()
	database.Migrate(database.DB)

	err := http.ListenAndServe(fmt.Sprintf(":%s", config.AppConfig.Port), nil)

	if err != nil {
		log.Println("Server error:", err)
	}
}