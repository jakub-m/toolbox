package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePeriod(t *testing.T) {
	node, rest, err := Period("100sec")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, node, PeriodNode{100})
}
