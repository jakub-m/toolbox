package parse

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePeriod(t *testing.T) {
	node, rest, err := Period("100sec")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, node, PeriodNode{100})
}

func TestFormatPeriod(t *testing.T) {
	for _, tc := range []struct {
		seconds  float64
		expected string
	}{
		{0, "0sec"},
		{100, "100sec"},
		{-100, "-100sec"},
	} {
		t.Run(fmt.Sprintf("%f == %s", tc.seconds, tc.expected), func(t *testing.T) {
			assert.Equal(t, tc.expected, PeriodNode{Seconds: tc.seconds}.String())
		})
	}

}
