package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Pages int32

var ErrInvalidRuntimeFormat = errors.New("invalid pages format")

func (p Pages) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d pages", p)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}

// Implement a UnmarshalJSON() method on the Runtime type so that it satisfies the
// json.Unmarshaler interface.
func (p *Pages) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedJSONValue, " ")
	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidRuntimeFormat error again.
	if len(parts) != 2 || parts[1] != "pages" {
		return ErrInvalidRuntimeFormat
	}
	// Otherwise, parse the string containing the number into an int32. Again, if this
	// fails return the ErrInvalidRuntimeFormat error.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*p = Pages(i)
	return nil
}
