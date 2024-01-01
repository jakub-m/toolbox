package parse

import (
	"fmt"
	"strings"
)

type Node interface {
	// Cursor method returns the cursor that can be used to find which part of the input contributed to the Node.
	// This can be used for printing nice error messages pointing at the particular place in the input.
	Cursor() Cursor
}

// CursorError is an interface of an error that has information about Cursor. This information can be used to print
// helpful error messages.
type CursorError interface {
	error
	Node
}

// Parser is not a pure function because there might be parsers that will have some minimal state.
type Parser interface {
	// Parse returns the node if found, the remaining string and the error if any. If the node is not found, then
	// do not return error, just return nil node. An error means that the program should stop immediatelly.
	Parse(Cursor) (Node, Cursor, error)
}

// Cursor is used to point where in the parsed string are we now.
type Cursor struct {
	Input string
	Pos   int
}

func NewCursor(s string) Cursor {
	return Cursor{Input: s, Pos: 0}
}

func (c Cursor) Ended() bool {
	return c.Pos >= len(c.Input)
}

func (c Cursor) String() string {
	return c.Input[c.Pos:]
}

func (c Cursor) Advance(nBytes int) Cursor {
	return Cursor{
		Input: c.Input,
		Pos:   c.Pos + nBytes,
	}
}

// ContinuedBy returns result of the main parser only if the reminder of main is parsed by the continuation parser.
// The continueation match is not returned, it is consumed and ignored.
func ContinuedBy(main, continuation Parser) Parser {
	pf := func(input Cursor) (Node, Cursor, error) {
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
	return FuncParser{Fn: pf, Name: name}
}

func FirstOf(parsers ...Parser) Parser {
	pf := func(input Cursor) (Node, Cursor, error) {
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
	return FuncParser{Fn: pf, Name: name}
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

func (p *RefStr) Parse(input Cursor) (Node, Cursor, error) {
	if p.Parser == nil {
		return nil, input, fmt.Errorf("BUG! RefStr.Parse is nil")
	}
	return p.Parser.Parse(input)
}

// FuncParser wraps function into a Parser.
type FuncParser struct {
	Fn   func(input Cursor) (Node, Cursor, error)
	Name string
}

func (p FuncParser) String() string {
	if p.Name == "" {
		return "funcParser{...}"
	} else {
		return p.Name
	}
}

func (p FuncParser) Parse(input Cursor) (Node, Cursor, error) {
	return p.Fn(input)
}

type optionalParser struct {
	parser Parser
}

func Optional(p Parser) optionalParser {
	return optionalParser{parser: p}
}

func (p optionalParser) String() string {
	return fmt.Sprintf("(%s)?", p.parser)
}

// EmptyNodeis the type of the `EmptyNode`. Usually it's more handy to use EmptyNode but in `switch...case...` you can use `EmptyNodeType`.
type EmptyNode struct {
	cursor Cursor
}

func (n EmptyNode) Cursor() Cursor {
	return n.cursor
}

func (n EmptyNode) String() string {
	return "<nil>"
}

func (p optionalParser) Parse(input Cursor) (Node, Cursor, error) {
	Logf("Optional %s", p.parser)
	LogIndentInc()
	defer LogIndentDec()
	node, rest, err := p.parser.Parse(input)
	if err != nil {
		return node, rest, err
	}
	if node == nil {
		return EmptyNode{input}, rest, nil
	}
	return node, rest, err
}
