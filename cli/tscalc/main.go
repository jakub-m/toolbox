package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
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
	root, err := parse(line)
	if err != nil {
		return "", err
	}
	// if root is the only node, which it is for now...
	switch n := root.(type) {
	case EpochTimeNode:
		return n.isoFormat(), nil
	default:
		return "", fmt.Errorf("unknown type of node %T", root)
	}
}

type Node interface {
}

type EpochTimeNode float64

func (t EpochTimeNode) isoFormat() string {
	sec, frac := math.Modf(float64(t))
	return time.Unix(int64(sec), int64(1_000_000_000*frac)).UTC().Format("2006-01-02T15:04:05-07:00")
}

func parse(input string) (Node, error) {
	epochTimeNode, rest, err := parseEpochTime(input)
	if err != nil {
		return nil, err
	}
	if rest != "" {
		return nil, fmt.Errorf("failed to parse whole input, the reminder: %s", rest)
	}
	return epochTimeNode, nil
}

// parseEpochTime returns the node if found, the output string and the error if any. If the node is not found, then
// do not return error. Erorrs mean that the program should stop.
func parseEpochTime(input string) (EpochTimeNode, string, error) {
	reEpochTime := regexp.MustCompile(`\d+(\.\d+)?`)
	indices := reEpochTime.FindStringIndex(input)
	if indices == nil {
		return 0, input, nil
	}
	match := input[indices[0]:indices[1]]
	t, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, input, fmt.Errorf("error while parsing %s: %w", err)
	}
	return EpochTimeNode(t), input[indices[1]:], nil
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
