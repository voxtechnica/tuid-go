// Package tuid provides facilities for generating and working with Time-based Unique Identifiers (TUID).
package tuid

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// MinID is the first ID at 2000-01-01T00:00:00Z
const MinID = "5Hr02eJHAfTt1tTM"

// MaxID is the first ID at 2100-01-01T00:00:00Z
const MaxID = "MuklDY5bgW1s9Ev2"

// TUID is a Time-based Unique Identifier (e.g. 91Mq07yx9IxHCi5Y) that has an embedded timestamp and sorts
// chronologically. It's a 16-digit base-62 big integer, where the leftmost bits are a timestamp with nanosecond
// resolution (e.g. 2021-03-08T05:54:09.208207000Z) and the rightmost 32 bits are "entropy" (a random number),
// providing some assurance of uniqueness if multiple IDs are created at the same moment. Collisions in a single
// information system are extremely unlikely. The Zero value of a TUID is an empty string.
type TUID string

// TUIDInfo is a convenience type for parsing a TUID's timestamp and entropy.
type TUIDInfo struct {
	ID        TUID      `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Entropy   uint32    `json:"entropy"`
}

// Int decodes the specified base-62 encoded TUID into a big integer
func (t TUID) Int() (*big.Int, error) {
	id, err := decode(string(t))
	if err != nil {
		return new(big.Int), err
	}
	return id, nil
}

// Time extracts the embedded timestamp from the specified TUID
func (t TUID) Time() (time.Time, error) {
	id, err := decode(string(t))
	if err != nil {
		return time.Time{}, err
	}
	nsec := new(big.Int).Rsh(id, 32)
	return time.Unix(0, nsec.Int64()), nil
}

// Entropy extracts the random 32 bits from the specified TUID
func (t TUID) Entropy() (uint32, error) {
	id, err := decode(string(t))
	if err != nil {
		return 0, err
	}
	mask := big.NewInt(1<<32 - 1)
	entropy := new(big.Int).And(id, mask)
	return uint32(entropy.Int64()), nil
}

// Info extracts the timestamp and entropy from the specified TUID
func (t TUID) Info() (TUIDInfo, error) {
	id, err := decode(string(t))
	if err != nil {
		return TUIDInfo{}, err
	}
	nsec := new(big.Int).Rsh(id, 32)
	timestamp := time.Unix(0, nsec.Int64())
	mask := big.NewInt(1<<32 - 1)
	entropy := uint32(new(big.Int).And(id, mask).Int64())
	return TUIDInfo{t, timestamp, entropy}, nil
}

// String implements the fmt.Stringer interface
func (t TUID) String() string {
	return string(t)
}

// NewID creates a new TUID with the current system time
func NewID() TUID {
	return NewIDWithTime(time.Now())
}

// NewIDWithTime creates a TUID with the provided timestamp
func NewIDWithTime(t time.Time) TUID {
	ts := new(big.Int).Lsh(big.NewInt(t.UnixNano()), 32)
	entropy, _ := rand.Int(rand.Reader, big.NewInt(1<<32))
	id := ts.Or(ts, entropy)
	tuid, _ := encode(id)
	return TUID(tuid)
}

// FirstIDWithTime creates a TUID with the provided timestamp and zero entropy, useful for query offsets
func FirstIDWithTime(t time.Time) TUID {
	id := new(big.Int).Lsh(big.NewInt(t.UnixNano()), 32)
	tuid, _ := encode(id)
	return TUID(tuid)
}

// IsValid checks to see if the provided TUID has valid characters and a reasonable embedded timestamp
func IsValid(t TUID) bool {
	id, err := decode(string(t))
	if err != nil {
		return false
	}
	minID, _ := decode(MinID)
	maxID, _ := decode(MaxID)
	return (id.Cmp(minID) >= 0) && (id.Cmp(maxID) <= 0)
}

// Compare supports sorting TUIDs chronologically
func Compare(t1 TUID, t2 TUID) int {
	if t1 == t2 {
		return 0
	}
	if t1 < t2 {
		return -1
	}
	return +1
}

// Duration returns the number of nanoseconds between two TUIDs as a time.Duration
func Duration(start TUID, stop TUID) (time.Duration, error) {
	startTime, err := start.Time()
	if err != nil {
		return 0, err
	}
	stopTime, err := stop.Time()
	if err != nil {
		return 0, err
	}
	return stopTime.Sub(startTime), nil
}

var base = big.NewInt(62)
var digits = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// encode the provided big integer into a base-62 encoded string
func encode(value *big.Int) (string, error) {
	if value.Sign() < 0 {
		return "", errors.New("base 62 encoding error: positive value required")
	}
	var result []byte
	for value.Sign() > 0 {
		q, r := new(big.Int).DivMod(value, base, new(big.Int))
		d := digits[r.Int64()]
		result = append([]byte{d}, result...) // prepend the new digit
		value = q
	}
	if len(result) == 0 {
		return string(digits[0]), nil
	}
	return string(result), nil
}

// decode the provided base-62 encoded string into a big integer
func decode(text string) (*big.Int, error) {
	textBytes := []byte(text)
	size := len(textBytes)
	if size == 0 {
		return new(big.Int), errors.New("base 62 decoding error: no digits")
	}
	result := new(big.Int)
	for i := 0; i < size; i++ {
		b := textBytes[size-1-i] // examine digits from right to left
		j := int64(bytes.IndexByte(digits, b))
		if j == -1 {
			msg := fmt.Sprintf("base 62 decoding error: invalid digit `%s` in %s", string(b), string(textBytes))
			return new(big.Int), errors.New(msg)
		}
		pow := new(big.Int).Exp(base, big.NewInt(int64(i)), nil)
		prod := new(big.Int).Mul(big.NewInt(j), pow)
		result = new(big.Int).Add(result, prod)
	}
	return result, nil
}
