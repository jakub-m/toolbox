package parse

import (
	"fmt"
	"strings"
)

type SequenceNode []Node

func (s SequenceNode) String() string {
	substrings := make([]string, len(s))
	for i, n := range s {
		substrings[i] = fmt.Sprintf("%s", n)
	}
	return fmt.Sprintf("[%s]", strings.Join(substrings, ", "))
}

// Sequence returns a sequence if all the parsers successfully parse.
func Sequence(parsers ...Parser) Parser {
	pf := func(input Cursor) (Node, Cursor, error) {
		seq := SequenceNode{}
		Logf("Sequence of %d parsers on: %s$", len(parsers), input)
		LogIndentInc()
		defer LogIndentDec()
		actualCur := input
		for i := 0; i < len(parsers) && !actualCur.Ended(); i++ {
			parser := parsers[i]
			Logf("Sequence[%d/%d] %s on: %s$", i+1, len(parsers), parser, actualCur)
			LogIndentInc()
			node, rest, err := parser.Parse(actualCur)
			if err != nil || node == nil {
				LogIndentDec()
				return node, input, err
			}
			actualCur = rest
			Logf("Sequence[%d/%d] match, rest: %s$", i+1, len(parsers), rest)
			seq = append(seq, node)
			LogIndentDec()
		}
		return seq, actualCur, nil
	}
	parserStrings := make([]string, len(parsers))
	for i := range parsers {
		parserStrings[i] = fmt.Sprint(parsers[i])
	}
	name := "(" + strings.Join(parserStrings, " ") + ")"
	return FuncParser{Fn: pf, Name: name}
}

func (s SequenceNode) RemoveEmpty() SequenceNode {
	filtered := SequenceNode{}
	for _, n := range s {
		if n == EmptyNode {
			continue
		}
		filtered = append(filtered, n)
	}
	return filtered
}

type repeatedParser struct {
	parser Parser
}

func Repeated(p Parser) repeatedParser {
	return repeatedParser{parser: p}
}

func (p repeatedParser) Parse(input Cursor) (Node, Cursor, error) {
	Logf("Repeat %s", p.parser)
	LogIndentInc()
	defer LogIndentDec()
	seq := SequenceNode{}
	rest := input
	for !rest.Ended() {
		node, c, err := p.parser.Parse(rest)
		rest = c
		if err != nil {
			return nil, rest, err
		}
		if node == nil {
			break
		}
		seq = append(seq, node)
	}
	return seq, rest, nil
}
