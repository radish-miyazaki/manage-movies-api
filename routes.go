package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.getMovie)

	// Create & Update HandleFunc
	router.HandlerFunc(http.MethodPost, "/v1/admin/movie/edit", app.editMovie)
	router.HandlerFunc(http.MethodDelete, "/v1/admin/movie/delete/:id", app.deleteMovie)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodGet, "/v1/genres/:id", app.getAllMoviesByGenre)

	return app.enableCORS(router)
}
