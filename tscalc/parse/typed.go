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
	return fmt.Sprintf("(%s):%s", n.Node, n.Type)
}

// Recursively flatten SequenceNode and TypedNode with SequenceNode, leaving only one top-level TypedNode with SequenceNode.
// Mind it removes the type from TypeNode containing SequenceNode.
func FlattenTyped(node Node) Node {
	if typed, ok := node.(TypedNode); ok {
		node := FlattenTyped(typed.Node)
		if seq, ok := node.(SequenceNode); ok {
			flat := []Node{}
			flatten := FlattenTyped(seq)
			if seq2, ok := flatten.(SequenceNode); ok {
				flat = append(flat, seq2.Nodes...)
			} else {
				flat = append(flat, flatten)
			}
			return SequenceNode{
				Nodes:  flat,
				cursor: seq.Cursor(),
			}
		} else if typed2, ok := node.(TypedNode); ok {
			return typed2
		} else {
			return TypedNode{
				Node: node,
				Type: typed.Type,
			}
		}
	} else if seq, ok := node.(SequenceNode); ok {
		flat := []Node{}
		for _, el := range seq.Nodes {
			el := FlattenTyped(el)
			if seq2, ok := el.(SequenceNode); ok {
				flat = append(flat, seq2.Nodes...)
			} else {
				flat = append(flat, el)
			}
		}
		return SequenceNode{
			Nodes:  flat,
			cursor: seq.Cursor(),
		}
	} else {
		return node
	}
}
