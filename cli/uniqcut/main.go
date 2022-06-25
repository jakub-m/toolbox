package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	selectorSpec := ""
	shouldCut := false
	debug := false
	flag.StringVar(&selectorSpec, "f", "", "select only those fields")
	flag.BoolVar(&shouldCut, "x", false, "print selected values, not full lines")
	flag.BoolVar(&debug, "v", false, "verbose debug mode")
	flag.Parse()

	if selectorSpec == "" {
		log.Fatalln("no fields selected")
	}

	setDebug(debug)
	stringSelectors, err := selectorSpecToStringSelectors(selectorSpec)
	setDebug(true)

	if err != nil {
		log.Fatalln("bad field specifier:", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	selector := selector{
		previous:        nil,
		stringSelectors: stringSelectors,
		shouldCut:       shouldCut,
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

func selectorSpecToStringSelectors(selectorSpec string) ([]stringSliceSelector, error) {
	selectors := []stringSliceSelector{}
	for _, singleSpec := range strings.Split(selectorSpec, ",") {
		somethingSelected := false
		type factory func(string) (stringSliceSelector, error)
		for _, fac := range []factory{
			specToSelectFrom,
		} {
			if selector, err := fac(singleSpec); err == nil {
				selectors = append(selectors, selector)
				somethingSelected = true
				break
			} else {
				log.Println(err)
			}
		}
		if !somethingSelected {
			return nil, fmt.Errorf("invalid field spec \"%s\"", singleSpec)
		}
	}
	return selectors, nil
}

type selectFrom int

func specToSelectFrom(spec string) (stringSliceSelector, error) {
	re := regexp.MustCompile(`^(\d+)-$`)
	submatch := re.FindStringSubmatch(spec)
	if submatch == nil {
		return selectFrom(-1), fmt.Errorf("bad spec: %s", spec)
	}
	val, err := strconv.Atoi(submatch[1])
	if err != nil {
		return selectFrom(-1), fmt.Errorf("bad spec, not a number: %s", spec)
	}
	if val < 1 {
		return selectFrom(-1), fmt.Errorf("bad spec, must be at least 1: %s", spec)
	}
	return selectFrom(val), nil
}

func (s selectFrom) selectValues(values []string) []string {
        i := int(s) - 1
	if len(values) < i {
		return []string{}
	}
	return values[i:]
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

func setDebug(enabled bool) {
	if enabled {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}
}
