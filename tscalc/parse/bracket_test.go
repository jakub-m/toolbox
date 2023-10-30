package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBracketSimple(t *testing.T) {
	node, rest, err := Bracket(RegexLiteral(`aaa`))("(aaa)")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, node, RegexLiteralNode{"aaa"})
}

func TestBracketEmbedded(t *testing.T) {
	node, rest, err := Bracket(Bracket(RegexLiteral(`aaa`)))("( ( aaa))")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, node, RegexLiteralNode{"aaa"})
}

func TestBracketAdd(t *testing.T) {
	node, rest, err := Bracket(Addition(RegexLiteral(`aaa`), RegexLiteral(`bbb`)))("(aaa + bbb)")
	assert.NoError(t, err)
	assert.Equal(t, "", rest)
	assert.Equal(t, node, AddNode{RegexLiteralNode{"aaa"}, RegexLiteralNode{"bbb"}})
}

// func TestBracketRecursive(t *testing.T) {
// 	lit := RegexLiteral(`a`)

// 	b := FirstOfWith(lit)
// 	b = b.Or(b)
// 	node, rest, err := Bracket(Addition(RegexLiteral(`aaa`), RegexLiteral(`bbb`)))("(aaa + bbb)")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "", rest)
// 	assert.Equal(t, node, AddNode{RegexLiteralNode{"aaa"}, RegexLiteralNode{"bbb"}})

// }
