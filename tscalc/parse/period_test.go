package parse

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParsePeriod(t *testing.T) {
	node, rest, err := Period.Parse("10s")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, node, PeriodNode(10*time.Second))
}

func TestFormatPeriod(t *testing.T) {
	for _, tc := range []struct {
		duration time.Duration
		expected string
	}{
		{0, "0s"},
		{100 * time.Second, "1m40s"},
		{-100 * time.Second, "-1m40s"},
	} {
		t.Run(fmt.Sprintf("%s == %s", tc.duration, tc.expected), func(t *testing.T) {
			assert.Equal(t, tc.expected, fmt.Sprint(PeriodNode(tc.duration)))
		})
	}

}
