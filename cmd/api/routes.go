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
	// Wrap the router with the panic recovery middleware.
	return app.recoverPanic(router)
}
