package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func GetParser() Parser {
	time_ := FirstOf(IsoTime, EpochTime)
	return FirstOf(Addition(time_, Period), Subtraction(time_, time_), time_)
}

type Node any

// Parser is not a pure function because there might be parsers that will have some minimal state.
type Parser interface {
	// Parse returns the node if found, the remaining string and the error if any. If the node is not found, then
	// do not return error, just return nil node. An error means that the program should stop immediatelly.
	Parse(input string) (Node, string, error)
}

type AddNode struct {
	Left, Right Node
}

// Addition returns <something> + <something else>. Addition parser takes care of whitespace.
func Addition(leftParser, rightParser Parser) Parser {
	pf := func(input string) (Node, string, error) {
		node, rest, err := Sequence(
			ContinuedBy(leftParser, WhitespaceEOL),
			ContinuedBy(RegexLiteral(`\+`), WhitespaceEOL),
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
	return funcParser{pf}
}

type SubNode struct {
	Left, Right Node
}

// Subtraction returns <time> - <time> parser. It takes care of the whitespace around - sign.
func Subtraction(leftParser, rightParser Parser) Parser {
	pf := func(input string) (Node, string, error) {
		node, rest, err := Sequence(
			ContinuedBy(leftParser, WhitespaceEOL),
			ContinuedBy(RegexLiteral(`\-`), WhitespaceEOL),
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
	return funcParser{pf}
}

type PeriodNode time.Duration

func (n PeriodNode) String() string {
	return fmt.Sprint(time.Duration(n))
}

type periodStr struct{}

var Period = periodStr{}

func (p periodStr) Parse(input string) (Node, string, error) {
	pat := regexp.MustCompile(`^(?:\d+\w+)+`)
	indices := pat.FindStringSubmatchIndex(input)
	if indices == nil {
		return nil, input, nil
	}
	match, rest := input[indices[0]:indices[1]], input[indices[1]:]
	d, err := time.ParseDuration(match)
	if err != nil {
		return 0, input, fmt.Errorf("error while parsing duration %s: %w", d, err)
	}
	return PeriodNode(d), rest, nil
}

type SequenceNode []Node

// Sequence returns a sequence if all the parsers successfully parse.
func Sequence(parsers ...Parser) Parser {
	seq := SequenceNode{}
	pf := func(input string) (Node, string, error) {
		actualInput := input
		for _, parser := range parsers {
			node, rest, err := parser.Parse(actualInput)
			if err != nil || node == nil {
				return node, input, err
			}
			actualInput = rest
			seq = append(seq, node)
		}
		return seq, actualInput, nil
	}
	return funcParser{pf}
}

type RegexLiteralNode struct{ match string }

func RegexLiteral(pat string) Parser {
	pf := func(input string) (Node, string, error) {
		pat := regexp.MustCompile(pat)
		indices := pat.FindStringIndex(input)
		if indices == nil {
			return nil, input, nil
		}
		match, rest := input[indices[0]:indices[1]], input[indices[1]:]
		return RegexLiteralNode{match: match}, rest, nil
	}
	return funcParser{pf}
}

type WhitespaceNode struct{}
type whitespaceEOLStr struct{}

var WhitespaceEOL = whitespaceEOLStr{}

func (p whitespaceEOLStr) Parse(input string) (Node, string, error) {
	pat := regexp.MustCompile(`^(\s+)|^$`)
	indices := pat.FindStringIndex(input)
	if indices == nil {
		return nil, input, nil
	}
	rest := input[indices[1]:]
	return WhitespaceNode{}, rest, nil
}

// ContinuedBy returns result of the main parser only if the reminder of main is parsed by the continuation parser.
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
	return funcParser{pf}
}

func FirstOf(parsers ...Parser) Parser {
	pf := func(input string) (Node, string, error) {
		for _, p := range parsers {
			node, rest, err := p.Parse(input)
			if err != nil || node != nil {
				return node, rest, err
			}
		}
		return nil, input, nil
	}
	return funcParser{pf}
}

type EpochTimeNode float64

const isoFormat = "2006-01-02T15:04:05-07:00"

// func (t EpochTimeNode) FormatISO() string {
// 	sec, frac := math.Modf(float64(t))
// 	return time.Unix(int64(sec), int64(1_000_000_000*frac)).UTC().Format(isoFormat)
// }

func (t EpochTimeNode) ToIsoTimeNode() IsoTimeNode {
	return IsoTimeNode(time.UnixMicro(int64(1_000_000 * float64(t))))
}

func (n EpochTimeNode) String() string {
	return fmt.Sprintf("%f", n)
}

var EpochTime = epochTimeStr{}

type epochTimeStr struct{}

func (p epochTimeStr) Parse(input string) (Node, string, error) {
	pat := regexp.MustCompile(`^\d+(\.\d+)?`)
	indices := pat.FindStringIndex(input)
	if indices == nil {
		return 0, input, nil
	}
	match := input[indices[0]:indices[1]]
	t, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, input, fmt.Errorf("error while parsing %s: %w", input, err)
	}
	return EpochTimeNode(t), input[indices[1]:], nil
}

type IsoTimeNode time.Time

func (n IsoTimeNode) String() string {
	return time.Time(n).UTC().Format(isoFormat)
}

func (n IsoTimeNode) ToEpochTimeNode() EpochTimeNode {
	t := time.Time(n).Unix()
	return EpochTimeNode(t)
}

type isoTimeStr struct{}

var IsoTime = isoTimeStr{}

func (p isoTimeStr) Parse(input string) (Node, string, error) {
	pat := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2}`)
	indices := pat.FindStringIndex(input)
	if indices == nil {
		return nil, input, nil
	}
	match := input[indices[0]:indices[1]]
	t, err := time.Parse(isoFormat, match)
	if err != nil {
		return nil, input, fmt.Errorf("error while parsing %s: %w", input, err)
	}
	return IsoTimeNode(t), input[indices[1]:], nil
}

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
	return funcParser{pf}
}

// funcParser wraps function into a Parser.
type funcParser struct {
	pf func(input string) (Node, string, error)
}

func (p funcParser) Parse(input string) (Node, string, error) {
	return p.pf(input)
}
