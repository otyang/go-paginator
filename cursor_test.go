package paginator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Books struct {
	ID    int
	Title string
}

var allBooks = []Books{
	{ID: 1, Title: "48 Laws of power -1"},
	{ID: 2, Title: "48 Laws of power -2"},
	{ID: 3, Title: "48 Laws of power -3"},
	{ID: 4, Title: "48 Laws of power -4"},
	{ID: 5, Title: "48 Laws of power -5"},
	{ID: 6, Title: "48 Laws of power -6"},
	{ID: 7, Title: "48 Laws of power -7"},
}

func testGetCursorID(next bool, index int) *Cursor {
	return &Cursor{next, fmt.Sprintf("%d", allBooks[index].ID)}
}

func TestNewPagerInfo(t *testing.T) {
	t.Run("expecting error: invalid cursor field.", func(t *testing.T) {
		_, err := NewPagerInfo[Books](allBooks, 4, "FieldDoesNotExist", "")
		assert.Error(t, err)
	})

	t.Run("empty or nil records", func(t *testing.T) {
		cursorValue := ""
		limit := 4

		cursor, err := NewPagerInfo[Books](nil, limit, "ID", cursorValue)
		assert.NoError(t, err)
		assert.Equal(
			t, PagerInfo[Books]{
				Total:             0,
				Limit:             limit,
				PrevCursor:        nil,
				NextCursor:        nil,
				Results:           nil,
				EncodedPrevCursor: "",
				EncodedNextCursor: "",
			},
			cursor)

		cursor, err = NewPagerInfo[Books]([]Books{}, limit, "ID", cursorValue)
		assert.NoError(t, err)
		assert.Equal(t, PagerInfo[Books]{Limit: limit, Results: []Books{}}, cursor)
	})

	t.Run("length of records below or equal the limit", func(t *testing.T) {
		limit := 1000
		cursorValue := ""

		wantPagerInfo := PagerInfo[Books]{
			Total:             len(allBooks),
			Limit:             limit,
			PrevCursor:        nil,
			NextCursor:        nil,
			Results:           allBooks,
			EncodedPrevCursor: "",
			EncodedNextCursor: "",
		}

		gotPagerInfo, err := NewPagerInfo[Books](allBooks, limit, "ID", cursorValue)
		assert.NoError(t, err)
		assert.Equal(t, wantPagerInfo, gotPagerInfo)
	})

	t.Run("length record above the limit and on first page", func(t *testing.T) {
		var (
			limit           = 3
			indexNextCursor = limit
			cursorValue     = ""
		)

		wantPagerInfo := PagerInfo[Books]{
			Total:             limit,
			Limit:             limit,
			PrevCursor:        nil,
			NextCursor:        testGetCursorID(true, indexNextCursor),
			Results:           allBooks[:limit],
			EncodedPrevCursor: "",
			EncodedNextCursor: EncodeCursor(testGetCursorID(true, indexNextCursor)),
		}

		gotPagerInfo, err := NewPagerInfo[Books](allBooks, limit, "ID", cursorValue)
		assert.NoError(t, err)
		assert.Equal(t, wantPagerInfo, gotPagerInfo)

		// debugpurpose
		// p, n := wantPagerInfo.GetCursors()
		// t.Error("===", p, n)
	})

	t.Run("length record above the limit and not on first page", func(t *testing.T) {
		var (
			limit           = 3
			indexNextCursor = limit
			cursorValue     = "1"
		)

		wantPagerInfo := PagerInfo[Books]{
			Total:             limit,
			Limit:             limit,
			PrevCursor:        testGetCursorID(false, 0),
			NextCursor:        testGetCursorID(true, indexNextCursor),
			Results:           allBooks[:limit],
			EncodedPrevCursor: EncodeCursor(testGetCursorID(false, 0)),
			EncodedNextCursor: EncodeCursor(testGetCursorID(true, indexNextCursor)),
		}

		gotPagerInfo, err := NewPagerInfo[Books](allBooks, limit, "ID", cursorValue)
		assert.NoError(t, err)
		assert.Equal(t, wantPagerInfo, gotPagerInfo)
	})
}
