package parse

import "fmt"

func Bracket(inner Parser) Parser {
	pf := func(input string) (Node, string, error) {
		node, rest, err := Sequence(
			RegexLiteral(`\(\s*`),
			inner,
			RegexLiteral(`\s*\)`),
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
	return funcParser{pf: pf, name: name}
}
