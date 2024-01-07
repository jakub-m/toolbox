package parse_sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSql(t *testing.T) {
	t.Skip()
	for _, tc := range []string{
		// `select 1`,
		// `select * from foo`,
		`select a, b, c, d from foo`,
	} {
		tree, err := Parse(tc)
		assert.NoError(t, err)
		assert.NotNil(t, tree)
	}
}
