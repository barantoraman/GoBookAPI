package main

import (
	"context"
	"net/http"

	"github.com/barantoraman/GoBookAPI/internal/data"
)

type contextKey string

// Transform the string "user" into contextKey type and assign it as userContextKey
// constant for managing user information within the request context.
const userContextKey = contextKey("user")

// contextSetUser() creates a fresh request copy embedded with the given User struct
// within the context, utilizing the userContextKey constant for referencing.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser() retrieves the User struct from the request context.
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
