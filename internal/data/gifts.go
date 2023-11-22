package data

import (
	"database/sql"
	"personalized_gifts.sanzhar.net/internal/validator"
	"time"
)

// Annotate the Gift struct with struct tags to control how the keys appear in the
// JSON-encoded output.

type Gift struct {
	ID          int64       `json:"id"`
	CreatedAt   time.Time   `json:"-"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Superiority string      `json:"superiority"`
	Preparation Preparation `json:"preparation"`
	Status      string      `json:"status"`
	Category    string      `json:"category"`
}

func ValidateGift(v *validator.Validator, gift *Gift) {
	v.Check(gift.Title != "", "title", "must be provided")
	v.Check(len(gift.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(gift.Description != "", "description", "must be provided")
	v.Check(len(gift.Description) <= 1000, "description", "must not be more than 1000 bytes long")

	v.Check(gift.Superiority != "", "superiority", "must be provided")
	v.Check(gift.Status != "", "status", "must be provided")
	v.Check(gift.Category != "", "category", "must be provided")
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type GiftModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m GiftModel) Insert(gift *Gift) error {
	return nil
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m GiftModel) Get(id int64) (*Gift, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m GiftModel) Update(gift *Gift) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m GiftModel) Delete(id int64) error {
	return nil
}
