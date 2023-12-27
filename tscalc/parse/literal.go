package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type LiteralNode string

func Literal(exact string) Parser {
	pf := func(input Cursor) (Node, Cursor, error) {
		Logf("Literal(%s) on: %s$", exact, input)
		if foundPrefix := strings.HasPrefix(input.String(), exact); foundPrefix {
			return LiteralNode(exact), input.Advance(len(exact)), nil
		} else {
			return nil, input, nil
		}
	}
	name := fmt.Sprintf("\"%s\"", exact)
	return FuncParser{Fn: pf, Name: name}
}

func Regex(pat string) Parser {
	return getRegexGroupParser(pat, 0)
}

func RegexGroup(pat string) Parser {
	return getRegexGroupParser(pat, 1)
}

func getRegexGroupParser(pat string, group int) Parser {
	if !strings.HasPrefix(pat, "^") {
		pat = "^" + pat
	}
	pf := func(input Cursor) (Node, Cursor, error) {
		Logf("Regex(%s) on input: %s$", pat, input)
		pat := regexp.MustCompile(pat)
		submatches := pat.FindStringSubmatchIndex(input.String())
		if submatches == nil {
			return nil, input, nil
		}
		k := group * 2
		match := input.String()[submatches[k]:submatches[k+1]]
		Logf("Regex(%s) match: %s$", pat, match)
		return LiteralNode(match), input.Advance(submatches[1]), nil
	}
	name := fmt.Sprintf("/%s/", pat)
	return FuncParser{Fn: pf, Name: name}
}
