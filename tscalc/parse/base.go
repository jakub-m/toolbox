package parse

import (
	"fmt"
	"strings"
)

type Node any

// Parser is not a pure function because there might be parsers that will have some minimal state.
type Parser interface {
	// Parse returns the node if found, the remaining string and the error if any. If the node is not found, then
	// do not return error, just return nil node. An error means that the program should stop immediatelly.
	Parse(input string) (Node, string, error)
}

type SequenceNode []Node

// Sequence returns a sequence if all the parsers successfully parse.
func Sequence(parsers ...Parser) Parser {
	pf := func(input string) (Node, string, error) {
		seq := SequenceNode{}
		Logf("Sequence of %d parsers on: %s$", len(parsers), input)
		LogIndentInc()
		defer LogIndentDec()
		actualInput := input
		for i, parser := range parsers {
			Logf("Sequence[%d/%d] %s on: %s$", i+1, len(parsers), parser, actualInput)
			node, rest, err := parser.Parse(actualInput)
			if err != nil || node == nil {
				return node, input, err
			}
			actualInput = rest
			Logf("Sequence[%d/%d] match, rest: %s$", i+1, len(parsers), rest)
			// Logf("Sequence[%d/%d] was %s, next node: %s", i+1, len(parsers), seq, node)
			seq = append(seq, node)
		}
		return seq, actualInput, nil
	}
	parserStrings := make([]string, len(parsers))
	for i := range parsers {
		parserStrings[i] = fmt.Sprint(parsers[i])
	}
	name := "(" + strings.Join(parserStrings, " ") + ")"
	return funcParser{pf: pf, name: name}
}

// ContinuedBy returns result of the main parser only if the reminder of main is parsed by the continuation parser.
// The continueation match is not returned, it is consumed and ignored.
func ContinuedBy(main, continuation Parser) Parser {
	pf := func(input string) (Node, string, error) {
		node, rest, err := main.Parse(input)
		if err != nil {
			return node, rest, err
		}
		contNode, contRest, contErr := continuation.Parse(rest)
		if contErr != nil {
			return contNode, contRest, contErr
		}
		if contNode != nil {
			return node, contRest, nil
		}
		return nil, input, nil
	}
	name := fmt.Sprintf("%s (?=%s)", main, continuation)
	return funcParser{pf: pf, name: name}
}

func FirstOf(parsers ...Parser) Parser {
	pf := func(input string) (Node, string, error) {
		Logf("FirstOf %d parsers on: %s$", len(parsers), input)
		LogIndentInc()
		defer LogIndentDec()
		for i, p := range parsers {
			Logf("FirstOf[%d/%d]: %s", i+1, len(parsers), p)
			node, rest, err := p.Parse(input)
			if err != nil || node != nil {
				if node != nil {
					Logf("FirstOf[%d/%d] match, rest: %s$", i+1, len(parsers), rest)
				}
				return node, rest, err
			}
		}
		return nil, input, nil
	}
	parserStrings := make([]string, len(parsers))
	for i := range parsers {
		parserStrings[i] = fmt.Sprint(parsers[i])
	}
	name := "(" + strings.Join(parserStrings, " | ") + ")"
	return funcParser{pf: pf, name: name}
}

// RefStr is used to build recursive parsers.
type RefStr struct {
	Parser Parser
	Name   string
}

func (p *RefStr) String() string {
	//return fmt.Sprint(p.Parser)
	if p.Name == "" {
		return "<ref>"
	} else {
		return p.Name
	}
}

func Ref() *RefStr {
	return &RefStr{}
}

func (p *RefStr) Parse(input string) (Node, string, error) {
	if p.Parser == nil {
		return nil, input, fmt.Errorf("BUG! RefStr.Parse is nil")
	}
	return p.Parser.Parse(input)
}

// funcParser wraps function into a Parser.
type funcParser struct {
	pf   func(input string) (Node, string, error)
	name string
}

func (p funcParser) String() string {
	if p.name == "" {
		return "funcParser{...}"
	} else {
		return p.name
	}
}

func (p funcParser) Parse(input string) (Node, string, error) {
	return p.pf(input)
}
