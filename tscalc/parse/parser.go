package parse

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

func GetParser() ParseFunc {
	time_ := FirstOf(IsoTime, EpochTime)
	return FirstOf(Addition(time_, Period), time_)
}

type Node any

// ParseFunc returns the node if found, the remaining string and the error if any. If the node is not found, then
// do not return error, just return nil node. An error means that the program should stop immediatelly.
type ParseFunc func(input string) (Node, string, error)

type AddNode struct {
	Left, Right Node
}

// Addition returns <something> + <something else>. Addition parser takes care of whitespace.
func Addition(leftParser, rightParser ParseFunc) ParseFunc {
	return func(input string) (Node, string, error) {
		node, rest, err := Sequence(
			ContinuedBy(leftParser, WhitespaceEOL),
			ContinuedBy(RegexLiteral(`\+`), WhitespaceEOL),
			rightParser,
		)(input)
		if err != nil || node == nil {
			return node, rest, err
		}
		seq, ok := node.(SequenceNode)
		if !(ok && len(seq) == 3) {
			return nil, input, fmt.Errorf("BUG! Expected 3 element sequence in Addition, got: %v", seq)
		}
		return &AddNode{seq[0], seq[2]}, rest, nil
	}
}

type PeriodNode struct {
	Seconds float64
}

func Period(input string) (Node, string, error) {
	pat := regexp.MustCompile(`^(\d+)sec`)
	indices := pat.FindStringSubmatchIndex(input)
	if indices == nil {
		return nil, input, nil
	}
	matchSec, rest := input[indices[2]:indices[3]], input[indices[1]:]
	seconds, err := strconv.ParseFloat(matchSec, 64)
	if err != nil {
		return 0, input, fmt.Errorf("error while parsing %s: %w", input, err)
	}
	node := PeriodNode{seconds}
	return node, rest, nil
}

type SequenceNode []Node

// Sequence returns a sequence if all the parsers successfully parse.
func Sequence(parsers ...ParseFunc) ParseFunc {
	seq := SequenceNode{}
	return func(input string) (Node, string, error) {
		actualInput := input
		for _, parser := range parsers {
			node, rest, err := parser(actualInput)
			if err != nil || node == nil {
				return node, input, err
			}
			actualInput = rest
			seq = append(seq, node)
		}
		return seq, actualInput, nil
	}
}

//
//type EmptyNode struct{}
//
//// Optional retuns a dummy empty node if the parser does not parse
//func Optional(parser ParseFunc) ParseFunc {
//	return func(input string) (Node, string, error) {
//		node, rest, err := parser(input)
//		if err != nil {
//			return node, rest, err
//		}
//		if node == nil {
//			if rest == input {
//				return EmptyNode{}, rest, nil
//			} else {
//				return nil, rest, fmt.Errorf("BUG! The parser in the optional node did not parse anything but consumed input")
//			}
//		}
//		return node, rest, nil
//	}
//}

type RegexLiteralNode struct{ match string }

func RegexLiteral(pat string) ParseFunc {
	return func(input string) (Node, string, error) {
		pat := regexp.MustCompile(pat)
		indices := pat.FindStringIndex(input)
		if indices == nil {
			return nil, input, nil
		}
		match, rest := input[indices[0]:indices[1]], input[indices[1]:]
		return RegexLiteralNode{match: match}, rest, nil
	}
}

type WhitespaceNode struct{}

func WhitespaceEOL(input string) (Node, string, error) {
	pat := regexp.MustCompile(`^(\s+)|^$`)
	indices := pat.FindStringIndex(input)
	if indices == nil {
		return nil, input, nil
	}
	rest := input[indices[1]:]
	node := WhitespaceNode{}
	return &node, rest, nil
}

// ContinuedBy returns result of the main parser only if the reminder of main is parsed by the continuation parser.
func ContinuedBy(main, continuation ParseFunc) ParseFunc {
	return func(input string) (Node, string, error) {
		node, rest, err := main(input)
		if err != nil {
			return node, rest, err
		}
		contNode, contRest, contErr := continuation(rest)
		if contErr != nil {
			return contNode, contRest, contErr
		}
		if contNode != nil {
			return node, contRest, nil
		}
		return nil, input, nil
	}
}

func FirstOf(parsers ...ParseFunc) ParseFunc {
	return func(input string) (Node, string, error) {
		for _, p := range parsers {
			node, rest, err := p(input)
			if err != nil || node != nil {
				return node, rest, err
			}
		}
		return nil, input, nil
	}
}

type EpochTimeNode float64

const isoFormat = "2006-01-02T15:04:05-07:00"

func (t EpochTimeNode) FormatISO() string {
	sec, frac := math.Modf(float64(t))
	return time.Unix(int64(sec), int64(1_000_000_000*frac)).UTC().Format(isoFormat)
}

// func (t EpochTimeNode) ToISOTimeNode() IsoTimeNode {
// 	return IsoTimeNode(time.UnixMicro(int64(1_000_000 * float64(t))))
// }

func EpochTime(input string) (Node, string, error) {
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

func (n IsoTimeNode) FormatTimestamp() string {
	t := time.Time(n)
	return fmt.Sprintf("%d", t.Unix())
}

func (n IsoTimeNode) ToEpochTimeNode() EpochTimeNode {
	t := time.Time(n).Unix()
	return EpochTimeNode(t)
}

func IsoTime(input string) (Node, string, error) {
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
	node := IsoTimeNode(t)
	return &node, input[indices[1]:], nil
}
