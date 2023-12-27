package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBracketSimple(t *testing.T) {
	node, rest, err := Bracket(Regex(`aaa`)).Parse(NewCursor("(aaa)"))
	assert.NoError(t, err)
	assert.True(t, rest.Ended())
	assert.Equal(t, node, LiteralNode("aaa"))
}

func TestBracketEmbedded(t *testing.T) {
	node, rest, err := Bracket(Bracket(Regex(`aaa`))).Parse(NewCursor("( ( aaa))"))
	assert.NoError(t, err)
	assert.True(t, rest.Ended())
	assert.Equal(t, node, LiteralNode("aaa"))
}

func TestBracketAdd(t *testing.T) {
	node, rest, err := Bracket(Addition(Regex(`aaa`), Regex(`bbb`))).Parse(NewCursor("(aaa + bbb)"))
	assert.NoError(t, err)
	assert.True(t, rest.Ended())
	assert.Equal(t, node, AddNode{LiteralNode("aaa"), LiteralNode("bbb")})
}
