package parse

func newLiteralNode(literal string) LiteralNode {
	return LiteralNode{Literal: literal}
}

func newSequenceNode(seq ...Node) SequenceNode {
	return SequenceNode{
		Nodes: seq,
	}
}
