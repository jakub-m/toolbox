package parse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdditionLiterals(t *testing.T) {
	addition := Addition(Regex(`a`), Regex(`b`))
	node, rest, err := addition.Parse(NewCursor("a + b"))
	assert.NoError(t, err)
	assert.True(t, rest.Ended())
	assert.Equal(t, AddNode{LiteralNode("a"), LiteralNode("b")}, node)
}

func TestAdditionTimeAndPeriod(t *testing.T) {
	addition := Addition(EpochTime, Period)

	node, rest, err := addition.Parse(NewCursor("100 + 10s"))
	assert.NoError(t, err)
	assert.True(t, rest.Ended())
	assert.Equal(t, AddNode{EpochTimeNode(100), PeriodNode(10 * time.Second)}, node)
}

func TestAddPeriods(t *testing.T) {
	node, rest, err := Addition(Period, Period).Parse(NewCursor("10s + 10s"))
	assert.NoError(t, err)
	assert.True(t, rest.Ended())
	assert.Equal(t, AddNode{PeriodNode(10 * time.Second), PeriodNode(10 * time.Second)}, node)
}
