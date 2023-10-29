package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"lib/tscalc/parse"
	"log"
	"os"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: paste the timestamp or operations on timestamp at the input. If there is no operation, the timestamp will be converted between epoch seconds and UTC time.`)
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
			fmt.Printf("%s\n", parse.IsoTimeNode(time.Now()))
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
	root, err := parseInput(line)
	if err != nil {
		return "", err
	}

	// If the input is a plain value, then only convert the formats.
	switch n := root.(type) {
	case parse.EpochTimeNode:
		return fmt.Sprint(n.ToIsoTimeNode()), nil
	case parse.IsoTimeNode:
		return fmt.Sprint(n.ToEpochTimeNode()), nil
	}

	// When at the input there are more values, then perform the proper calculations.
	reduced, err := reduceTree(root)
	if err != nil {
		return "", err
	}

	// Format output
	switch n := reduced.(type) {
	case parse.IsoTimeNode, parse.PeriodNode:
		return fmt.Sprint(n), nil
	}

	return "", fmt.Errorf("BUG! After reduction expected other node type, got %T: %v", reduced, reduced)
}

func parseInput(input string) (parse.Node, error) {
	root, rest, err := parse.GetParser()(input)
	if err != nil {
		return nil, err
	}
	if rest != "" {
		return nil, fmt.Errorf("failed to parse whole input, the reminder: %s", rest)
	}
	return root, nil
}

// reduceTree performs actual operations on nodes.
func reduceTree(root parse.Node) (parse.Node, error) {
	switch node := root.(type) {
	case parse.EpochTimeNode:
		return node.ToIsoTimeNode(), nil
	case parse.IsoTimeNode, parse.PeriodNode:
		return node, nil
	case parse.AddNode:
		reducedLeft, err := reduceTree(node.Left)
		if err != nil {
			return nil, err
		}
		reducedRight, err := reduceTree(node.Right)
		if err != nil {
			return nil, err
		}
		return addNodes(reducedLeft, reducedRight)
	case parse.SubNode:
		reducedLeft, err := reduceTree(node.Left)
		if err != nil {
			return nil, err
		}
		reducedRight, err := reduceTree(node.Right)
		if err != nil {
			return nil, err
		}
		return subNodes(reducedLeft, reducedRight)
	default:
		return nil, fmt.Errorf("BUG! Unexpected node type in reduceTree %T: %v", root, root)
	}
}

func addNodes(leftNode, rightNode parse.Node) (parse.Node, error) {
	left := leftNode.(parse.IsoTimeNode)
	right := rightNode.(parse.PeriodNode)
	return parse.IsoTimeNode(time.Time(left).Add(time.Duration(right))), nil
}

func subNodes(leftNode, rightNode parse.Node) (parse.Node, error) {
	left := leftNode.(parse.IsoTimeNode)
	right := rightNode.(parse.IsoTimeNode)
	return parse.PeriodNode(time.Time(left).Sub(time.Time(right))), nil
}
