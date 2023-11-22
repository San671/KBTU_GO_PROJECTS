package main

import (
	"fmt"
	"net/http"
	"personalized_gifts.sanzhar.net/internal/validator"
	"time"

	"personalized_gifts.sanzhar.net/internal/data"
)

func (app *application) createGiftHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Superiority string `json:"superiority"`
		Status      string `json:"status"`
		Category    string `json:"category"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Movie struct.
	gift := &data.Gift{
		Title:       input.Title,
		Description: input.Description,
		Superiority: input.Superiority,
		Status:      input.Status,
		Category:    input.Category,
	}
	// Initialize a new Validator.
	v := validator.New()
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateGift(v, gift); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Call the Insert() method on our movies model, passing in a pointer to the
	// validated movie struct. This will create a record in the database and update the
	// movie struct with the system-generated information.
	err = app.models.Gifts.Insert(gift)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/gifts/%d", gift.ID))
	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"gift": gift}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showGiftHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		// Use the new notFoundResponse() helper.
		app.notFoundResponse(w, r)
		return
	}

	gift := data.Gift{
		ID:          id,
		CreatedAt:   time.Now(),
		Title:       "Ring ",
		Description: "This ring will make a wonderful gift",
		Superiority: "gold",
		Status:      "ready",
		Category:    "decoration",
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"gift": gift}, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
