package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Update the routes() method to return a http.Handler instead of a *httprouter.Router.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/gifts", app.listGiftsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/gifts", app.createGiftHandler)
	router.HandlerFunc(http.MethodGet, "/v1/gifts/:id", app.showGiftHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/gifts/:id", app.updateGiftHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/gifts/:id", app.deleteGiftHandler)
	// Add the route for the POST /v1/users endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	// Add the route for the PUT /v1/users/activated endpoint.
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	return app.recoverPanic(app.rateLimit(router))
}
