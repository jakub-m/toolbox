package parse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdditionLiterals(t *testing.T) {
	addition := Addition(RegexLiteral(`a`), RegexLiteral(`b`))
	node, rest, err := addition.Parse("a + b")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, AddNode{RegexLiteralNode{"a"}, RegexLiteralNode{"b"}}, node)
}

func TestAdditionTimeAndPeriod(t *testing.T) {
	addition := Addition(EpochTime, Period)

	node, rest, err := addition.Parse("100 + 10s")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, AddNode{EpochTimeNode(100), PeriodNode(10 * time.Second)}, node)
}

func TestAddPeriods(t *testing.T) {
	node, rest, err := Addition(Period, Period).Parse("10s + 10s")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, AddNode{PeriodNode(10 * time.Second), PeriodNode(10 * time.Second)}, node)
}
