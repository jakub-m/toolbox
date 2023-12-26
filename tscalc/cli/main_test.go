package main

import (
	"fmt"
	"lib/tscalc/parse"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	parse.LogEnabled = true
}

func TestCalc(t *testing.T) {
	nowFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	for _, tc := range []struct {
		input    string
		expected string
	}{
		{"100", "1970-01-01T00:01:40+00:00"},
		{"  100", "1970-01-01T00:01:40+00:00"},
		{"1970-01-01T00:01:40+00:00", "100.000000"},
		{"1m + 1s", "1m1s"},
		// {"1m+1s", "1m1s"}, // TODO handle spaces
		{"1970-01-01T00:00:00+00:00 + 1m40s", "1970-01-01T00:01:40+00:00"},
		{"1970-01-01T00:01:40+00:00 - 1970-01-01T00:00:00+00:00", "1m40s"},
		{"200 - 100", "1m40s"},
		{"NOW - 1h", "1969-12-31T23:00:00+00:00"},
		{"NOW - NOW", "0s"},
		//{"NOW - NOW + 1h", "1h0m0s"}, // TODO shoudl be ((now - now):period-expr + 1h):period-expr, not (now - (now + 1h):date-expr):period-expr
		{"1s + 1s - 1s", "1s"},
		// {"1s - 1s + 1s", "1s"}, // TODO recognise "sub period" as "add minus period"?
	} {
		t.Run(fmt.Sprintf("%s == %s", tc.input, tc.expected), func(t *testing.T) {
			actual, err := handleLine(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestNow(t *testing.T) {
	nowFunc = func() time.Time {
		return time.Unix(42, 0)
	}
	actual, err := handleLine("NOW")
	assert.NoError(t, err)
	assert.Equal(t, "1970-01-01T00:00:42+00:00", actual)
}
