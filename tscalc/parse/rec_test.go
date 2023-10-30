package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecursiveWithBracket(t *testing.T) {
	lit := RegexLiteral(`a`)
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
			node, rest, err := exprRef.Parse("a")
			assert.NoError(t, err)
			assert.Equal(t, "", rest)
			assert.Equal(t, RegexLiteralNode{"a"}, node)
		})
	}
}
