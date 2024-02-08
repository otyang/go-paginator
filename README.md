# go-paginator

This package provides `PagerInfo` and related functionalities for efficient data pagination with cursors in Go applications.

**Key Features:**

* Represents a paginated response with details about the overall data and pagination parameters.
* Supports various data types through its generic `T` parameter.
* Allows extracting string values of cursors for direct use in API responses/requests.
* Provides `NewPagerInfo` function for constructing paginated responses from data slices, limits, and cursor information.
* Handles cursor calculations for both previous and next pages based on the limit and provided cursor.
* Encodes and exposes cursor values in base64 format for easy storage and transmission.

**Benefits:**

* Simplifies implementing efficient data pagination with cursors.
* Ensures consistent and flexible handling of pagination logic for different data types.
* Provides a clean and reusable approach for managing pagination in Go applications.

**Usage:**

1. Import the `paginator` package in your Go code.
2. Define the type of your data elements (`T`).
3. Use `NewPagerInfo` to create a `PagerInfo[T]` object:
    * Pass the data slice, limit, cursor column name, and cursor value (if any).
    * The function calculates total items, cursor values, and returns a populated `PagerInfo[T]` object.
4. Access various properties of the `PagerInfo[T]` object:
    * `Total`: Total number of data items.
    * `Limit`: Requested limit for returned results.
    * `PrevCursor`, `NextCursor`: Cursor pointers for previous and next pages (optional).
    * `Results`: Actual data slice for the current page.
    * `EncodedPrevCursor`, `EncodedNextCursor`: Base64-encoded strings of `PrevCursor` and `NextCursor` (optional).
5. Extract cursor strings using `GetCursors()` method for direct use in API responses or requests.

**Example:**

```go
package main

import (
	"fmt"

	"github.com/otyang/go-paginator"
)

type User struct {
	// ... user fields
}

func main() {
	// Sample data and parameters
	users := []User{{}, {}, {}, {}}
	limit := 2

	// Create paginated response
	pagerInfo, err := paginator.NewPagerInfo(users, limit, "id", "cursorValue")
	if err != nil {
		panic(err)
	}

	// Access and utilize pagination information
	fmt.Println("Total:", pagerInfo.Total)
	fmt.Println("Limit:", pagerInfo.Limit)
	fmt.Println("PrevCursor:", pagerInfo.PrevCursor)
	fmt.Println("NextCursor:", pagerInfo.NextCursor)
	fmt.Println("Results:", pagerInfo.Results)
	fmt.Println("EncodedPrevCursor:", pagerInfo.EncodedPrevCursor)
	fmt.Println("EncodedNextCursor:", pagerInfo.EncodedNextCursor)

	// Get cursor strings for API usage
	prevCursor, nextCursor := pagerInfo.GetCursors()
	fmt.Println("PrevCursor:", prevCursor)
	fmt.Println("NextCursor:", nextCursor)
}
```

This package provides a powerful and versatile solution for managing pagination with cursors in Go applications, making it easier to handle and present large datasets efficiently. 