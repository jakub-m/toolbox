package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func FormatLine(line string) string {
	parsers := []Parser{
		ParseWordWithTrailingSpace,
		GetSeriesParser(ParseKeyColonValue),
	}

	var output strings.Builder

	input := line
	nodes := []Node{}
	for {
		nodeCountAtStart := len(nodes)
		for _, parse := range parsers {
			log.Println("Parsing input:", input)
			node, reminder, err := parse(input)
			if err != nil {
				panic(err)
			}
			if node != nil {
				log.Printf("Got node %T: %v", node, node)
				nodes = append(nodes, node)
			}
			input = reminder
		}
		if len(nodes) == nodeCountAtStart {
			nodes = append(nodes, StringNode(input))
			break
		}
	}
	for _, node := range nodes {
		output.WriteString(node.StringIndent(0))
	}

	return output.String()
}

type Node interface {
	StringIndent(level int) string
}

type StringNode string

func (n StringNode) StringIndent(level int) string {
	return string(n)
}

// Parser takes a string at input, and consumes some part of it turning the consumed part into a Node. The reminder is returned,
// or empty string if consumed everything.
type Parser func(input string) (Node, string, error)

type Series struct {
	nodes []Node
}

func (s Series) StringIndent(level int) string {
	var b strings.Builder
	b.WriteString("\n")
	for _, n := range s.nodes {
		b.WriteString(strings.Repeat(" ", level+1))
		nodeString := n.StringIndent(level + 1)
		b.WriteString(nodeString)
		b.WriteString("\n")
	}
	return b.String()
}

func GetSeriesParser(parser Parser) Parser {
	seriesParser := func(input string) (Node, string, error) {
		series := Series{nodes: []Node{}}
		for {
			node, remainder, err := parser(input)
			if err != nil {
				return nil, input, err
			}
			if remainder == input {
				if len(series.nodes) == 0 {
					return nil, remainder, nil
				} else {
					return series, remainder, nil
				}
			}
			if node == nil {
				return nil, input, fmt.Errorf("parser %v returned no node but yet consumed input, this should not happen", parser)
			}
			series.nodes = append(series.nodes, node)
			input = remainder
		}
	}
	return seriesParser
}

type KeyColonValue string

func (k KeyColonValue) StringIndent(level int) string {
	return string(k)
}

func ParseKeyColonValue(input string) (Node, string, error) {
	re := regexp.MustCompile(`^(\w+:\s+\w+)(?:,\s+)?`)
	match := re.FindStringSubmatchIndex(input)
	if len(match) == 0 {
		return nil, input, nil
	}
	if len(match) != 4 {
		return nil, input, fmt.Errorf("the returned match for string %v is wrong: %v", input, match)
	}
	_, matchEnd, groupStart, groupEnd := match[0], match[1], match[2], match[3]
	if r, ok := getFirstRune(input[groupEnd:]); ok && r == ':' {
		// Do not consume strings like "foo: bar:". Golang regex does not support negative look-aheads.
		return nil, input, nil
	}
	return KeyColonValue(input[groupStart:groupEnd]), input[matchEnd:], nil
}

type WordWithTrailingSpace string

func (w WordWithTrailingSpace) StringIndent(level int) string {
	return string(w)
}

func ParseWordWithTrailingSpace(input string) (Node, string, error) {
	re := regexp.MustCompile(`^(\S+\s)`)
	match := re.FindStringSubmatchIndex(input)
	if len(match) == 0 {
		return nil, input, nil
	}
	if len(match) != 4 {
		return nil, input, fmt.Errorf("the returned match for string %v is wrong: %v", input, match)
	}
	start, end := match[2], match[3]
	return WordWithTrailingSpace(input[start:end]), input[end:], nil
}

func getFirstRune(s string) (rune, bool) {
	for _, r := range s {
		return r, true
	}
	return '\u0000', false
}
