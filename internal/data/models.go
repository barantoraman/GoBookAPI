package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Books BookModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books: BookModel{DB: db},
	}
}
