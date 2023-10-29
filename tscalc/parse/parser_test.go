package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdditionLiterals(t *testing.T) {
	addition := Addition(RegexLiteral(`a`), RegexLiteral(`b`))
	node, rest, err := addition("a + b")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, AddNode{RegexLiteralNode{"a"}, RegexLiteralNode{"b"}}, node)
}

func TestAdditionTimeAndPeriod(t *testing.T) {
	addition := Addition(EpochTime, Period)

	node, rest, err := addition("100 + 100sec")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, AddNode{EpochTimeNode(100), PeriodNode{100}}, node)
}

func TestAddPeriods(t *testing.T) {
	node, rest, err := Addition(Period, Period)("100sec + 100sec")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, AddNode{PeriodNode{100}, PeriodNode{100}}, node)
}

func TestParsePeriod(t *testing.T) {
	node, rest, err := Period("100sec")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, node, PeriodNode{100})
}

func TestSequenceComplete(t *testing.T) {
	node, rest, err := Sequence(RegexLiteral(`aaa`), RegexLiteral(`bbb`), RegexLiteral(`ccc`))("aaabbbccc")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, node, SequenceNode{RegexLiteralNode{"aaa"}, RegexLiteralNode{"bbb"}, RegexLiteralNode{"ccc"}})
}

func TestSequenceInomplete(t *testing.T) {
	node, rest, err := Sequence(RegexLiteral(`aaa`), RegexLiteral(`bbb`), RegexLiteral(`ccc`))("aaabbbc")
	assert.NoError(t, err)
	assert.Equal(t, "aaabbbc", rest)
	assert.Equal(t, node, nil)
}
