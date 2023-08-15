package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// for errors
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// for healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// for books
	router.HandlerFunc(http.MethodPost, "/v1/books", app.createBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", app.showBookHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", app.updateBookHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", app.deleteBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books", app.listBooksHandler)

	// for users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	// Wrap the router with the middlewares
	return app.recoverPanic(app.rateLimit(router))
}
