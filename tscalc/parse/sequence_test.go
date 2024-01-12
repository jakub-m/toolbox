package parse

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequenceComplete(t *testing.T) {
	node, rest, err := Sequence(Regex(`aaa`), Regex(`bbb`), Regex(`ccc`)).Parse(NewCursor("aaabbbccc"))
	assert.NoError(t, err)
	assert.True(t, rest.Ended(), rest.String())
	assert.Equal(t, `["aaa" "bbb" "ccc"]`, fmt.Sprint(node))
}

func TestSequenceIncomplete(t *testing.T) {
	node, rest, err := Sequence(Regex(`aaa`), Regex(`bbb`), Regex(`ccc`)).Parse(NewCursor("aaabbbc"))
	assert.NoError(t, err)
	assert.Equal(t, "aaabbbc", rest.String())
	assert.Equal(t, node, nil)
}

func TestRepeated(t *testing.T) {
	node, rest, err := Repeated(Literal("x")).Parse(NewCursor("xxx"))
	assert.NoError(t, err)
	assert.True(t, rest.Ended(), rest.String())
	assert.Equal(t, `["x" "x" "x"]`, fmt.Sprint(node))
}

func TestSequenceWithRepeated(t *testing.T) {
	node, rest, err := Sequence(Repeated(Literal("x")), Repeated(Literal("y"))).Parse(NewCursor("xxyy"))
	assert.NoError(t, err)
	assert.True(t, rest.Ended(), rest.String())
	assert.Equal(t, `[["x" "x"] ["y" "y"]]`, fmt.Sprint(node))
}
