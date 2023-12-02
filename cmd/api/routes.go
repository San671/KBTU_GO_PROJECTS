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
	// Use the requirePermission() middleware on each of the /v1/movies** endpoints,
	// passing in the required permission code as the first parameter.
	router.HandlerFunc(http.MethodGet, "/v1/gifts", app.requirePermission("gifts:read", app.listGiftsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/gifts", app.requirePermission("gifts:write", app.createGiftHandler))
	router.HandlerFunc(http.MethodGet, "/v1/gifts/:id", app.requirePermission("gifts:read", app.showGiftHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/gifts/:id", app.requirePermission("gifts:write", app.updateGiftHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/gifts/:id", app.requirePermission("gifts:write", app.deleteGiftHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	// Add the enableCORS() middleware.
	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
