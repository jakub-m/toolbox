package parse

import (
	"fmt"
	"strings"
)

type SequenceNode struct {
	Nodes  []Node
	cursor Cursor
}

func (s SequenceNode) Cursor() Cursor {
	return s.cursor
}

func (s SequenceNode) Len() int {
	return len(s.Nodes)
}

func (s SequenceNode) String() string {
	substrings := make([]string, s.Len())
	for i, n := range s.Nodes {
		substrings[i] = fmt.Sprintf("%s", n)
	}
	return fmt.Sprintf("[%s]", strings.Join(substrings, " "))
}

// Sequence returns a sequence if all the parsers successfully parse.
func Sequence(parsers ...Parser) Parser {
	pf := func(input Cursor) (Node, Cursor, error) {
		nodes := []Node{}
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
			nodes = append(nodes, node)
			LogIndentDec()
		}
		return SequenceNode{Nodes: nodes, cursor: input}, actualCur, nil
	}
	parserStrings := make([]string, len(parsers))
	for i := range parsers {
		parserStrings[i] = fmt.Sprint(parsers[i])
	}
	name := "(" + strings.Join(parserStrings, " ") + ")"
	return FuncParser{Fn: pf, Name: name}
}

func (s SequenceNode) RemoveEmpty() SequenceNode {
	filtered := []Node{}
	for _, n := range s.Nodes {
		if _, ok := n.(EmptyNode); ok {
			continue
		}
		filtered = append(filtered, n)
	}
	s.Nodes = filtered
	return s
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
	nodes := []Node{}
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
		nodes = append(nodes, node)
	}
	return SequenceNode{Nodes: nodes, cursor: input}, rest, nil
}
