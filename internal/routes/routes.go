package routes

import "net/http"

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.StripPrefix("/storage/", http.FileServer(http.Dir("storage")))
	mux.Handle("/storage/", fileServer)

	RegisterMovieRoutes(mux)

	return mux
}