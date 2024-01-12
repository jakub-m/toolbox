package parse

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequenceFlatten(t *testing.T) {
	input := SequenceNode{
		Nodes: []Node{
			LiteralNode{
				Literal: "a",
			},
			SequenceNode{
				Nodes: []Node{
					LiteralNode{
						Literal: "b",
					},
					SequenceNode{
						Nodes: []Node{
							LiteralNode{
								Literal: "c",
							},
							LiteralNode{
								Literal: "d",
							},
							SequenceNode{
								Nodes: []Node{
									LiteralNode{
										Literal: "e",
									},
								},
							},
						},
					},
				},
			},
			LiteralNode{
				Literal: "f",
			},
		},
	}
	assert.Equal(t, `["a" "b" "c" "d" "e" "f"]`, fmt.Sprint(FlattenNode(input)))
}

func TestFlattenTypedLiteral(t *testing.T) {
	input := TypedNode{
		Node: newLiteralNode("lit1"),
		Type: "typ1",
	}
	flat := FlattenNodeTyped(input)
	assert.Equal(t, `("lit1"):typ1`, fmt.Sprint(flat))
}

func TestFlattenLiteral(t *testing.T) {
	input := newLiteralNode("lit1")
	flat := FlattenNodeTyped(input)
	assert.Equal(t, `"lit1"`, fmt.Sprint(flat))
}

func TestFlattenSequence(t *testing.T) {
	input := newSequenceNode(
		newLiteralNode("lit1"),
		newSequenceNode(
			newLiteralNode("lit2"),
			newLiteralNode("lit3"),
		),
		newLiteralNode("lit4"),
	)
	flat := FlattenNodeTyped(input)
	assert.Equal(t, `["lit1" "lit2" "lit3" "lit4"]`, fmt.Sprint(flat))

}

func TestFlattenTyped2(t *testing.T) {
	input := TypedNode{
		Node: newSequenceNode(
			newSequenceNode(
				newLiteralNode("lit1"),
			),
		),
		Type: "typ1",
	}
	flat := FlattenNodeTyped(input)
	assert.Equal(t, `["lit1"]`, fmt.Sprint(flat))
}

func TestFlattenTyped(t *testing.T) {
	input := TypedNode{
		Node: newSequenceNode(
			newSequenceNode(
				newLiteralNode("lit1"),
				TypedNode{
					Node: newLiteralNode("lit2"),
					Type: "type_lit2",
				},
				newSequenceNode(
					newLiteralNode("lit3"),
				),
			),
			newLiteralNode("lit4"),
			TypedNode{
				Node: newLiteralNode("lit5"),
				Type: "type_lit5",
			},
			TypedNode{
				Node: newSequenceNode(
					TypedNode{
						Node: newLiteralNode("lit6"),
						Type: "type_lit6",
					},
				),
				Type: "type_seq_lit6",
			},
		),
		Type: "type_top_seq",
	}

	flat := FlattenNodeTyped(input)
	assert.Equal(t, `["lit1" ("lit2"):type_lit2 "lit3" "lit4" ("lit5"):type_lit5 ("lit6"):type_lit6]`, fmt.Sprint(flat))
}
