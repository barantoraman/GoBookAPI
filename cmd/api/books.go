package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/barantoraman/GoBookAPI/internal/data"
	"github.com/barantoraman/GoBookAPI/internal/validator"
)

// createBookHandler handles the HTTP POST request to create a new book.
func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the input data from the JSON request.
	var input struct {
		ISBN      string     `json:"isbn"`
		Title     string     `json:"title"`
		Author    string     `json:"author"`
		Genres    []string   `json:"genres"`
		Pages     data.Pages `json:"pages"`
		Language  string     `json:"language"`
		Publisher string     `json:"publisher"`
		Year      int32      `json:"year"`
	}

	// Read and parse the JSON request body into the input struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Create a new Book instance using the input data.
	book := &data.Book{
		ISBN:      input.ISBN,
		Title:     input.Title,
		Author:    input.Author,
		Genres:    input.Genres,
		Pages:     input.Pages,
		Language:  input.Language,
		Publisher: input.Publisher,
		Year:      input.Year,
	}

	v := validator.New()

	// Validate the book data using the ValidateBook() function and return errors if any.
	if data.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the validated book into the database.
	err = app.models.Books.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Set the 'Location' header for the newly created book resource.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	// Write the JSON response with the created book data and headers.
	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showBookHandler handles the HTTP GET request to retrieve a book.
func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the book ID from the request.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Retrieve the book from the database.
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			// Respond with "Not Found" if book doesn't exist.
		default:
			app.serverErrorResponse(w, r, err)
			// Respond with server error for other errors.
		}
		return
	}

	// Respond with the book data in JSON format.
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		// Respond with server error if JSON writing fails.
		app.serverErrorResponse(w, r, err)
	}
}

// updateBookHandler handles the HTTP PATCH request to update a book's details.
func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the book ID from the request.
	id, err := app.readIDParam(r)
	if err != nil {
		// Respond with "Not Found" if ID is invalid.
		app.notFoundResponse(w, r)
		return
	}

	// Retrieve the existing book record from the database
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case
			errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Define a struct to hold the updated fields received from the JSON request
	var input struct {
		ISBN      *string     `json:"isbn"`
		Title     *string     `json:"title"`
		Author    *string     `json:"author"`
		Genres    []string    `json:"genres"`
		Pages     *data.Pages `json:"pages"`
		Language  *string     `json:"language"`
		Publisher *string     `json:"publisher"`
		Year      *int32      `json:"year"`
	}

	// Decode the JSON request body and update the book fields if values are provided.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Update the book fields with the new values if provided.
	if input.ISBN != nil {
		book.ISBN = *input.ISBN
	}
	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Author != nil {
		book.Author = *input.Author
	}
	if input.Genres != nil {
		book.Genres = input.Genres
		// we don't need to dereference a slice.
	}
	if input.Pages != nil {
		book.Pages = *input.Pages
	}
	if input.Language != nil {
		book.Language = *input.Language
	}
	if input.Publisher != nil {
		book.Publisher = *input.Publisher
	}

	if input.Year != nil {
		book.Year = *input.Year
	}

	v := validator.New()

	if data.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the book record in the database.
	err = app.models.Books.Update(book)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Respond with the updated book data in JSON format.
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteBookHandler handles the HTTP DELETE request to remove a book by its ID.
func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the book ID from the request.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Delete the book record from the database.
	err = app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Respond with a success message in JSON format.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listBooksHandler retrieves a list of books based on query parameters,
// applies filters and sorting, and sends a JSON response with the results.
func (app *application) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ISBN   string
		Title  string
		Author string
		Genres []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.ISBN = app.readString(qs, "isbn", "")
	input.Title = app.readString(qs, "title", "")
	input.Author = app.readString(qs, "author", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Cursor = app.readInt(qs, "cursor", 1, v)
	input.Filters.CursorSize = app.readInt(qs, "cursor_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "author", "year", "-id", "-title", "-author", "-year"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	books, metadata, err := app.models.Books.GetAll(input.ISBN, input.Title, input.Author, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"books": books, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
