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
	"unicode/utf8"
)

type args struct {
	selectorSpec string
	shouldCut    bool
	debug        bool
	separator    string
	showCount    bool
}

func main() {
	opts := args{
		selectorSpec: "1-",
		shouldCut:    false,
		debug:        false,
		separator:    "",
		showCount:    false,
	}
	flag.StringVar(&opts.selectorSpec, "f", opts.selectorSpec, "Select only those fields. Example: `1-`, `-2,3,4-`.")
	flag.StringVar(&opts.separator, "d", opts.separator, "Delimiter string (not a regex). If left unset then consecutive whitespace is used.")
	flag.BoolVar(&opts.shouldCut, "x", opts.shouldCut, "Print selected values, not full lines.")
	flag.BoolVar(&opts.debug, "v", opts.debug, "Verbose debug mode.")
	flag.BoolVar(&opts.showCount, "c", opts.showCount, "Prefix lines by number of occurances.")
	flag.Usage = func() {
		fmt.Println(`Utility that combines uniq with cut. It returns unique lines but taking into account only fragments of the input.`)
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()
	mainInternal(opts, os.Stdin, os.Stdout)
}

func mainInternal(args args, in io.Reader, out io.Writer) {
	if args.selectorSpec == "" {
		fatalln("no fields selected")
	}

	setDebug(args.debug)
	stringSelectors, err := selectorSpecToStringSelectors(args.selectorSpec)

	if err != nil {
		fatalln("bad field specifier:", err)
	}

	var splitter splitterFunc

	switch n := utf8.RuneCountInString(args.separator); {
	case n == 0:
		whitespace := regexp.MustCompile(`\s+`)
		splitter = func(s string) []string {
			return whitespace.Split(s, -1)
		}
	default:
		splitter = func(s string) []string {
			return strings.Split(s, args.separator)
		}
	}

	scanner := bufio.NewScanner(in)
	selector := selector{
		prevSelected:    nil,
		stringSelectors: stringSelectors,
		showSelected:    args.shouldCut,
		showCount:       args.showCount,
		out:             out,
	}
	for scanner.Scan() {
		selector = selector.selectUnique(scanner.Text(), splitter)
	}
	selector.finish()

	if err := scanner.Err(); err != nil {
		fatalln(err)
	}
}

type splitterFunc func(string) []string

type selector struct {
	prevLine        string
	prevSelected    []string
	stringSelectors []stringSliceSelector
	showSelected    bool
	showCount       bool
	uniqLineCount   int
	out             io.Writer
}

func (s selector) selectUnique(line string, splitter splitterFunc) selector {
	parts := splitter(line)
	log.Println(line, parts)
	selected := []string{}
	for _, ssel := range s.stringSelectors {
		selected = append(selected, ssel.selectValues(parts)...)
	}

	prevDiffers := s.prevSelected != nil && !slicesEqual(s.prevSelected, selected)
	if prevDiffers {
		s.print()
		s.uniqLineCount = 1
	} else {
		s.uniqLineCount++
	}
	s.prevSelected = selected
	s.prevLine = line
	return s
}

func (s selector) finish() {
	s.print()
}

func (s selector) print() {
	textToDisplay := ""
	if s.showSelected {
		textToDisplay = strings.Join(s.prevSelected, "\t")
	} else {
		textToDisplay = s.prevLine
	}
	if s.showCount {
		textToDisplay = fmt.Sprintf("%d\t%s", s.uniqLineCount, textToDisplay)
	}
	fmt.Fprintln(s.out, textToDisplay)
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
			specToSelectExact,
			specToSelectTo,
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

type selectExact int

func specToSelectExact(spec string) (stringSliceSelector, error) {
	re := regexp.MustCompile(`^(\d+)$`)
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
	return selectExact(val), nil
}

func (s selectExact) selectValues(values []string) []string {
	i := int(s) - 1
	if len(values) <= i {
		return []string{}
	}
	return []string{values[i]}
}

var _ stringSliceSelector = (*selectExact)(nil)

type selectTo int

func specToSelectTo(spec string) (stringSliceSelector, error) {
	re := regexp.MustCompile(`^-(\d+)$`)
	submatch := re.FindStringSubmatch(spec)
	if submatch == nil {
		return selectTo(-1), fmt.Errorf("bad spec: %s", spec)
	}
	val, err := strconv.Atoi(submatch[1])
	if err != nil {
		return selectTo(-1), fmt.Errorf("bad spec, not a number: %s", spec)
	}
	if val < 1 {
		return selectTo(-1), fmt.Errorf("bad spec, must be at least 1: %s", spec)
	}
	return selectTo(val), nil
}

func (s selectTo) selectValues(values []string) []string {
	if s < 1 {
		return []string{}
	}
	if len(values) < int(s) {
		return values
	}
	return values[:s]
}

var _ stringSliceSelector = (*selectTo)(nil)

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
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func fatalln(message ...any) {
	// don't use log.Fatal because it would display only in debug mode.
	fmt.Fprintln(os.Stderr, message...)
	os.Exit(1)
}
