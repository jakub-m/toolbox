package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	fields := ""
	shouldCut := false
	flag.StringVar(&fields, "f", "", "select only those fields")
	flag.BoolVar(&shouldCut, "x", false, "print selected values, not full lines")
	flag.Parse()

	if fields == "" {
		log.Fatalln("no fields selected")
	}

	scanner := bufio.NewScanner(os.Stdin)
	selector := selector{
		previous: nil,
		stringSelectors: []stringSliceSelector{
			selectFrom(1),
		},
		shouldCut: shouldCut,
	}
	for scanner.Scan() {
		selector = selector.selectUnique(scanner.Text(), os.Stdout)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

type selector struct {
	previous        []string
	stringSelectors []stringSliceSelector
	shouldCut       bool
}

var reSpaces = regexp.MustCompile(`\s+`)

func (s selector) selectUnique(text string, w io.Writer) selector {
	parts := reSpaces.Split(text, -1)
	log.Println(text, parts)
	selected := []string{}
	for _, ssel := range s.stringSelectors {
		selected = append(selected, ssel.selectValues(parts)...)
	}

	if s.previous == nil || !slicesEqual(s.previous, selected) {
		s.previous = selected
		if s.shouldCut {
			fmt.Fprintln(w, strings.Join(selected, "\t"))
		} else {
			fmt.Fprintln(w, text)
		}
	}
	return s
}

type stringSliceSelector interface {
	selectValues([]string) []string
}

type selectFrom int

func (s selectFrom) selectValues(values []string) []string {
	if len(values) < int(s) {
		return []string{}
	}
	return values[s:]
}

var _ stringSliceSelector = (*selectFrom)(nil)

func slicesEqual[T comparable](some, other []T) bool {
	if len(some) != len(other) {
		return false
	}
	for i := 0; i < len(some); i++ {
		if some[i] != other[i] {
			return false
		}
	}
	return true
}
