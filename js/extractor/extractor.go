package extractor

import (
	"fmt"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/js/nodedata"
)

type Extractor interface {
	Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error
}

type ExtractorInitError struct {
	What string
}

func (e ExtractorInitError) Error() string {
	return fmt.Sprintf("Extractor initialization error: %s", e.What)
}

type Identifier struct {
	Key string
}

func (i Identifier) Extract(n *nodedata.NodeData, agg aggregator.Aggregator) error {
	agg.Aggregate(i.Key, n.Identifier)

	return nil
}

func NewIdentifier(config map[string]string) (Extractor, error) {
	if config["key"] == "" {
		return nil, ExtractorInitError{What: "missing key in config"}
	}

	return Identifier{config["key"]}, nil
}

type Literal struct {
	Key string
}

func (l Literal) Extract(n *nodedata.NodeData, agg aggregator.Aggregator) error {
	agg.Aggregate(l.Key, n.Content)

	return nil
}

func NewLiteral(config map[string]string) (Extractor, error) {
	if config["key"] == "" {
		return nil, ExtractorInitError{What: "missing key in config"}
	}

	return Literal{config["key"]}, nil
}

func stringToLiteralType(str string) nodedata.AstNodeType {
	switch str {
	case "number":
		return nodedata.NumberLiteral
	case "string":
		return nodedata.StringLiteral
	case "null":
		return nodedata.NullLiteral
	case "boolean":
		return nodedata.BooleanLiteral
	}

	return nodedata.UnknownAstNodeType
}

type TypedLiteral struct {
	Type nodedata.AstNodeType
	Literal
}

func (t TypedLiteral) Extract(n *nodedata.NodeData, agg aggregator.Aggregator) error {
	if n.Type == t.Type {
		t.Literal.Extract(n, agg)
	}

	return nil
}

func NewTypedLiteral(config map[string]string) (Extractor, error) {
	if config["key"] == "" {
		return nil, ExtractorInitError{What: "missing key in config"}
	}

	if config["t"] == "" {
		return nil, ExtractorInitError{What: "missing t key in config"}
	}

	if stringToLiteralType(config["t"]) == nodedata.UnknownAstNodeType {
		return nil, ExtractorInitError{What: fmt.Sprintf("key 't' is not a valid js literal type '%v'", config["t"])}
	}

	return TypedLiteral{Type: stringToLiteralType(config["t"]), Literal: Literal{config["key"]}}, nil
}
