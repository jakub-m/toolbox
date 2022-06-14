package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		res, err := calcuate(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error: %v", err)
	}

}

func calcuate(expr string) (string, error) {
	reMain := regexp.MustCompile(`([\d:]+)\s+-\s+([\d:]+)`)
	matchMain := reMain.FindStringSubmatch(expr)
	if !(len(matchMain) == 3 && matchMain[0] == expr) {
		return "", fmt.Errorf("not a correct input expression: `%s`", expr)
	}
	timeLeft, err := parseTime(matchMain[1])
	if err != nil {
		return "", err
	}
	timeRight, err := parseTime(matchMain[2])
	if err != nil {
		return "", err
	}
	d := timeLeft.minutes - timeRight.minutes
	return fmt.Sprintf("%d", d), nil
}

func parseTime(timeStr string) (time, error) {
	errMessage := fmt.Errorf("not a correct time: `%s`", timeStr)
	re := regexp.MustCompile(`(\d+):(\d+)`)
	m := re.FindStringSubmatch(timeStr)
	if !(len(m) == 3 && m[0] == timeStr) {
		return time{}, errMessage
	}
	hh, err := strconv.Atoi(m[1])
	if err != nil {
		return time{}, errMessage
	}
	mm, err := strconv.Atoi(m[2])
	if err != nil {
		return time{}, errMessage
	}
	return time{
		minutes: hh*60 + mm,
	}, nil
}

type time struct {
	minutes int
}

func subMinutes(left, right time) int {
	return left.minutes - right.minutes
}
