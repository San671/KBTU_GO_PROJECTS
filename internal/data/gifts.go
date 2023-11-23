package data

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
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

func (m GiftModel) Get(id int64) (*Gift, error) {
	// The PostgreSQL bigserial type that we're using for the movie ID starts
	// auto-incrementing at 1 by default, so we know that no movies will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	query := `
SELECT id, created_at, title, description, superiority, status, category, version
FROM gifts
WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var gift Gift
	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := m.DB.QueryRow(query, id).Scan(
		&gift.ID,
		&gift.CreatedAt,
		&gift.Title,
		&gift.Description,
		&gift.Superiority,
		&gift.Status,
		&gift.Category,
		&gift.Version,
	)
	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &gift, nil
}

func (m GiftModel) Update(gift *Gift) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
UPDATE gifts
SET title = $1, description = $2, superiority = $3, status = $4, category =$5, version = version + 1
WHERE id = $6
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		gift.Title,
		gift.Description,
		gift.Superiority,
		gift.Status,
		gift.Category,
		gift.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&gift.Version)
}

func (m GiftModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
        DELETE FROM gifts
        WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
