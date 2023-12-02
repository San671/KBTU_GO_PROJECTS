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
	// Use the requireActivatedUser() middleware on our five /v1/movies** endpoints.
	router.HandlerFunc(http.MethodGet, "/v1/gifts", app.requireActivatedUser(app.listGiftsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/gifts", app.requireActivatedUser(app.createGiftHandler))
	router.HandlerFunc(http.MethodGet, "/v1/gifts/:id", app.requireActivatedUser(app.showGiftHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/gifts/:id", app.requireActivatedUser(app.updateGiftHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/gifts/:id", app.requireActivatedUser(app.deleteGiftHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
