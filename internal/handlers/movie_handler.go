package handlers

import (
	"api/online-cinema-theather/internal/database"
	"api/online-cinema-theather/internal/models"
	"encoding/json"
	"net/http"
)

func GetMovies(w http.ResponseWriter, r *http.Request) {
	var movies []models.Movie

	result := database.DB.Preload("Video").Limit(30).Find(&movies)

	if result.Error != nil {
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}