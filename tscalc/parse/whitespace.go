package parse

import "regexp"

type WhitespaceNode struct{}
type whitespaceEOLStr struct{}

var WhitespaceEOL = whitespaceEOLStr{}

func (w whitespaceEOLStr) String() string {
	return "<ws>"
}

func (p whitespaceEOLStr) Parse(input Cursor) (Node, Cursor, error) {
	pat := regexp.MustCompile(`^(\s+)|^$`)
	indices := pat.FindStringIndex(input.String())
	if indices == nil {
		return nil, input, nil
	}
	rest := input.Advance(indices[1])
	return WhitespaceNode{}, rest, nil
}
