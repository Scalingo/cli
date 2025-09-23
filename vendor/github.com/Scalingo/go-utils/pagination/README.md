## Package `pagination` v1.1.1

This is a pagination library for Go.

This library provides a data structure to transport paginated data. It is useful when you need to return a page of items
along with the total number of items and other Metadata.

**What this library does:**

* it provides a data structure to transport paginated data.
* it provides functions to create "paginated data" and "page request" structures.
* it accepts data as a slice of any type. The data supplied must be 1 page of items.

**What this library does NOT do:**

* it is NOT responsible for fetching the data from the database.
* it is NOT responsible for splitting a result set into pages.

The `Request` struct is used to define the page number and the number of items per page. Used when requesting a page
of items.

The `Paginated` struct is used to return a page of items along with the total number of items.

// Example data returned

```json 
{
  "data": [],
  "meta": {
    "current_page": 1,
    "next_page": 1,
    "per_page": 20,
    "prev_page": 1,
    "total_count": 0,
    "total_pages": 1
  }
}
```

### Usage

#### Example of the first page of data of a small dataset:

In this example, we create a new `Paginated` struct with a slice of `Item` and a `Request` struct.

* The total number of items in the database is 4.
* The `Request` struct specifies that we want the first page with 2 items per page.
* The `[]Item` slice contains 2 items (the first page of data).

```go
package main

import (
	"fmt"
	"github.com/Scalingo/go-utils/pagination"
)

type Item struct {
	ID   int
	Name string
}

func main() {
	items := []Item{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
	}

	page := 1
	perPage := 2
	p := pagination.New[[]Item](items, pagination.NewRequest(page, perPage), 4)
	fmt.Println(p)
}
```

#### An example of the second page of data of a larger dataset:

In this example:

* The total number of items in the database is 123.
* The `PageRequest` struct specifies that we want the second page with 10 items per page.
* The `[]Item` slice contains 10 items (the second page of data).

```go
    items := []Item{
        {ID: 21, Name: "Item 21"},
        {ID: 22, Name: "Item 22"},
        {ID: 23, Name: "Item 23"},
        {ID: 24, Name: "Item 24"},
        ...
        {ID: 29, Name: "Item 29"},
        {ID: 30, Name: "Item 30"},
    }
        
    page := 2
	perPage := 10
    p := pagination.New[[]Item](items, pagination.NewRequest(page, perPage), 123)
```

In both examples we rely on the database query to return the correct page of data. The `Paginated` struct is used to
transport the data to the client.

