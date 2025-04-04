package handlers

import (
	"api/online-cinema-theather/internal/database"
	"api/online-cinema-theather/internal/models"
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
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

func generateRandomFileName(originalFilename string) string {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	randomString := hex.EncodeToString(randomBytes)
	ext := filepath.Ext(originalFilename)
	return randomString + ext
}

func putImageFile(w http.ResponseWriter, r *http.Request, uploadDir string, fieldName string) string {
	file, handler, err := r.FormFile(fieldName)
	if err != nil {
		http.Error(w, "The '"+fieldName+"' field is required", http.StatusUnprocessableEntity)
		return ""
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusInternalServerError)
		return ""
	}

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" && contentType != "image/webp" {
		http.Error(w, "Only image files are allowed", http.StatusUnsupportedMediaType)
		return ""
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		http.Error(w, "Failed to rewind uploaded file", http.StatusInternalServerError)
		return ""
	}

	randomName := generateRandomFileName(handler.Filename)
	path := filepath.Join(uploadDir, randomName)

	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "Failed to save "+fieldName, http.StatusInternalServerError)
		return ""
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to write "+fieldName+" file", http.StatusInternalServerError)
		return ""
	}

	return path
}


func getMovieId(w http.ResponseWriter, r *http.Request) int64 {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 0)

	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest) 
		return 0
	}

	return id
}

func getMovie(w http.ResponseWriter, id int64) (*models.Movie, bool) {
	var movie models.Movie

	if err := database.DB.First(&movie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Movie not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch movie", http.StatusInternalServerError)
		}
		return nil, false
	}

	return &movie, true
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
	id := getMovieId(w, r)
	if id == 0 {
		return
	}

	movie, ok := getMovie(w, id)
	if !ok {
		return
	}

	sendJSON(w, http.StatusOK, movie)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id := getMovieId(w, r)
	if id == 0 {
		return
	}

	movie, ok := getMovie(w, id)
	if !ok {
		return
	}

	result := database.DB.Delete(movie)
	if result.Error != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	os.Remove(movie.Thumbnail)
	os.Remove(movie.Preview)

	sendJSON(w, http.StatusOK, map[string]string{"message": "Movie deleted successfully"})
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 20)

	thumbnailFilePath := putImageFile(w, r, "storage/imgs/thumbnails/", "thumbnail")
	if thumbnailFilePath == "" {
		return
	}

	previewFilePath := putImageFile(w, r, "storage/imgs/previews/", "preview")
	if previewFilePath == "" {
		os.Remove(thumbnailFilePath)
		return
	}

	videoID, err := strconv.ParseUint(r.FormValue("video_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid video_id", http.StatusBadRequest)
		os.Remove(thumbnailFilePath)
		os.Remove(previewFilePath)
		return
	}

	movie := models.Movie{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		VideoID:     uint(videoID),
		Thumbnail:   thumbnailFilePath,
		Preview:     previewFilePath,
	}

	if err := validate.Struct(movie); err != nil {
		os.Remove(thumbnailFilePath)
		os.Remove(previewFilePath)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}


	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&movie).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		os.Remove(thumbnailFilePath)
		os.Remove(previewFilePath)
		http.Error(w, "Failed to save movie to the database", http.StatusInternalServerError)
		return
	}

	sendJSON(w, http.StatusCreated, movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	id := getMovieId(w, r)
	if id == 0 {
		return
	}

	movie, ok := getMovie(w, id)
	if !ok {
		return
	}

	r.ParseMultipartForm(20 << 20)

	if title := r.FormValue("title"); title != "" {
		movie.Title = title
	}
	if description := r.FormValue("description"); description != "" {
		movie.Description = description
	}
	if videoID := r.FormValue("video_id"); videoID != "" {
		if parsedID, err := strconv.ParseUint(videoID, 10, 32); err == nil {
			movie.VideoID = uint(parsedID)
		}
	}

	newThumbnail := putImageFile(w, r, "storage/imgs/thumbnails/", "thumbnail")
	if newThumbnail != "" {
		os.Remove(movie.Thumbnail)
		movie.Thumbnail = newThumbnail
	}

	newPreview := putImageFile(w, r, "storage/imgs/previews/", "preview")
	if newPreview != "" {
		os.Remove(movie.Preview)
		movie.Preview = newPreview
	}

	if err := validate.Struct(movie); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := database.DB.Save(movie).Error; err != nil {
		http.Error(w, "Failed to update movie", http.StatusInternalServerError)
		return
	}

	sendJSON(w, http.StatusOK, movie)
}
