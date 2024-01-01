package parse

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatTime(t *testing.T) {
	for _, tc := range []struct {
		t        IsoTimeNode
		expected string
	}{
		{IsoTimeNode{Time: time.Unix(0, 0)}, "1970-01-01T00:00:00+00:00"},
	} {
		t.Run(fmt.Sprintf("%s==%s", tc.t, tc.expected), func(t *testing.T) {
			assert.Equal(t, tc.expected, fmt.Sprint(tc.t))
		})
	}
}
