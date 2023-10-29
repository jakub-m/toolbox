package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
