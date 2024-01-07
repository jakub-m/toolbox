package parse_sql

import (
	"fmt"
	p "lib/tscalc/parse"
)

func Parse(input string) (p.Node, error) {
	node, cur, err := getParser().Parse(p.NewCursor(input))
	if !cur.Ended() {
		return nil, p.NewCursorError(cur, fmt.Errorf("did not parse whole input"))
	}
	node = p.Flatten(node)
	return node, err
}

const (
	typeFromExpr     = "fromExpr"
	typeFromLit      = "fromLit"
	typeSelectExpr   = "selectExpr"
	typeSelectLit    = "selectLit"
	typeSelector     = "selector"
	typeSelectorName = "selectorName"
	typeTableName    = "tableName"
)

func getParser() p.Parser {
	whitespace := p.Optional(p.Regex(`\s+`))
	lit_select := p.Typed(
		p.Literal(`select`),
		typeSelectLit)
	selectorName := p.Typed(
		p.Regex(`[*]|[0-9]+|[a-zA-Z][a-zA-Z0-9_]*`),
		typeSelectorName,
	)
	selector := p.Typed(
		p.FirstOf(
			p.Sequence(selectorName,
				p.Repeated(
					p.Sequence(
						whitespace,
						p.Literal(","),
						whitespace,
						selectorName,
					),
				)),
			selectorName,
		),
		typeSelector)
	lit_from := p.Typed(
		p.Literal(`from`),
		typeFromLit)
	table_name := p.Typed(
		p.Regex(`[a-zA-Z][a-zA-Z0-9_]*`),
		typeTableName)
	from_expr :=
		p.Typed(
			p.Sequence(
				whitespace,
				lit_from,
				whitespace,
				table_name,
			), typeFromExpr)
	select_expr := p.Typed(
		p.Sequence(
			lit_select,
			whitespace,
			selector,
			p.Optional(
				from_expr,
			),
		), typeSelectExpr)
	syntax := select_expr
	return syntax
}
