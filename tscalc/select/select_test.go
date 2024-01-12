package parse_sql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSql(t *testing.T) {
	for _, tc := range []string{
		// `select 1`,
		// `select * from foo`,
		`select a, b, c, d from foo`,
	} {
		tree, err := Parse(tc)
		assert.NoError(t, err)
		assert.NotNil(t, tree)
		assert.Equal(t, "", fmt.Sprint(tree)) // just to print
	}
}
