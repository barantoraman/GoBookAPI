package data

import (
	"math"
	"strings"

	"github.com/barantoraman/GoBookAPI/internal/validator"
)

type Filters struct {
	Cursor       int
	CursorSize   int
	Sort         string
	SortSafeList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Cursor > 0, "cursor", "must be greater than zero")
	v.Check(f.Cursor <= 10_000_000, "cursor", "must be a maximum of 10 million")
	v.Check(f.CursorSize > 0, "cursor_size", "must be greater than zero")
	v.Check(f.CursorSize <= 100, "cursor_size", "must be a maximum of 100")
	// Check that the sort parameter matches a value in the safelist.
	v.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.CursorSize
}

func (f Filters) offset() int {
	return (f.Cursor - 1) * f.CursorSize
}

// for metadata
type Metadata struct {
	CurrentCursor int `json:"current_cursor,omitempty"`
	CursorSize    int `json:"cursor_size,omitempty"`
	FirstCursor   int `json:"first_cursor,omitempty"`
	LastCursor    int `json:"last_cursor,omitempty"`
	TotalRecords  int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, cursor, cursorSize int) Metadata {
	if totalRecords == 0 {
		// If no records are found, we provide an empty Metadata struct.
		return Metadata{}
	}
	return Metadata{
		CurrentCursor: cursor,
		CursorSize:    cursorSize,
		FirstCursor:   1,
		LastCursor:    int(math.Ceil(float64(totalRecords) / float64(cursorSize))),
		TotalRecords:  totalRecords,
	}
}
