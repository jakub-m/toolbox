package parse

import "fmt"

type AddNode struct {
	Left, Right Node
}

func (n AddNode) String() string {
	return fmt.Sprintf("{%s + %s}", n.Left, n.Right)
}

// Addition returns <something> + <something else>. Addition parser takes care of whitespace.
func Addition(leftParser, rightParser Parser) Parser {
	pf := func(input Cursor) (Node, Cursor, error) {
		Logf("Addition on: %s$", input)
		LogIndentInc()
		defer LogIndentDec()
		node, rest, err := Sequence(
			ContinuedBy(leftParser, WhitespaceEOL),
			ContinuedBy(Literal(`+`), WhitespaceEOL),
			rightParser,
		).Parse(input)
		if err != nil || node == nil {
			return node, rest, err
		}
		seq, ok := node.(SequenceNode)
		if !(ok && len(seq) == 3) {
			return nil, input, fmt.Errorf("BUG! Expected 3 element sequence in Addition, got: %v", seq)
		}
		return AddNode{seq[0], seq[2]}, rest, nil
	}
	name := "(" + fmt.Sprintf("%s \"+\" %s", leftParser, rightParser) + ")"
	return FuncParser{Fn: pf, Name: name}
}

type SubNode struct {
	Left, Right Node
}

func (n SubNode) String() string {
	return fmt.Sprintf("{%s - %s}", n.Left, n.Right)
}

// Subtraction returns <time> - <time> parser. It takes care of the whitespace around - sign.
func Subtraction(leftParser, rightParser Parser) Parser {
	pf := func(input Cursor) (Node, Cursor, error) {
		Logf("Subtraction on: %s$", input)
		LogIndentInc()
		defer LogIndentDec()
		node, rest, err := Sequence(
			ContinuedBy(leftParser, WhitespaceEOL),
			ContinuedBy(Literal(`-`), WhitespaceEOL),
			rightParser,
		).Parse(input)
		if err != nil || node == nil {
			return node, rest, err
		}
		seq, ok := node.(SequenceNode)
		if !(ok && len(seq) == 3) {
			return nil, input, fmt.Errorf("BUG! Expected 3 element sequence in Subtraction, got: %v", seq)
		}
		return SubNode{seq[0], seq[2]}, rest, nil
	}
	name := "(" + fmt.Sprintf("%s \"-\" %s", leftParser, rightParser) + ")"
	return FuncParser{Fn: pf, Name: name}
}
