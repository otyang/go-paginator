package paginator

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
)

type Cursor struct {
	DirectionNext bool
	Value         string
}

// EncodeCursor encodes a Cursor object into a base64-encoded string.
//
// If the provided cursor object is nil, an empty string is returned.
func EncodeCursor(cursorOBJ *Cursor) string {
	if cursorOBJ == nil {
		return ""
	}
	serializedCursor, err := json.Marshal(cursorOBJ)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(serializedCursor)
}

// DecodeCursor decodes a base64-encoded string back into a Cursor object.
func DecodeCursor(base64EncodedCursor string) (*Cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(base64EncodedCursor)
	if err != nil {
		return nil, err
	}

	var unserialisedCursor Cursor
	if err := json.Unmarshal(decodedCursor, &unserialisedCursor); err != nil {
		return nil, err
	}

	return &unserialisedCursor, nil
}

// getValue retrieves the value of a specific field from a provided struct.
//
// This is a helper function used internally by the other functions
// to access the `Value` field from the `Cursor` struct.
//
// It uses reflection to dynamically access the field by name and ensures
// the field exists before returning its value as a string pointer.
//
// In case of any errors (e.g., field not found), it panics with an informative error message.
func getValue(vStruct any, fieldName string) (string, error) {
	v := reflect.ValueOf(vStruct)

	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("GetFieldValue: expected a struct, got %T", v)
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return "", fmt.Errorf("GetFieldValue: field '%s' does not exist", fieldName)
	}

	return fmt.Sprintf("%v", field), nil
}

// Reverse reverses the order of elements in a slice of type T.
// It works for any type of slice, including strings, structs, and other custom types.
//
// Note: Be aware that this modifies the original slice, not a copy.
// If you need to preserve the original order, create a copy before calling Reverse.
//
// Example usage:
//
// numbers := []int{1, 2, 3, 4, 5}
// reversed := Reverse(numbers)
// fmt.Println(reversed) // Output: [5 4 3 2 1]
func Reverse[T any](s []T) []T {
	// Create a new slice to avoid modifying the original.
	reversed := make([]T, len(s))

	// Copy elements in reverse order.
	for i, j := 0, len(s)-1; i <= j; i, j = i+1, j-1 {
		reversed[i] = s[j]
	}

	return reversed
}
