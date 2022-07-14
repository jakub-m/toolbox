package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniq(t *testing.T) {
	var out bytes.Buffer
	mainInternal(options{
		selectorSpec: "2",
		separator:    "",
		showCount:    false,
	}, reader(`
1 foo
2 foo
3 bar
4 bar 
5 bar
6 quux
`), io.Writer(&out))
	assert.Equal(t, ltrim(`
1 foo
3 bar
6 quux
`), out.String())
}

func reader(s string) io.Reader {
	return strings.NewReader(ltrim(s))
}

func ltrim(s string) string {
	return strings.TrimLeft(s, "\n ")
}
