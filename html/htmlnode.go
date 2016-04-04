//Package scraper provides tools to easily scrap HTML pages
package html

import (
	"bytes"
	"fmt"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/html/extractor"
	"github.com/jloup/scraper/html/nodedata"
	"github.com/jloup/scraper/html/validator"
	"github.com/jloup/scraper/node"
	html "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type ErrorHtmlNode struct {
	What string
}

func (e ErrorHtmlNode) Error() string {
	return fmt.Sprintf("HtmlNode: %s", e.What)
}

type HtmlNode struct {
	TagString []byte
	TagAtom   atom.Atom
	NodeType  html.NodeType

	Validators []validator.Validator
	Extractors []extractor.Extractor

	input *nodedata.NodeData
}

func NewHtmlNode(tag []byte, nodeType html.NodeType) *HtmlNode {
	var s *HtmlNode

	if tag == nil {
		s = nil
	} else if at := atom.Lookup([]byte(tag)); at == 0 {
		s = &HtmlNode{TagString: tag, TagAtom: 0, NodeType: nodeType}
	} else {
		s = &HtmlNode{TagString: tag, TagAtom: at, NodeType: nodeType}
	}

	return s
}

func (h HtmlNode) Copy() node.Node {
	return &HtmlNode{
		TagString:  h.TagString,
		TagAtom:    h.TagAtom,
		NodeType:   h.NodeType,
		Validators: h.Validators,
		Extractors: h.Extractors,
	}
}

func (h *HtmlNode) InputNode(data interface{}) {
	h.input = data.(*nodedata.NodeData)
}

func (h *HtmlNode) ValidateNode() bool {
	if h.input.Type != h.NodeType ||
		(h.NodeType != html.TextNode && ((h.input.TagAtom != 0 && h.input.TagAtom != h.TagAtom) || bytes.Compare(h.input.TagString, h.TagString) != 0)) {
		return false
	}

	for _, validator := range h.Validators {
		if t := validator.Validate(h.input); t == false {
			return false
		}
	}

	return true
}

func (h *HtmlNode) ExtractNode(agg aggregator.Aggregator) error {
	for _, extractor := range h.Extractors {
		if err := extractor.Extract(h.input, agg); err != nil {
			return err
		}
	}

	return nil
}
