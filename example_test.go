package tuid_test

import (
	"fmt"
	"time"

	"github.com/voxtechnica/tuid-go"
)

func ExampleNewID() {
	id := tuid.NewID()
	
	// A TUID is always 16 characters long
	fmt.Println(len(id))
	// Output: 16
}

func ExampleFirstIDWithTime() {
	// Create a specific time for our example
	t := time.Date(2021, 3, 8, 5, 54, 9, 208207000, time.UTC)
	
	// Create the first TUID for this exact nanosecond (zero entropy)
	id := tuid.FirstIDWithTime(t)
	fmt.Println(id)
	// Output: 91Mq07yx9IvZeYJE
}

func ExampleTUID_Time() {
	id := tuid.TUID("91Mq07yx9IxHCi5Y")
	
	// Extract the embedded timestamp
	ts, _ := id.Time()
	
	// Format the timestamp
	fmt.Println(ts.Format(time.RFC3339Nano))
	// Output: 2021-03-08T05:54:09.208207Z
}

func ExampleTUID_Entropy() {
	id := tuid.TUID("9AxhASm3k3MVH8se")
	
	// Extract the entropy (random bits)
	entropy, _ := id.Entropy()
	
	fmt.Println(entropy)
	// Output: 1588339652
}

func ExampleDuration() {
	startID := tuid.TUID("9AxgffrWr9qCnfIT")
	stopID := tuid.TUID("9AxjEL0lPtoGAbLE")
	
	// Calculate duration between two TUIDs
	duration, _ := tuid.Duration(startID, stopID)
	
	fmt.Println(duration)
	// Output: 32m2.322640879s
}

func ExampleIsValid() {
	validID := tuid.TUID("91Mq07yx9IxHCi5Y")
	invalidID := tuid.TUID("not-a-valid-tuid")
	
	fmt.Println(tuid.IsValid(validID))
	fmt.Println(tuid.IsValid(invalidID))
	// Output:
	// true
	// false
}
