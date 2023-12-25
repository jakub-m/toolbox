package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecursiveParserRecursionOnRight(t *testing.T) {
	//LogEnabled = true
	input := `1+2+3`
	number := RegexLiteral(`[0-9]+`)
	plus := Literal(`+`)
	plusExprRef := Ref()
	numberPlusSomething := Sequence(number, plus, plusExprRef)
	plusExpr := FirstOf(numberPlusSomething, number)
	plusExprRef.Parser = plusExpr
	parser := plusExpr
	node, rem, err := parser.Parse(input)
	assert.NoError(t, err)
	assert.Equal(t, "", rem)
	assert.Equal(t,
		SequenceNode{
			RegexLiteralNode{match: "1"},
			LiteralNode{Exact: "+"},
			SequenceNode{
				RegexLiteralNode{match: "2"},
				LiteralNode{Exact: "+"},
				RegexLiteralNode{match: "3"},
			},
		},
		node)
}
