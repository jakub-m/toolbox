package parse

import "fmt"

func Flat(parser Parser) Parser {
	return flattenParser{parser}
}

type flattenParser struct {
	parser Parser
}

func (p flattenParser) Parse(input Cursor) (Node, Cursor, error) {
	node, cur, err := p.parser.Parse(input)
	if node == nil || err != nil {
		return node, cur, err
	}
	return FlattenNode(node), cur, err
}

func (p flattenParser) String() string {
	return fmt.Sprint(p.parser)
}

// FlattenNodeTyped recursively flattens SequenceNode and TypedNode with SequenceNode, leaving only one top-level TypedNode with SequenceNode.
// Mind it removes the type from TypeNode containing SequenceNode.
//
// Deprecated. Remove?
func FlattenNodeTyped(root Node) Node {
	unpackTyped := func(node Node) any {
		if typed, ok := node.(TypedNode); ok {
			if seq, ok := typed.Node.(SequenceNode); ok {
				return seq.Nodes
			} else {
				return node
			}
		} else if seq, ok := node.(SequenceNode); ok {
			return seq.Nodes
		} else {
			return node
		}

	}
	return FlattenNodeFunc(root, unpackTyped)
}

// FlattenNode flattens sequences in the node.
// Deprecated. Remove?
func FlattenNode(node Node) Node {
	unpackSequence := func(n Node) any {
		if seq, ok := n.(SequenceNode); ok {
			return seq.Nodes
		} else {
			return n
		}
	}
	return FlattenNodeFunc(node, unpackSequence)
}

// FlattenNodeFunc is a generic function to write custom flattening functions. The function checks if the node is a sequence and if so,
// runs the fn on each element of the sequence and flattens the result.
// If fn returns nil, then the element is ignored. If fn retuns []Node, then the result is flattened (recursively).
// If fn returns a single Node then the element is used as is.
//
// If fn is an identity function, then the whole operation is no-op, since it won't recognise SequenceNode.
func FlattenNodeFunc(node Node, fn func(Node) any) Node {
	var flattenRec func(node Node) any
	flattenRec = func(node Node) any {
		if n := fn(node); n == nil {
			return nil
		} else {
			switch transformed := n.(type) {
			case Node:
				return transformed
			case []Node:
				flat := []Node{}
				for _, o := range transformed {
					if p := flattenRec(o); p != nil {
						switch q := p.(type) {
						case Node:
							flat = append(flat, q)
						case []Node:
							flat = append(flat, q...)
						default:
							panic(fmt.Sprintf("unexpected type in FlattenFunc: %T %s", p, p))
						}
					}
				}
				return flat
			default:
				panic(fmt.Sprintf("unexpected type in FlattenFunc: %T %s", n, n))
			}
		}
	}
	if f := flattenRec(node); f == nil {
		return nil
	} else {
		switch g := f.(type) {
		case Node:
			return g
		case []Node:
			return SequenceNode{
				Nodes:  g,
				cursor: node.Cursor(),
			}
		default:
			panic(fmt.Sprintf("unexpected type in FlattenFunc: %T %s", f, f))
		}
	}
}
