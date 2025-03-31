package main

import (
	"api/online-cinema-theather/internal/config"
	"api/online-cinema-theather/internal/database"
	"api/online-cinema-theather/internal/routes"
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

	mux := routes.RegisterRoutes()

	err := http.ListenAndServe(fmt.Sprintf(":%s", config.AppConfig.Port), mux)

	if err != nil {
		log.Println("Server error:", err)
	}
}