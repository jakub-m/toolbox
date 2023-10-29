package parse

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

func GetParser() ParseFunc {
	return FirstOf(
		ContinuedBy(EpochTime, WhitespaceEOL),
		ContinuedBy(IsoTime, WhitespaceEOL),
	)
}

type Node any

type EpochTimeNode float64

const isoFormat = "2006-01-02T15:04:05-07:00"

func (t EpochTimeNode) FormatISO() string {
	sec, frac := math.Modf(float64(t))
	return time.Unix(int64(sec), int64(1_000_000_000*frac)).UTC().Format(isoFormat)
}

type IsoTimeNode time.Time

func (n IsoTimeNode) FormatTimestamp() string {
	t := time.Time(n)
	return fmt.Sprintf("%d", t.Unix())
}

// ParseFunc returns the node if found, the remaining string and the error if any. If the node is not found, then
// do not return error, just return nil node. An error means that the program should stop immediatelly.
type ParseFunc func(input string) (Node, string, error)

type WhitespaceNode struct{}

func WhitespaceEOL(input string) (Node, string, error) {
	pat := regexp.MustCompile(`(\s+)|$^`)
	indices := pat.FindStringIndex(input)
	if indices == nil {
		return nil, input, nil
	}
	_, rest := input[indices[0]:indices[1]], input[indices[1]:]
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

func EpochTime(input string) (Node, string, error) {
	pat := regexp.MustCompile(`\d+(\.\d+)?`)
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

func IsoTime(input string) (Node, string, error) {
	pat := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2}`)
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
