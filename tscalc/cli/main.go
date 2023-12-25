package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	p "lib/tscalc/parse"
	"log"
	"os"
	"strings"
	"time"
)

var nowFunc func() time.Time = time.Now

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: Paste the timestamp or operations on timestamp at the input. If there is no operation, the timestamp will be converted between epoch seconds and UTC time.\n\n`)
		flag.PrintDefaults()
	}
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()

	if !verbose {
		log.SetOutput(io.Discard)
	}

	if stat, err := os.Stdin.Stat(); err == nil {
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			// If stdin not opened, just print current time.
			log.Println("No stdin, print current time")
			fmt.Printf("%s\n", p.IsoTimeNode(nowFunc()))
			return
		}
	} else {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		res, err := handleLine(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error: %v", err)
	}

}

func handleLine(line string) (string, error) {
	line = strings.TrimSpace(line)
	root, err := parseInput(line)
	if err != nil {
		return "", err
	}

	// If the input is a plain value, then only convert the formats.
	switch n := root.(type) {
	case p.EpochTimeNode:
		return fmt.Sprint(n.ToIsoTimeNode()), nil
	case p.IsoTimeNode:
		return fmt.Sprint(n.ToEpochTimeNode()), nil
		// default:
		// 	return "", fmt.Errorf("unexpected node type: %T", n)
	}

	// When at the input there are more values, then perform the proper calculations.
	reduced, err := reduceTree(root, nowFunc())
	if err != nil {
		return "", err
	}

	// Format output
	switch n := reduced.(type) {
	case p.IsoTimeNode, p.PeriodNode:
		return fmt.Sprint(n), nil
	}
	return "", fmt.Errorf("BUG! After reduction expected other node type, got %T: %v", reduced, reduced)
}

func parseInput(input string) (p.Node, error) {
	parser := getParser()
	p.Logf("parser: %s", parser)
	root, rest, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}
	if rest != "" {
		return nil, fmt.Errorf("failed to parse whole input, the reminder: %s", rest)
	}
	p.Logf("parsed: %T %s", root, root)
	return root, nil
}

/*
getParser returns the parser of the following form:

	<syntax> ::= <date-expr> | <period-expr>
	<period-expr> ::= <period> | <period> "+" <period-expr> | <period> "-" <period-expr> | <date-expr> "-" <date-expr>
	<date-expr> ::= <date> | <date> "+" <period-expr>
*/
func getParser() p.Parser {
	periodExprRef := p.Ref()
	periodExprRef.Name = "<period-expr>"
	dateExprRef := p.Ref()
	dateExprRef.Name = "<date-expr>"

	periodExpr := p.FirstOf(
		p.Addition(p.Period, periodExprRef),
		p.Subtraction(p.Period, periodExprRef),
		p.Subtraction(dateExprRef, dateExprRef),
		p.Period,
	)
	periodExprRef.Parser = periodExpr

	time_ := p.FirstOf(
		p.IsoTime,
		p.EpochTime,
		p.Literal("NOW"),
	)

	dateExpr := p.FirstOf(
		p.Addition(time_, periodExpr),
		p.Subtraction(time_, periodExpr),
		time_,
	)
	dateExprRef.Parser = dateExpr

	syntax := p.FirstOf(
		periodExpr,
		dateExpr,
	)

	return syntax
}

// reduceTree performs actual operations on nodes.
func reduceTree(root p.Node, now time.Time) (p.Node, error) {
	switch node := root.(type) {
	case p.EpochTimeNode:
		return node.ToIsoTimeNode(), nil
	case p.IsoTimeNode, p.PeriodNode:
		return node, nil
	case p.AddNode:
		reducedLeft, err := reduceTree(node.Left, now)
		if err != nil {
			return nil, err
		}
		reducedRight, err := reduceTree(node.Right, now)
		if err != nil {
			return nil, err
		}
		return addNodes(reducedLeft, reducedRight)
	case p.SubNode:
		reducedLeft, err := reduceTree(node.Left, now)
		if err != nil {
			return nil, err
		}
		reducedRight, err := reduceTree(node.Right, now)
		if err != nil {
			return nil, err
		}
		return subNodes(reducedLeft, reducedRight)
	case p.LiteralNode:
		switch node.Exact {
		case "NOW":
			return p.IsoTimeNode(now), nil
		default:
			return nil, fmt.Errorf("BUG! Unexpected literal: %s", node.Exact)
		}

	default:
		return nil, fmt.Errorf("BUG! Unexpected node type in reduceTree %T: %v", root, root)
	}
}

func addNodes(leftNode, rightNode p.Node) (p.Node, error) {
	right := rightNode.(p.PeriodNode)
	switch left := leftNode.(type) {
	case p.IsoTimeNode:
		return addNodesIsoTimePeriod(left, right)
	case p.PeriodNode:
		return addNodesPeriodPeriod(left, right)
	default:
		return nil, fmt.Errorf("BUG! Unexpected node type %T in addNodes", leftNode)
	}
}

func addNodesIsoTimePeriod(left p.IsoTimeNode, right p.PeriodNode) (p.Node, error) {
	return p.IsoTimeNode(time.Time(left).Add(time.Duration(right))), nil
}

func addNodesPeriodPeriod(left, right p.PeriodNode) (p.Node, error) {
	return p.PeriodNode(time.Duration(left) + (time.Duration(right))), nil
}

func subNodes(leftNode, rightNode p.Node) (p.Node, error) {
	switch left := leftNode.(type) {
	case p.IsoTimeNode:
		switch right := rightNode.(type) {
		case p.IsoTimeNode:
			return p.PeriodNode(time.Time(left).Sub(time.Time(right))), nil
		case p.PeriodNode:
			return p.IsoTimeNode(time.Time(left).Add(0 - time.Duration(right))), nil
		}
	case p.PeriodNode:
		switch right := rightNode.(type) {
		case p.PeriodNode:
			return p.PeriodNode(time.Duration(left) - time.Duration(right)), nil
		}
	}
	return nil, fmt.Errorf("BUG! unexpected note types in subNodes: %T and %T", leftNode, rightNode)
}
