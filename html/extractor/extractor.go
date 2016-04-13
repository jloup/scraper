//Package extractor implements sraping functions
package extractor

import (
	"fmt"
	"strings"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

// extract from NodeData useful data and store it via an aggregator
type Extractor interface {
	Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error
}

type ExtractorInitError struct {
	What string
}

func (e ExtractorInitError) Error() string {
	return fmt.Sprintf("Extractor initialization error: %s", e.What)
}

// extract a specified attribute from a HTML tag
type Attribute struct {
	Rename string
	Attr   []byte
}

func (a Attribute) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	if a.Rename != "" {
		agg.Aggregate(a.Rename, string(node.Get(a.Attr)))
	} else {
		agg.Aggregate(string(a.Attr), string(node.Get(a.Attr)))
	}

	return nil
}

type AttributeA struct {
	Rename string
	Attr   atom.Atom
}

func (a AttributeA) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	if a.Rename != "" {
		agg.Aggregate(a.Rename, string(node.GetAtom(a.Attr)))
	} else {
		agg.Aggregate(a.Attr.String(), string(node.GetAtom(a.Attr)))
	}

	return nil
}

func NewAttribute(config map[string]string) (Extractor, error) {
	if config["attr"] == "" {
		return nil, ExtractorInitError{What: "Missing attr key in config"}
	}

	if a := atom.Lookup([]byte(config["attr"])); a == 0 {
		return Attribute{Attr: []byte(config["attr"]), Rename: config["rename"]}, nil
	} else {
		return AttributeA{Attr: a, Rename: config["rename"]}, nil
	}
}

// extract text content contained inside a node ans store it as 'key'
type TextContent struct {
	Key string
}

func (s TextContent) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	agg.Aggregate(s.Key, string(node.TextContent))

	return nil
}

func NewTextContent(config map[string]string) (Extractor, error) {

	if config["key"] == "" {
		return nil, ExtractorInitError{What: "Missing key in config"}
	}

	e := TextContent{Key: config["key"]}

	return e, nil
}

// extract strings.TrimSpace(text content contained inside a node) and store it as 'key'.
type TextContentStripNewLine struct {
	Key string
}

func (s TextContentStripNewLine) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	str := strings.TrimSpace(string(node.TextContent))
	agg.Aggregate(s.Key, str)

	return nil
}

func NewTextContentStripNewLine(config map[string]string) (Extractor, error) {
	if config["key"] == "" {
		return nil, ExtractorInitError{What: "Missing key in config"}
	}

	e := TextContentStripNewLine{Key: config["key"]}

	return e, nil
}

//add {<key>: <value>} to aggregat
type SetKV struct {
	Key   string
	Value string
}

func (s SetKV) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	agg.Aggregate(s.Key, s.Value)

	return nil
}

func NewSetKV(config map[string]string) (Extractor, error) {

	if config["key"] == "" {
		return nil, ExtractorInitError{What: "Missing key in config"}
	}

	if config["value"] == "" {
		return nil, ExtractorInitError{What: "Missing value in config"}
	}

	e := SetKV{Key: config["key"], Value: config["value"]}

	return e, nil
}

func NewSetType(config map[string]string) (Extractor, error) {
	config["key"] = "type"
	config["value"] = config["type"]
	return NewSetKV(config)
}
