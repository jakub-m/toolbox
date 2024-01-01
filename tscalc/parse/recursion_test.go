package parse

import (
	"fmt"
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
	assert.Equal(t, `["1" "+" ["2" "+" ["3"]]]`, fmt.Sprint(node))
}
