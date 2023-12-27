package parse

import "fmt"

func Bracket(inner Parser) Parser {
	pf := func(input Cursor) (Node, Cursor, error) {
		node, rest, err := Sequence(
			Regex(`\(\s*`),
			inner,
			Regex(`\s*\)`),
		).Parse(input)
		if err != nil || node == nil {
			return node, rest, err
		}
		seq := node.(SequenceNode)
		if len(seq) != 3 {
			return nil, input, fmt.Errorf("BUG! Sequenece in bracket should be of len 3, was %d", len(seq))
		}
		return seq[1], rest, nil
	}
	name := fmt.Sprintf("\"(\" %s \")\"", inner)
	return FuncParser{Fn: pf, Name: name}
}
