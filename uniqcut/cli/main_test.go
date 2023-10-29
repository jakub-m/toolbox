package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqSimple(t *testing.T) {
	var out bytes.Buffer
	mainInternal(args{
		selectorSpec: "2",
	}, reader(`
a foo
b foo
c bar
d bar 
e bar
f quux
`), io.Writer(&out))
	assert.Equal(t, ltrim(`
b foo
e bar
f quux
`), out.String())
}

func TestUniqCount(t *testing.T) {
	var out bytes.Buffer
	mainInternal(args{
		selectorSpec: "2",
		showCount:    true,
	}, reader(`
a foo
c bar
d bar 
e bar
f quux
`), io.Writer(&out))
	assert.Equal(t, ltrim(`
1	a foo
3	e bar
1	f quux
`), out.String())
}

func reader(s string) io.Reader {
	return strings.NewReader(ltrim(s))
}

func ltrim(s string) string {
	return strings.TrimLeft(s, "\n ")
}
