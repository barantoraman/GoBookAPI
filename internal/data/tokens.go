package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/barantoraman/GoBookAPI/internal/validator"
)

// Define a Token struct to hold the data for an individual token.
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

// Define the TokenModel type.
type TokenModel struct {
	DB *sql.DB
}

// The New() method efficiently generates a fresh Token struct and seamlessly
// populates the tokens table with the associated data.
func (t TokenModel) New(userID int64, ttl time.Duration) (*Token, error) {
	token, err := generateToken(userID, ttl)
	if err != nil {
		return nil, err
	}
	err = t.Insert(token)
	return token, err
}

// Insert method adds the provided token data into the tokens table.
func (t TokenModel) Insert(token *Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry) 
	VALUES ($1, $2, $3)`
	args := []interface{}{token.Hash, token.UserID, token.Expiry}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := t.DB.ExecContext(ctx, query, args...)
	return err
}

// DeleteAllForUser method removes all tokens associated with
// the given userID from the tokens table.
func (t TokenModel) DeleteAllForUser(userID int64) error {
	query := `
	DELETE FROM tokens
	WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := t.DB.ExecContext(ctx, query, userID)
	return err
}

func generateToken(userID int64, ttl time.Duration) (*Token, error) {
	// Create a Token instance containing the user ID and expiry information.
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}
	// Fill the byte slice using the Read() function from the crypto/rand package,
	// which retrieves random bytes from the operating system's CSPRNG. An error
	// will be returned in case of CSPRNG malfunction.
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generate a SHA-256 hash of the plaintext token string.
	// we convert it to a slice using the [:] operator before storing it.
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}
