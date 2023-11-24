package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	// Set the Movies field to be an interface containing the methods that both the
	// 'real' model and mock model need to support.
	Gifts interface {
		Insert(gift *Gift) error
		Get(id int64) (*Gift, error)
		Update(gift *Gift) error
		Delete(id int64) error
	}
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Gifts: GiftModel{DB: db},
	}
}
