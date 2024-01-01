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
			fmt.Printf("%s\n", p.IsoTimeNode{Time: nowFunc()})
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
			fmt.Println(err)
			if nerr, ok := err.(p.NodeError); ok {
				fmt.Println(nerr.Cursor().Input)
				r := nerr.Cursor().Pos - 1
				if r <= 0 {
					r = 0
				}
				fmt.Printf("%s^\n", strings.Repeat("_", r))
			}
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

	// If there is a single element at the input, just convert the format.
	if seq, ok := root.(p.SequenceNode); ok {
		if nonEmpty := seq.RemoveEmpty(); nonEmpty.Len() == 1 {
			switch n := nonEmpty.Nodes[0].(type) {
			case p.EpochTimeNode:
				return fmt.Sprint(n.ToIsoTimeNode()), nil
			case p.IsoTimeNode:
				return fmt.Sprint(n.ToEpochTimeNode()), nil
			case p.LiteralNode:
				if n.Literal == strNow {
					return fmt.Sprint(p.IsoTimeNode{Time: nowFunc()}), nil
				}
			}
		}
	}

	seq := root.(p.SequenceNode)
	if seq.Len() != 2 {
		return "", fmt.Errorf("expected sequence of 2 elements, got: %s", seq)
	}

	acc := seq.Nodes[0]
	// Initial acc can be either term or [+=] period (a sequence). Here make it a single term.
	if seq, ok := acc.(p.SequenceNode); ok {
		if seq.Len() == 2 {
			literal := seq.Nodes[0].(p.LiteralNode)
			if literal.Literal == strMinus {
				acc = p.PeriodNode{Duration: -1 * seq.Nodes[1].(p.PeriodNode).Duration}
			} else {
				acc = seq.Nodes[1]
			}
		}
	}

	reduced, err := reduce(acc, seq.Nodes[1], nowFunc())

	// When at the input there are more values, then perform the proper calculations.
	//reduced, err := reduce(root, nowFunc())
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
	root, rest, err := parser.Parse(p.NewCursor(input))
	if err != nil {
		return nil, err
	}
	if !rest.Ended() {
		return nil, fmt.Errorf("failed to parse whole input, the reminder: %s", rest)
	}
	p.Logf("Input parsed to: %T %s", root, root)
	return root, nil
}

func getParser() p.Parser {
	plusMinus := p.RegexGroup(`\s*([+-])\s*`)

	term := p.FirstOf(
		p.Period,
		p.IsoTime,
		p.Literal(strNow),
		p.EpochTime,
	)
	signedTerm := p.Sequence(
		plusMinus,
		term,
	)
	syntax := p.Sequence(
		p.FirstOf(
			p.Sequence(
				plusMinus,
				p.Period,
			),
			term,
		),
		p.Optional(
			p.Repeated(signedTerm),
		),
	)
	return syntax
}

// reduce performs actual operations on nodes.
func reduce(acc p.Node, seq p.Node, now time.Time) (p.Node, error) {
	log.Printf("Reduce: %s (%T) and %s (%T)", acc, acc, seq, seq)
	for _, opTerm := range seq.(p.SequenceNode).Nodes {
		opTermSeq := opTerm.(p.SequenceNode)
		if opTermSeq.Len() != 2 {
			return nil, fmt.Errorf("expected two nodes, got %d: %s", opTermSeq.Len(), opTermSeq)
		}
		first := opTermSeq.Nodes[0]
		second := opTermSeq.Nodes[1]
		literal, ok := first.(p.LiteralNode)
		if !ok {
			return nil, fmt.Errorf("expected literal node, got %s (%T)", first, first)
		}
		combined, err := combine(acc, literal, second, nowFunc())
		if err != nil {
			return nil, err
		}
		acc = combined
	}
	return acc, nil
}

const (
	strPlus  = "+"
	strMinus = "-"
	strNow   = "now"
)

func combine(leftNode p.Node, literal p.LiteralNode, rightNode p.Node, now time.Time) (p.Node, p.NodeError) {
	log.Printf("Combine %s (%T) %s %s (%T)", leftNode, leftNode, literal, rightNode, rightNode)
	leftNode = forceIsoTime(leftNode, now)
	rightNode = forceIsoTime(rightNode, now)

	switch left := leftNode.(type) {
	case p.PeriodNode:
		switch right := rightNode.(type) {
		case p.PeriodNode:
			switch literal.Literal {
			case strPlus:
				return p.PeriodNode{
					Duration: left.Duration + right.Duration,
					Cur:      right.Cursor(),
				}, nil
			case strMinus:
				return p.PeriodNode{
					Duration: left.Duration - right.Duration,
					Cur:      right.Cursor(),
				}, nil
			}
		case p.IsoTimeNode:
			return p.IsoTimeNode{
				Time: right.Time.Add(left.Duration),
				Cur:  right.Cursor(),
			}, nil
		}
	case p.IsoTimeNode:
		switch right := rightNode.(type) {
		case p.PeriodNode:
			switch literal.Literal {
			case strPlus:
				return p.IsoTimeNode{
					Time: left.Time.Add(right.Duration),
					Cur:  right.Cursor(),
				}, nil
			case strMinus:
				return p.IsoTimeNode{
					Time: left.Time.Add(-1 * right.Duration),
					Cur:  right.Cursor(),
				}, nil
			}
		case p.IsoTimeNode:
			switch literal.Literal {
			case strMinus:
				return p.PeriodNode{
					Duration: left.Time.Sub(right.Time),
					Cur:      right.Cur,
				}, nil
			}
		}
	}
	err := combineError{
		node: leftNode,
		err:  fmt.Errorf("cannot combine %s (%T) and %s and %s (%T)", leftNode, leftNode, literal, rightNode, rightNode),
	}
	return nil, err
}

func forceIsoTime(node p.Node, now time.Time) p.Node {
	switch n := node.(type) {
	case p.EpochTimeNode:
		return n.ToIsoTimeNode()
	case p.LiteralNode:
		if n.Literal == strNow {
			return p.IsoTimeNode{Time: now, Cur: node.Cursor()}
		}
	}
	return node
}

type combineError struct {
	err  error
	node p.Node
}

func (e combineError) Error() string {
	return e.err.Error()
}

func (e combineError) Cursor() p.Cursor {
	return e.node.Cursor()
}
