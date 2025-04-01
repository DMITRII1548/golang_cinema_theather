package routes

import (
	"api/online-cinema-theather/internal/handlers"
	"net/http"
)

func RegisterMovieRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetMovies(w, r)
		case http.MethodPost:
			handlers.CreateMovie(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetMovie(w, r)
		case http.MethodDelete:
			handlers.DeleteMovie(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}