package data

import (
	"database/sql"
	"time"

	"github.com/barantoraman/GoBookAPI/internal/validator"
)

type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	ISBN      string    `json:"isbn"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Genres    []string  `json:"genres"`
	Pages     Pages     `json:"pages"`
	Language  string    `json:"language"`
	Publisher string    `json:"publisher"`
	Year      int32     `json:"year"`
	Version   int32     `json:"version"`
}

type BookModel struct {
	DB *sql.DB
}

func (b BookModel) Insert(book *Book) error {
	return nil
}

func (b BookModel) Get(id int64) (*Book, error) {
	return nil, nil
}

func (b BookModel) Update(book *Book) error {
	return nil
}

func (b BookModel) Delete(id int64) error {
	return nil
}

// ValidateBook checks the validity of the provided book data
func ValidateBook(v *validator.Validator, book *Book) {
	// ISBN validation
	v.Check(book.ISBN != "", "isbn", "must be provided")
	v.Check(len(book.ISBN) == 13, "isbn", "must be 13 digit")

	// Title validation
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 500, "title", "must not be more than 500 bytes long")

	// Author validation
	v.Check(book.Author != "", "author", "must be provided")
	v.Check(len(book.Author) <= 500, "author", "must not be more than 500 bytes long")

	// Genres validation
	v.Check(book.Genres != nil, "genres", "must be provided")
	v.Check(len(book.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(book.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(book.Genres), "genres", "must not contain duplicate values")

	// Pages validation
	v.Check(book.Pages != 0, "pages", "must be provided")
	v.Check(book.Pages > 0, "pages", "must be a positive integer")

	// Language validation
	v.Check(book.Language != "", "language", "must be provided")
	v.Check(len(book.Language) <= 500, "language", "must not be more than 500 bytes long")

	//Publisher validation
	v.Check(book.Publisher != "", "publisher", "must be provided")
	v.Check(len(book.Publisher) <= 500, "publisher", "must not be more than 500 bytes long")

	// Year validation
	v.Check(book.Year != 0, "year", "must be provided")
	v.Check(book.Year >= 0, "year", "must be greater than 0")
	v.Check(book.Year <= int32(time.Now().Year()), "year", "must not be in the future")
}
