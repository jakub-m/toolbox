package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type LiteralNode struct{ Exact string }

func Literal(exact string) Parser {
	pf := func(input string) (Node, string, error) {
		Logf("Literal(%s) on: %s$", exact, input)
		if rest, foundPrefix := strings.CutPrefix(input, exact); foundPrefix {
			return LiteralNode{Exact: exact}, rest, nil
		} else {
			return nil, input, nil
		}
	}
	name := fmt.Sprintf("\"%s\"", exact)
	return funcParser{pf: pf, name: name}
}

type RegexLiteralNode struct{ match string }

func RegexLiteral(pat string) Parser {
	pf := func(input string) (Node, string, error) {
		Logf("RegexLiteral(%s) on input: %s$", pat, input)
		pat := regexp.MustCompile(pat)
		indices := pat.FindStringIndex(input)
		if indices == nil {
			return nil, input, nil
		}
		match, rest := input[indices[0]:indices[1]], input[indices[1]:]
		Logf("RegexLiteral(%s) match: %s$", pat, match)
		return RegexLiteralNode{match: match}, rest, nil
	}
	name := fmt.Sprintf("/%s/", pat)
	return funcParser{pf: pf, name: name}
}
