package handlers

import (
	"api/online-cinema-theather/internal/database"
	"api/online-cinema-theather/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	var movies []models.Movie

	if err := database.DB.Preload("Video").Limit(30).Find(&movies).Error; err != nil {
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		return
	}

	sendJSON(w, http.StatusOK, movies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	
	if err := database.DB.Create(&movie).Error; err != nil {
		http.Error(w, "Failed to create movie", http.StatusInternalServerError)
		return
	}

	sendJSON(w, http.StatusCreated, movie)
}