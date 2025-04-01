package handlers

import (
	"api/online-cinema-theather/internal/database"
	"api/online-cinema-theather/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func getMovieId(w http.ResponseWriter, r *http.Request) int64 {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 0)

	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest) 
		return 0
	}

	return id
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	var movies []models.Movie

	if err := database.DB.Preload("Video").Limit(30).Find(&movies).Error; err != nil {
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		return
	}

	sendJSON(w, http.StatusOK, movies)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	id := getMovieId(w, r)

	if id == 0 {
		return
	}

	if err := database.DB.First(&movie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Movie not found", http.StatusNotFound) 
		} else {
			http.Error(w, "Failed to get a movie", http.StatusInternalServerError) 
		}
		return
	}

	sendJSON(w, http.StatusOK, movie)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id := getMovieId(w, r)

	if id == 0 {
		return
	}

	result := database.DB.Delete(&models.Movie{}, id)

	if result.RowsAffected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return 
	}

	if result.Error != nil {
		http.Error(w, "Interanl server error", http.StatusInternalServerError)
		return        
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "Movie deleted successfully"})
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
