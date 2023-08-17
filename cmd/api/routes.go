package main

import (
	"expvar"
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

	//for books
	router.HandlerFunc(http.MethodGet, "/v1/books", app.requireActivatedUser(app.listBooksHandler))
	router.HandlerFunc(http.MethodPost, "/v1/books", app.requireActivatedUser(app.createBookHandler))
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", app.requireActivatedUser(app.showBookHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", app.requireActivatedUser(app.updateBookHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", app.requireActivatedUser(app.deleteBookHandler))

	// for users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	// for authentication
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// for our custom metrics,
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// Wrap the router with the middlewares
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
