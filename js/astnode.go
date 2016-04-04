package js

import (
	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/js/extractor"
	"github.com/jloup/scraper/js/nodedata"
	"github.com/jloup/scraper/node"
)

type JsAstNode struct {
	Type       nodedata.AstNodeType
	Identifier string
	Extractors []extractor.Extractor

	input *nodedata.NodeData
}

func NewJsAstNode(t nodedata.AstNodeType, identifier string) *JsAstNode {
	return &JsAstNode{Type: t, Identifier: identifier}
}

func (j JsAstNode) Copy() node.Node {
	return &JsAstNode{
		Type:       j.Type,
		Identifier: j.Identifier,
		Extractors: j.Extractors,
	}
}

func (j *JsAstNode) InputNode(data interface{}) {
	j.input = data.(*nodedata.NodeData)
}

func (j *JsAstNode) ValidateNode() bool {
	if j.Type == nodedata.LeafLiteral && nodedata.IsLeafLiteral(j.input.Type) {
		return true
	}

	if (j.Type != nodedata.AllAstNodeType && j.input.Type != j.Type) ||
		(j.Identifier != "" && j.Identifier != j.input.Identifier) {
		return false
	}

	return true
}

func (j *JsAstNode) ExtractNode(agg aggregator.Aggregator) error {
	for _, extractor := range j.Extractors {
		if err := extractor.Extract(j.input, agg); err != nil {
			return err
		}
	}

	return nil
}
