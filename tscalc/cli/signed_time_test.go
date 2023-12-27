package main

import (
	p "lib/tscalc/parse"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseSignedPeriod(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected time.Duration
	}{
		{"-3s", -3 * time.Second},
		{"3s", 3 * time.Second},
		{"+3s", 3 * time.Second},
		{"- 3s", -3 * time.Second},
		{"+ 3s", 3 * time.Second},
	} {
		node, rest, err := SignedPeriod.Parse(p.NewCursor(tc.input))
		assert.NoError(t, err)
		assert.True(t, rest.Ended(), rest.String())
		assert.Equal(t, p.PeriodNode(tc.expected), node)
	}
}

func TestParseSignedPeriodSequence(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected time.Duration
	}{
		{"3s + 3s", 6 * time.Second},
		{"3s+3s", 6 * time.Second},
		{"3s - 1s", 2 * time.Second},
		{"3s-1s", 2 * time.Second},
		{"3s-1s-1s", 1 * time.Second},
		{"3s - 1s - 1s", 1 * time.Second},
		{"10s - 4s - 3s-2s-1s", 0 * time.Second},
	} {
		node, rest, err := SignedPeriodExpr.Parse(p.NewCursor(tc.input))
		assert.NoError(t, err)
		assert.True(t, rest.Ended())
		assert.Equal(t, p.PeriodNode(tc.expected), node)
	}
}
