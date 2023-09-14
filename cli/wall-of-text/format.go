package main

import (
	"fmt"
	"regexp"
	"strings"
)

// var regexLineWithListedFields *regexp.Regexp

// func init() {
// regexLineWithListedFields = regexp.MustCompile(`(\w+:\s+\w+)+`)
// }

// const newLine = "\n"

func FormatLine(line string) string {
	// var output strings.Builder
	// latestConsumedIndex := 0
	// indentLevel := 0
	// if allSubmatchIndices := regexLineWithListedFields.FindAllStringSubmatchIndex(line, -1); allSubmatchIndices != nil {
	// 	for _, submatchIndices := range allSubmatchIndices {
	// 		for i := 2; i < len(submatchIndices); i += 2 {
	// 			indexStart := submatchIndices[i]
	// 			indexEnd := submatchIndices[i+1]
	// 			if indexStart > latestConsumedIndex {
	// 				output.WriteString(line[latestConsumedIndex:indexStart])
	// 			}
	// 			latestConsumedIndex = indexEnd
	// 			output.WriteString(newLine)
	// 			output.WriteString(strings.Repeat("  ", indentLevel))
	// 			output.WriteString(line[indexStart:indexEnd])
	// 		}

	// 	}
	// }
	// return output.String()

	parsers := []Parser{
		ParseWordWithTrailingSpace,
		GetSeriesParser(ParseKeyColonValue),
	}

	var output strings.Builder

	input := line
	nodes := []Node{}
	finish := false
	for !finish {
		for _, parse := range parsers {
			node, reminder, err := parse(input)
			if err != nil {
				panic(err)
			}
			if node != nil {
				nodes = append(nodes, node)
			}
			if input == reminder {
				nodes = append(nodes, StringNode(reminder))
				finish = true
				break
			}
			input = reminder
		}
	}
	for _, node := range nodes {
		output.WriteString(node.String())
	}

	return output.String()
}

type Node interface {
	String() string
}

type StringNode string

func (n StringNode) String() string {
	return string(n)
}

// Parser takes a string at input, and consumes some part of it turning the consumed part into a Node. The reminder is returned,
// or empty string if consumed everything.
type Parser func(input string) (Node, string, error)

type Series struct {
	nodes []Node
}

func (s Series) String() string {
	var b strings.Builder
	for _, n := range s.nodes {
		b.WriteString(n.String())

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
				return series, remainder, nil
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

func (k KeyColonValue) String() string {
	return string(k)
}

func ParseKeyColonValue(input string) (Node, string, error) {
	re := regexp.MustCompile(`^(\w+:\s+\w+)(?:,\s+)?`)
	match := re.FindStringSubmatchIndex(input)
	if match == nil {
		return nil, input, nil
	}
	if len(match) != 4 {
		return nil, input, fmt.Errorf("the returned match for string %v is wrong: %v", input, match)
	}
	start, end := match[2], match[3]
	return KeyColonValue(input[start:end]), input[end:], nil
}

type WordWithTrailingSpace string

func (w WordWithTrailingSpace) String() string {
	return string(w)
}

func ParseWordWithTrailingSpace(input string) (Node, string, error) {
	re := regexp.MustCompile(`^(\S+\s)`)
	match := re.FindStringSubmatchIndex(input)
	if len(match) != 4 {
		return nil, input, fmt.Errorf("the returned match for string %v is wrong: %v", input, match)
	}
	start, end := match[2], match[3]
	return WordWithTrailingSpace(input[start:end]), input[end:], nil
}
