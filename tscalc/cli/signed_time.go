package main

import (
	"fmt"
	p "lib/tscalc/parse"
)

// SignedPeriodExpr parses an expression on signed periods, like `12h + 1m - 3s`
var SignedPeriodExpr = signedPeriodExprStr{}

type signedPeriodExprStr struct{}

func (p signedPeriodExprStr) String() string {
	return "<signed-period-expr>"
}

func (s signedPeriodExprStr) Parse(input p.Cursor) (p.Node, p.Cursor, error) {
	node, rest, err := p.Repeated(SignedPeriod).Parse(input)
	seq := node.(p.SequenceNode)
	if err != nil {
		return nil, rest, err
	}
	if len(seq) == 0 {
		return nil, rest, err
	}
	result := seq[0].(p.PeriodNode)
	for _, period := range seq[1:] {
		result += period.(p.PeriodNode)
	}
	return p.PeriodNode(result), rest, err
}

// SignedPeriod is a period with preceding sign, like `3m`, `+3m`, `-3m`, `- 3m` or `+ 3m`.
var SignedPeriod = signedPeriodStr{}

type signedPeriodStr struct{}

func (p signedPeriodStr) String() string {
	return "<signed-period>"
}

func (s signedPeriodStr) Parse(input p.Cursor) (p.Node, p.Cursor, error) {
	parser := p.FirstOf(
		p.Period,
		p.Sequence(
			p.Optional(p.Regex(`\s+`)),
			p.Regex(`[+-]`),
			p.Optional(p.Regex(`\s+`)),
			p.Period,
		),
	)

	p.Logf("%s is %s", s, parser)
	p.LogIndentInc()
	defer p.LogIndentDec()
	node, rest, err := parser.Parse(input)
	if err != nil || node == nil {
		return nil, rest, err
	}
	switch t := node.(type) {
	case p.PeriodNode:
		return t, rest, nil
	case p.SequenceNode:
		switch t[1].(p.LiteralNode) {
		case "+":
			return t[3].(p.PeriodNode), rest, nil
		case "-":
			return -1 * t[3].(p.PeriodNode), rest, nil
		default:
			return nil, rest, fmt.Errorf("signedPeriodStr: unexpected sequence: %s", t)
		}
	default:
		return nil, rest, fmt.Errorf("signedPeriodStr: unexpected type: %T %s", node, node)
	}

}
