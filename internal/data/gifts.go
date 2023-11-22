package data

import (
	"database/sql"
	"personalized_gifts.sanzhar.net/internal/validator"
	"time"
)

// Annotate the Gift struct with struct tags to control how the keys appear in the
// JSON-encoded output.

type Gift struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Superiority string    `json:"superiority"`
	Status      string    `json:"status"`
	Category    string    `json:"category"`
	Version     int32     `json:"version"`
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

// The Insert() method accepts a pointer to a movie struct, which should contain the
// data for the new record.
func (m GiftModel) Insert(gift *Gift) error {
	// Define the SQL query for inserting a new record in the gifts table and returning
	// the system-generated data.
	query := `
        INSERT INTO gifts (title, description, superiority, status, category)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, version`
	// Create an args slice containing the values for the placeholder parameters from
	// the gift struct.
	args := []interface{}{gift.Title, gift.Description, gift.Superiority, gift.Status, gift.Category}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at, and version values into the gift struct.
	return m.DB.QueryRow(query, args...).Scan(&gift.ID, &gift.CreatedAt, &gift.Version)
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
