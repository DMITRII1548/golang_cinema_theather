package routes

import "net/http"

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	RegisterMovieRoutes(mux)

	return mux
}