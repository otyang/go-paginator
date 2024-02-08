package paginator

// PagerInfo[T any] represents a paginated response containing a subset of data
// along with information about the overall data and pagination parameters.
//
// Type Parameters:
//
//	T: The type of data elements contained in the 'Results' slice.
//
// Fields:
//
//	Total: The total number of data items (before pagination).
//	Limit: The maximum number of data items to be included in a single page.
//	PrevCursor: A pointer to the cursor value for the previous page, or nil if there is no previous page.
//	NextCursor: A pointer to the cursor value for the next page, or nil if there is no next page.
//	Results: A slice containing the paginated data items for the current page.
//	EncodedPrevCursor: The base64-encoded string representation of the 'PrevCursor', or an empty string if 'PrevCursor' is nil.
//	EncodedNextCursor: The base64-encoded string representation of the 'NextCursor', or an empty string if 'NextCursor' is nil.
type PagerInfo[T any] struct {
	Total             int
	Limit             int
	PrevCursor        *Cursor
	NextCursor        *Cursor
	Results           []T
	EncodedPrevCursor string
	EncodedNextCursor string
}

// Extracts and returns the string values of the PrevCursor and NextCursor (if not nil).
// Useful for directly sending the cursor values in API responses or requests.
func (p *PagerInfo[T]) GetCursors() (prev string, next string) {
	if p.PrevCursor != nil {
		prev = p.PrevCursor.Value
	}
	if p.NextCursor != nil {
		next = p.NextCursor.Value
	}
	return prev, next
}

// NewPagerInfo constructs a PagerInfo object from a slice of records,
// limit, and cursor information.
//
// It provides a generic way to handle pagination logic for
// various data types (specified by T).
//
// The PagerInfo object returns the following:
//
//   - Total: The total number of items in the data set.
//   - Limit: The requested limit for returned results.
//   - PrevCursor: The cursor for previous page (optional).
//   - NextCursor: The cursor for next page (optional).
//   - Results: The actual slice of data to be returned.
//   - EncodedPrevCursor: The base64-encoded string representation of the PrevCursor (optional).
//   - EncodedNextCursor: The base64-encoded string representation of the NextCursor (optional).
func NewPagerInfo[T any](records []T, limit int, cursorColumn string, cursorValue string) (PagerInfo[T], error) {
	total := len(records)
	results := records

	// Calculate next cursor
	var nextCursor *Cursor
	{
		// You only get a next cursor when there is more results
		if len(records) > limit {
			total = limit
			results = records[:limit]

			v, err := getValue(records[limit], cursorColumn)
			if err != nil {
				return PagerInfo[T]{}, err
			}

			nextCursor = &Cursor{DirectionNext: true, Value: v}
		}
	}

	// Calculate previous cursor
	var prevCursor *Cursor
	{
		// You only get a prev cursor when not on firstpage AND there are results
		if cursorValue != "" && len(records) > 0 {
			v, err := getValue(results[0], cursorColumn)
			if err != nil {
				return PagerInfo[T]{}, err
			}

			prevCursor = &Cursor{DirectionNext: false, Value: v}
		}
	}

	return PagerInfo[T]{
		Total:             total,
		Limit:             limit,
		PrevCursor:        prevCursor,
		NextCursor:        nextCursor,
		Results:           results,
		EncodedPrevCursor: EncodeCursor(prevCursor),
		EncodedNextCursor: EncodeCursor(nextCursor),
	}, nil
}
