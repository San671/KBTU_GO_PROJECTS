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
	movie := &data.Gift{
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
	if data.ValidateGift(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
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
