package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	input := "Foo Bar: hello_world: 728, blabla: 728, pony: FOO, color: BAR"
	actual := FormatLine(input)
	expected := `Foo 
Bar: hello_world: 728, 
blabla: 728, 
pony: FOO, 
color: BAR`
	assert.Equal(t, expected, actual)
}
