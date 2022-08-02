package tuid

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"
)

type encodingTest struct {
	name         string
	i            int64
	s            string
	errorMessage string
}

func TestEncode(t *testing.T) {
	data := []encodingTest{
		{"zero", 0, "0", ""},
		{"positive", 1024, "GW", ""},
		{"negative", -1, "", "positive value required"},
		{"boundary", 62, "10", ""},
		{"pooch", 50014, "D0g", ""},
		{"million", 1_000_000, "4C92", ""},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			text, err := encode(big.NewInt(d.i))
			if text != d.s {
				t.Errorf("expected %s, received %s", d.s, text)
			}
			var msg string
			if err != nil {
				msg = err.Error()
			}
			if !strings.Contains(msg, d.errorMessage) {
				t.Errorf("expected error message `%s`, received `%s`", d.errorMessage, msg)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	data := []encodingTest{
		{"empty", 0, "", "no digits"},
		{"invalid", 0, "bad!", "invalid digit"},
		{"spaces", 0, " 10", "invalid digit"},
		{"zero", 0, "0", ""},
		{"positive", 1024, "GW", ""},
		{"boundary", 62, "10", ""},
		{"pooch", 50014, "D0g", ""},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			i, err := decode(d.s)
			if i.Int64() != d.i {
				t.Errorf("expected %d, received %d", d.i, i)
			}
			var msg string
			if err != nil {
				msg = err.Error()
			}
			if !strings.Contains(msg, d.errorMessage) {
				t.Errorf("expected error message `%s`, received `%s`", d.errorMessage, msg)
			}
		})
	}
}

func TestNewTuid(t *testing.T) {
	tuid := NewID()
	if tuid == "" {
		t.Error("expected a TUID, got a zero value")
	}
}

func TestNewTuidWithTime(t *testing.T) {
	// Test creating a TUID with a provided timestamp
	now := time.Now()
	tuid := NewIDWithTime(now)
	// Test extraction of the embedded timestamp
	ts, err := tuid.Time()
	if err != nil {
		t.Error(err)
	}
	if ts.UnixNano() != now.UnixNano() {
		t.Error("error: embedded timestamp did not match")
	}
}

func TestNewTuidWithTimeAndEntropy(t *testing.T) {
	tuid := NewID()
	info, err := tuid.Info()
	if err != nil {
		t.Error(err)
	}
	tuid2 := NewIDWithTimeAndEntropy(info.Timestamp, info.Entropy)
	if tuid != tuid2 {
		t.Error("tuids did not match")
	}
}

func TestFirstTuidWithTime(t *testing.T) {
	now := time.Now()
	tuid := NewIDWithTime(now)
	first := FirstIDWithTime(now)
	if strings.Compare(string(first), string(tuid)) > 0 {
		t.Error("expected first TUID to sort before a regular TUID")
	}
	entropy, err := first.Entropy()
	if err != nil {
		t.Error(err)
	}
	if entropy != 0 {
		t.Error("expected first TUID to have zero entropy")
	}
}

func TestMinTuid(t *testing.T) {
	minTimestamp := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedMinID := FirstIDWithTime(minTimestamp)
	if MinID != expectedMinID {
		t.Error("expected MinID to be the first TUID with the minimum timestamp")
	}
	minID, _ := TUID(MinID).Int()
	if minID.BitLen() != 92 {
		t.Error("expected min TUID to have 92 bits")
	}
	minEntropy, _ := TUID(MinID).Entropy()
	if minEntropy != 0 {
		t.Error("expected min TUID to have zero entropy")
	}
}

func TestMaxTuid(t *testing.T) {
	maxTimestamp := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedMaxID := FirstIDWithTime(maxTimestamp)
	if MaxID != expectedMaxID {
		t.Error("expected MaxID to be the first TUID with the maximum timestamp")
	}
	maxID, _ := TUID(MaxID).Int()
	if maxID.BitLen() != 94 {
		t.Error("expected max TUID to have 94 bits")
	}
	maxEntropy, _ := TUID(MaxID).Entropy()
	if maxEntropy != 0 {
		t.Error("expected max TUID to have zero entropy")
	}
}

func TestTuid_Int(t *testing.T) {
	tuid := TUID("91Mq07yx9IxHCi5Y")
	expected := "101100110101001001000001111100110011001011011110001101001100001011100111101100111011001011000"
	id, err := tuid.Int()
	if err != nil {
		t.Error(err)
	}
	bits := fmt.Sprintf("%b", id)
	if bits != expected {
		t.Error("received unexpected TUID integer bits")
	}
}

func TestTuid_Time(t *testing.T) {
	expected := time.Date(2021, 3, 8, 5, 54, 9, 208207000, time.UTC)
	ts, err := TUID("91Mq07yx9IxHCi5Y").Time()
	if err != nil {
		t.Error(err)
	}
	if ts.UnixNano() != expected.UnixNano() {
		t.Errorf("error: expected %s, received %s", expected, ts)
	}
}

func TestTuid_Entropy(t *testing.T) {
	tuid := TUID("9AxhASm3k3MVH8se")
	expected := "1011110101011000010001111000100"
	entropy, err := tuid.Entropy()
	if err != nil {
		t.Error(err)
	}
	received := fmt.Sprintf("%b", entropy)
	if received != expected {
		t.Errorf("error: expected %s, received %s", expected, received)
	}
}

func TestDuration(t *testing.T) {
	expected := "32m2.322640879s"
	duration, err := Duration("9AxgffrWr9qCnfIT", "9AxjEL0lPtoGAbLE")
	if err != nil {
		t.Error(err)
	}
	received := duration.String()
	if received != expected {
		t.Errorf("error: expected %s, received %s", expected, received)
	}
}

func TestCompare(t *testing.T) {
	tuid1 := NewID()
	tuid2 := NewID()
	if Compare(tuid1, tuid2) != -1 {
		t.Errorf("expected TUID %s to sort cronologically before %s", tuid1, tuid2)
	}
	if Compare(tuid1, tuid1) != 0 {
		t.Errorf("expected TUID %s to compare equally with %s", tuid1, tuid1)
	}
	if Compare(tuid2, tuid1) != 1 {
		t.Errorf("expected TUID %s to sort cronologically after %s", tuid2, tuid1)
	}
}

func TestIsValid(t *testing.T) {
	data := []struct {
		name  string
		tuid  TUID
		valid bool
	}{
		{"minID", TUID(MinID), true},
		{"maxID", TUID(MaxID), true},
		{"normal", TUID("91Mq07yx9IxHCi5Y"), true},
		{"unusual", TUID("AndIHave16Digits"), true},
		{"16z", TUID("zzzzzzzzzzzzzzzzz"), false},
		{"invalid", TUID("I'mNotATuid!"), false},
		{"sequence", TUID("1500593"), false},
		{"blank", TUID(""), false},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			v := IsValid(d.tuid)
			if v != d.valid {
				t.Errorf("expected %s validity to be %v", d.tuid, d.valid)
			}
		})
	}
}

func TestUniqueIDs(t *testing.T) {
	count := 100000
	ids := map[TUID]struct{}{}
	startTime := time.Now()
	for i := 0; i < count; i++ {
		tuid := NewID()
		ids[tuid] = struct{}{}
	}
	duration := time.Since(startTime)
	rate := duration.Nanoseconds() / int64(count)
	fmt.Printf("generated %d Tuids in %s (%d ns/TUID)\n", len(ids), duration, rate)
	if len(ids) != count {
		t.Errorf("expected %d Tuids, received %d", count, len(ids))
	}
}
