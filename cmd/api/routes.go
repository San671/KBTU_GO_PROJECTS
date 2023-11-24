package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// Add the route for the GET /v1/movies endpoint.
	router.HandlerFunc(http.MethodGet, "/v1/gifts", app.listGiftsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/gifts", app.createGiftHandler)
	router.HandlerFunc(http.MethodGet, "/v1/gifts/:id", app.showGiftHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/gifts/:id", app.updateGiftHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/gifts/:id", app.deleteGiftHandler)
	return router
}
