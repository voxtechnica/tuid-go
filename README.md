# Time-based Unique Identifiers (TUID)

[![Go Reference](https://pkg.go.dev/badge/github.com/voxtechnica/tuid-go.svg)](https://pkg.go.dev/github.com/voxtechnica/tuid-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/voxtechnica/tuid-go)](https://goreportcard.com/report/github.com/voxtechnica/tuid-go)

A TUID is an immutable, chronologically sortable, unique identifier with an embedded nanosecond-resolution timestamp.
Internally, it is a big integer value, expressed as a case-sensitive base-62 string (digits 0-9, characters A-Z, and
characters a-z). The left bits are the embedded nanoseconds since epoch, and the right 32 bits are entropy (a random
integer). Valid IDs in the 21st century are 16 characters long, and the bit length ranges from 92-94 bits.
Example: `91Mq07yx9IxHCi5Y` is a TUID with a UTC timestamp of `2021-03-08T05:54:09.208207Z`. Because the left bits are
the nanoseconds since epoch, the IDs sort in chronological order.

A TUID compares nicely with a 128-bit
UUID ([Universally Unique Identifier](https://en.wikipedia.org/wiki/Universally_unique_identifier)), but it's not
a [UUID](https://www.ietf.org/rfc/rfc4122.txt). While a Version 1 UUID has an embedded timestamp, it's not
chronologically sortable. A random Version 4 UUID is more common. Example: `05cd1093-c52c-48e5-a343-9e0017454067`. And,
while a UUID is readily available in a Go [uuid package](https://pkg.go.dev/github.com/google/uuid), a TUID has
significant advantages:

* A TUID is sortable chronologically (alphabetically), making it easy to order a collection of things with IDs.
* A TUID has an embedded nanosecond resolution timestamp, making it easy to track when things were created.
* A TUID is shorter than a UUID (16 characters vs. 36 characters), requiring less space in a database or URL.
* A TUID has no hyphens, making it easy to copy and paste into other systems.

A TUID also compares nicely with a FriendlyID ([Friendly Identifier](https://github.com/norman/friendly-id)).
Example: `5wbwf6yUxVBcr48AMbz9cb`. A FriendlyID is a base-62 encoded 128-bit UUID without the hyphens, making it shorter
and easier to work with than a UUID. A FriendlyID has the great advantage of being interoperable with any existing UUID,
but it's missing the embedded timestamp and is not chronologically sortable. Also, a FriendlyID version of a full
128-bit UUID is 22 characters, whereas a TUID is only 16 characters.

So, why use a TUID? Because you like the embedded timestamp, and you like being able to sort things with IDs
chronologically.

## Installation

```bash
go get github.com/voxtechnica/tuid-go
```

## Usage

Creating a new TUID is simple:

```go
package main

import (
	"fmt"
	"time"

	"github.com/voxtechnica/tuid-go"
)

func main() {
	// Cast a string to a TUID
	id := tuid.TUID("91Mq07yx9IxHCi5Y")

	// Create a new TUID with current time and random entropy
	newID := tuid.NewID()
	fmt.Println("New ID:", newID)

	// Create a TUID with a specific timestamp
	timeID := tuid.NewIDWithTime(time.Now().UTC())

	// Create the first TUID for a specific time (zero entropy, useful for pagination)
	firstID := tuid.FirstIDWithTime(time.Now().UTC())

	// Validate a TUID
	isValid := tuid.IsValid(newID)
	fmt.Println("Is Valid:", isValid)

	// Extract the timestamp
	timestamp, _ := id.Time()
	fmt.Println("Created at:", timestamp)

	// Extract both timestamp and entropy
	info, _ := id.Info()
	fmt.Printf("Info - Time: %v, Entropy: %d\n", info.Timestamp, info.Entropy)

	// Calculate the duration between two IDs
	duration, _ := tuid.Duration(id, newID)
	fmt.Println("Duration between IDs:", duration)
}
```

Since `TUID` is an alias for a `string`, it serializes to and from JSON transparently as a string. The package also provides a `TUIDInfo` struct with JSON tags for convenience when returning detailed ID information from an API.

If you have an ID string, and you want to check if it's a valid TUID, you can use the function `tuid.IsValid(id)`. Valid means that it has only base-62 characters (digits 0-9, letters A-Z, and letters a-z), and that it has an embedded timestamp in the 21st century.

If you have a pair of IDs (e.g. two events), you can compute the `time.Duration` between them with `tuid.Duration(startID, stopID)`.
