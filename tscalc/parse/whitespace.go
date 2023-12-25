package parse

import "regexp"

type WhitespaceNode struct{}
type whitespaceEOLStr struct{}

var WhitespaceEOL = whitespaceEOLStr{}

func (w whitespaceEOLStr) String() string {
	return "<ws>"
}

func (p whitespaceEOLStr) Parse(input string) (Node, string, error) {
	pat := regexp.MustCompile(`^(\s+)|^$`)
	indices := pat.FindStringIndex(input)
	if indices == nil {
		return nil, input, nil
	}
	rest := input[indices[1]:]
	return WhitespaceNode{}, rest, nil
}
