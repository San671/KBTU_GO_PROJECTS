package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&gift.ID, &gift.CreatedAt, &gift.Version)
}

func (m GiftModel) Get(id int64) (*Gift, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Define the SQL query for retrieving the movie data.
	query := `
        SELECT pg_sleep(10), id, created_at, title, description, superiority, status, category, version
        FROM gifts
        WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var gift Gift
	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 3-second timeout deadline. Note that we're using the empty context.Background()
	// as the 'parent' context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&[]byte{}, // Add this line.
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
        WHERE id = $6 AND version = $7
        RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		gift.Title,
		gift.Description,
		gift.Superiority,
		gift.Status,
		gift.Category,
		gift.ID,
		gift.Version,
	}
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&gift.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
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
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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

func (m GiftModel) GetAll(title string, filters Filters) ([]*Gift, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, created_at, title, description, superiority, status, category, version
	FROM gifts
	WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
    ORDER BY %s %s, id ASC
    LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()
	totalRecords := 0
	gifts := []*Gift{}

	for rows.Next() {
		var gift Gift

		err := rows.Scan(
			&totalRecords, // Извлекаем общее количество записей
			&gift.ID,
			&gift.CreatedAt,
			&gift.Title,
			&gift.Description,
			&gift.Superiority,
			&gift.Status,
			&gift.Category,
			&gift.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		gifts = append(gifts, &gift)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return gifts, metadata, nil
}
