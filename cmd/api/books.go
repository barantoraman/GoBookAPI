package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/barantoraman/GoBookAPI/internal/data"
	"github.com/barantoraman/GoBookAPI/internal/validator"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
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
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
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

	// Call the ValidateBook() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
