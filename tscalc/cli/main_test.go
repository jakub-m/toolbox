package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimestampTransforms(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected string
	}{
		{"100", "1970-01-01T00:01:40+00:00"},
		{"1970-01-01T00:01:40+00:00", "100"},
	} {
		t.Run(fmt.Sprintf("%s == %s", tc.input, tc.expected), func(t *testing.T) {
			actual, err := handleLine(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}