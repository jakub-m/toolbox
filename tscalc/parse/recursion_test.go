package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecursiveParserRecursionOnRight(t *testing.T) {
	//LogEnabled = true
	input := `1+2+3`
	number := Regex(`[0-9]+`)
	plus := Literal(`+`)
	plusExprRef := Ref()
	numberPlusSomething := Sequence(number, plus, plusExprRef)
	plusExpr := FirstOf(numberPlusSomething, number)
	plusExprRef.Parser = plusExpr
	parser := plusExpr
	node, rem, err := parser.Parse(NewCursor(input))
	assert.NoError(t, err)
	assert.True(t, rem.Ended())
	assert.Equal(t,
		SequenceNode{
			LiteralNode("1"),
			LiteralNode("+"),
			SequenceNode{
				LiteralNode("2"),
				LiteralNode("+"),
				SequenceNode{
					LiteralNode("3"),
				},
			},
		},
		node)
}

func TestRecursiveWithBracket(t *testing.T) {
	lit := Regex(`a`)
	exprRef := Ref()
	expr := FirstOf(lit, Bracket(exprRef))
	exprRef.Parser = expr

	for _, tc := range []string{
		"a",
		"(a)",
		"((a))",
		"(((a)))",
	} {
		t.Run(tc, func(t *testing.T) {
			node, rest, err := exprRef.Parse(NewCursor("a"))
			assert.NoError(t, err)
			assert.True(t, rest.Ended())
			assert.Equal(t, LiteralNode("a"), node)
		})
	}
}
