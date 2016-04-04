package js

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/js/nodedata"
	"github.com/jloup/scraper/node"
)

type configNode struct {
	Config     map[string]string
	Type       string
	Identifier string
	Agg        string
	Extractors []map[string]string

	Childs []configNode
}

type configRoot struct {
	Config []configNode
}

func stringToAstNodeType(str string) nodedata.AstNodeType {
	switch str {
	case "{}":
		return nodedata.ObjectLiteral
	case "[]":
		return nodedata.ArrayLiteral
	case "fn":
		return nodedata.FunctionLiteral
	case "var":
		return nodedata.VariableExpression
	case "literal":
		return nodedata.LeafLiteral
	case "string":
		return nodedata.StringLiteral
	case "bool":
		return nodedata.BooleanLiteral
	case "number":
		return nodedata.NumberLiteral
	case "for":
		return nodedata.ForStatement
	case "new":
		return nodedata.NewExpression
	case "*":
		return nodedata.AllAstNodeType
	}

	return nodedata.UnknownAstNodeType
}

func JSONToJsScraper(r io.Reader) ([]*node.ScraperNode, error) {
	var config configRoot

	dec := json.NewDecoder(r)

	err := dec.Decode(&config)
	if err != nil {
		return nil, err
	}

	var scrapers []*node.ScraperNode
	for _, child := range config.Config {
		scraper, err := parseConfig(child)
		if err != nil {
			return nil, err
		}

		scrapers = append(scrapers, scraper...)
	}

	return scrapers, err
}

func resolveExtractors(config configNode, jsAstNode *JsAstNode) error {

	for _, extractorConfig := range config.Extractors {
		if extractorRegistry[extractorConfig["name"]] != nil {
			generator := extractorRegistry[extractorConfig["name"]]
			if extractor, err := generator(extractorConfig); err == nil {
				jsAstNode.Extractors = append(jsAstNode.Extractors, extractor)
			} else {
				return err
			}

		} else {
			return generatorError{fmt.Sprintf("extractor '%s' not found", extractorConfig["name"])}
		}
	}

	return nil
}

func configToJsAstNode(config configNode) (*node.ScraperNode, error) {
	nodeType := stringToAstNodeType(config.Type)
	if nodeType == nodedata.UnknownAstNodeType {
		return nil, fmt.Errorf("not recignized node type '%v'", config.Type)
	}

	jsAstNode := NewJsAstNode(nodeType, config.Identifier)

	agg := aggregator.NewAggregatorFromConfig(config.Agg)

	if len(config.Extractors) > 0 {
		if err := resolveExtractors(config, jsAstNode); err != nil {
			return nil, err
		}
	}

	scraperNode := node.NewScraperNode(agg, jsAstNode)

	return &scraperNode, nil
}

func parseConfig(config configNode) ([]*node.ScraperNode, error) {
	var scrapers []*node.ScraperNode
	var err error

	if stringToAstNodeType(config.Type) == nodedata.UnknownAstNodeType {
		scrapers, err = GetScraperConfig(config.Type, config.Config)
		if err != nil {
			return nil, err
		}
	} else {
		scrapers = make([]*node.ScraperNode, 1, 1)
		scraper, err := configToJsAstNode(config)
		if err != nil {
			return nil, err
		}

		scrapers[0] = scraper

		for _, child := range config.Childs {
			newJsAstNode, err := parseConfig(child)
			if err != nil {
				return scrapers, err
			}
			scraper.AddChilds(newJsAstNode...)

		}
	}

	return scrapers, nil
}
