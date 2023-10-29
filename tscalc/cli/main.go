package main

import (
	"bufio"
	"flag"
	"fmt"
	"lib/tscalc/parse"
	"log"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: paste the timestamp or operations on timestamp at the input. If there is no operation, the timestamp will be converted between epoch seconds and UTC time.`)
		flag.PrintDefaults()
	}
	flag.Parse()
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
		return n.FormatISO(), nil
	case parse.IsoTimeNode:
		return n.FormatTimestamp(), nil
	}

	// When at the input there are more values, then perform the proper calculations.
	reduced, err := reduceTree(root)
	if err != nil {
		return "", err
	}

	node, ok := reduced.(parse.EpochTimeNode)
	if !ok {
		return "", fmt.Errorf("BUG! After reduction expected other node type, got %T: %v", reduced, reduced)
	}

	return node.FormatISO(), nil
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
	case parse.EpochTimeNode, parse.PeriodNode:
		return node, nil
	case parse.IsoTimeNode:
		return node.ToEpochTimeNode(), nil
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
	default:
		return nil, fmt.Errorf("BUG! Unexpected node type in reduceTree %T: %v", root, root)
	}
}

func addNodes(leftNode, rightNode parse.Node) (parse.Node, error) {
	left, ok := leftNode.(parse.EpochTimeNode)
	if !ok {
		return nil, fmt.Errorf("BUG! Expected the left addition node to be epoch time node, was %T: %v", leftNode, leftNode)
	}
	right, ok := rightNode.(parse.PeriodNode)
	if !ok {
		return nil, fmt.Errorf("BUG! Expected the right addition node to be period node, was %T: %v", rightNode, rightNode)
	}
	return parse.EpochTimeNode(float64(left) + right.Seconds), nil
}

//func calcuate(expr string) (string, error) {
//	reMain := regexp.MustCompile(`([\d:]+)\s+-\s+([\d:]+)`)
//	matchMain := reMain.FindStringSubmatch(expr)
//	if !(len(matchMain) == 3 && matchMain[0] == expr) {
//		return "", fmt.Errorf("not a correct input expression: `%s`", expr)
//	}
//	timeLeft, err := parseTime(matchMain[1])
//	if err != nil {
//		return "", err
//	}
//	timeRight, err := parseTime(matchMain[2])
//	if err != nil {
//		return "", err
//	}
//	d := timeLeft.minutes - timeRight.minutes
//	return fmt.Sprintf("%d", d), nil
//}
//
//func parseTime(timeStr string) (time, error) {
//	errMessage := fmt.Errorf("not a correct time: `%s`", timeStr)
//	re := regexp.MustCompile(`(\d+):(\d+)`)
//	m := re.FindStringSubmatch(timeStr)
//	if !(len(m) == 3 && m[0] == timeStr) {
//		return time{}, errMessage
//	}
//	hh, err := strconv.Atoi(m[1])
//	if err != nil {
//		return time{}, errMessage
//	}
//	mm, err := strconv.Atoi(m[2])
//	if err != nil {
//		return time{}, errMessage
//	}
//	return time{
//		minutes: hh*60 + mm,
//	}, nil
//}
//
//type time struct {
//	minutes int
//}
//
//func subMinutes(left, right time) int {
//	return left.minutes - right.minutes
//}
