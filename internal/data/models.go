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
	Gifts  GiftModel
	Tokens TokenModel // Add a new Tokens field.
	Users  UserModel  // Add a new Users field.
}

func NewModels(db *sql.DB) Models {
	return Models{
		Gifts:  GiftModel{DB: db},
		Tokens: TokenModel{DB: db}, // Initialize a new TokenModel instance.
		Users:  UserModel{DB: db},  // Initialize a new UserModel instance.
	}
}
