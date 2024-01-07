package parse

import "fmt"

// Typed explicitly assigns type (string label) of the node to the node returned by the parser.
// For example, in expression `2 / 5` both 2 and 5 can be integeres, but the first parsed node
// can be of type `nominator` and the second one of type `denominator`.
// The node returned by the original parser will be now wrapped in TypedNode type.
func Typed(parser Parser, typeName string) Parser {
	return &typedParser{
		parser:   parser,
		typeName: typeName,
	}
}

type typedParser struct {
	parser   Parser
	typeName string
}

func (p *typedParser) Parse(input Cursor) (Node, Cursor, error) {
	node, cur, err := p.parser.Parse(input)
	if node != nil {
		node = TypedNode{node, p.typeName}
	}
	return node, cur, err
}

func (p *typedParser) String() string {
	return fmt.Sprintf("%s(%s)", p.typeName, p.parser)
}

type TypedNode struct {
	Node Node
	Type string
}

func (n TypedNode) Cursor() Cursor {
	return n.Node.Cursor()
}

func (n TypedNode) String() string {
	return fmt.Sprintf("%s(%s)", n.Type, n.Node)
}
