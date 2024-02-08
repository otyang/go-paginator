package paginator

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func _base64Json(next bool, cursorValue string) string {
	b, _ := json.Marshal(&Cursor{DirectionNext: next, Value: cursorValue})
	return base64.StdEncoding.EncodeToString(b)
}

func TestEncodeCursor(t *testing.T) {
	tests := []struct {
		name      string
		cursorOBJ *Cursor
		want      string
	}{
		{"nil cursor", nil, ""},
		{"error json marshalling", (*Cursor)(nil), ""},
		{"valid prev Cursor", &Cursor{DirectionNext: false, Value: "prev"}, _base64Json(false, "prev")},
		{"valid next Cursor", &Cursor{DirectionNext: true, Value: "next"}, _base64Json(true, "next")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeCursor(tt.cursorOBJ)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDecodeCursor(t *testing.T) {
	tests := []struct {
		name                string
		base64EncodedCursor string
		want                *Cursor
		wantErr             bool
	}{
		{"empty string", "", nil, true},
		{"error unmarshalling json", "X_eyJEaXJlY3Rpb25OZXh0IjpmYWxzZSwiQ3Vyc29yVmFsdWUiOiJhcHJldiJ9", nil, true},
		{
			"valid prev cursor", "eyJEaXJlY3Rpb25OZXh0IjpmYWxzZSwiVmFsdWUiOiJwcmV2In0=",
			&Cursor{DirectionNext: false, Value: "prev"}, false,
		},
		{
			"valid next cursor", "eyJEaXJlY3Rpb25OZXh0Ijp0cnVlLCJWYWx1ZSI6Im5leHQifQ==",
			&Cursor{DirectionNext: true, Value: "next"}, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeCursor(tt.base64EncodedCursor)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_getValue(t *testing.T) {
	fmtS := func(t any) string {
		s := fmt.Sprintf("%v", t)
		return s
	}

	t.Run("integer", func(t *testing.T) {
		structure := struct{ ID int }{1}
		got, err := getValue(structure, "ID")
		assert.NoError(t, err)
		assert.Equal(t, fmtS(structure.ID), got)
	})

	// string
	t.Run("string", func(t *testing.T) {
		structure := struct{ Title string }{"title"}
		got, err := getValue(structure, "Title")
		assert.NoError(t, err)
		assert.Equal(t, fmtS(structure.Title), got)
	})

	// time
	t.Run("time", func(t *testing.T) {
		structure := struct{ CreatedAt time.Time }{time.Now()}
		got, err := getValue(structure, "CreatedAt")
		assert.NoError(t, err)
		assert.Equal(t, fmtS(structure.CreatedAt), got)
	})

	// boolean
	t.Run("boolean", func(t *testing.T) {
		structure := struct{ IsActivated bool }{true}
		got, err := getValue(structure, "IsActivated")
		assert.NoError(t, err)
		assert.Equal(t, fmtS(structure.IsActivated), got)
	})

	// panics
	t.Run("error", func(t *testing.T) {
		structure := struct{ ID string }{"_id1234567890"}
		v, err := getValue(structure, "FieldDoesNotExist")
		assert.Error(t, err)
		assert.Empty(t, v)
	})
}
