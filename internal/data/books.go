package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/barantoraman/GoBookAPI/internal/validator"
	"github.com/lib/pq"
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
	query := `
		INSERT INTO books (isbn, title, author, genres, pages, language , publisher, year)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, version`

	args := []interface{}{book.ISBN, book.Title, book.Author, pq.Array(book.Genres), book.Pages, book.Language, book.Publisher, book.Year}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return b.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

func (b BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, created_at, isbn, title, author, genres, pages, language, publisher, year, version
	FROM books
	WHERE id = $1`

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.ISBN,
		&book.Title,
		&book.Author,
		pq.Array(&book.Genres),
		&book.Pages,
		&book.Language,
		&book.Publisher,
		&book.Year,
		&book.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &book, nil
}

func (b BookModel) Update(book *Book) error {
	query := `
		UPDATE books
		SET isbn = $1, title = $2, author = $3, genres = $4, pages = $5, language = $6, publisher = $7 , year = $8, version = version + 1 
		WHERE id = $9 AND version = $10
		RETURNING version`

	args := []interface{}{
		book.ISBN,
		book.Title,
		book.Author,
		pq.Array(book.Genres),
		book.Pages,
		book.Language,
		book.Publisher,
		book.Year,
		book.ID,
		book.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version)
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

func (b BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM books
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := b.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (b BookModel) GetAll(isbn string, title string, author string, genres []string, filters Filters) ([]*Book, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, isbn, title, author, genres, pages, language, publisher, year, version
		FROM books
		WHERE (isbn = $1 OR $1 = '')
		AND (to_tsvector('simple',title) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (to_tsvector('simple',author) @@ plainto_tsquery('simple', $3) OR $3 = '')
		AND (genres @> $4 OR $4 = '{}')
		ORDER BY %s %s , id ASC
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{isbn, title, author, pq.Array(genres), filters.limit(), filters.offset()}

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := b.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err //update this to return an empty metadata struct
	}
	defer rows.Close()

	totalRecords := 0
	//Initialize an empty slice to hold the book data
	books := []*Book{}

	//Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&totalRecords,
			&book.ID,
			&book.CreatedAt,
			&book.ISBN,
			&book.Title,
			&book.Author,
			pq.Array(&book.Genres),
			&book.Pages,
			&book.Language,
			&book.Publisher,
			&book.Year,
			&book.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		books = append(books, &book)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Cursor, filters.CursorSize)
	return books, metadata, nil
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
