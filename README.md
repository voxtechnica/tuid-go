# Time-based Unique Identifiers (TUID)

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

So, why use a TUID? Because you like the embedded timestamp, and you like being able to sort things chronologically.

## Using TUIDs

Creating a new TUID is simple:

* `id := TUID("91Mq07yx9IxHCi5Y")` :: Cast a string to a TUID.
* `id := tuid.NewID()` :: Create a TUID with the current system time and a random entropy value.
* `id := tuid.NewIDWithTime(time.Now())` :: Create a TUID with a specified timestamp and random entropy value.
* `id := tuid.FirstIdWithTime(time.Now())` :: Create the first TUID with a specified timestamp and zero entropy. These
  are useful for paging through a collection of things with IDs.

If you have an ID string, and you want check if it's a valid TUID, you can use the `tuid.IsValid(id)` function. And, for
valid TUIDs, you can extract the embedded timestamp with `id.Time()`. If you have a pair of IDs, you can compute the
`time.Duration` between them with `tuid.Duration(startId, stopId)`.
