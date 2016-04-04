package js

import (
	"fmt"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/js/extractor"
	"github.com/jloup/scraper/js/nodedata"
	"github.com/jloup/scraper/node"
)

// extractor constructor prototype
type ExtractorGenerator func(config map[string]string) (extractor.Extractor, error)
type ScraperGenerator func(config map[string]string) ([]*node.ScraperNode, error)

var extractorRegistry = map[string]ExtractorGenerator{
	"identifier": extractor.NewIdentifier,
	"literal":    extractor.NewLiteral,
	"tliteral":   extractor.NewTypedLiteral,
}

var scraperRegistry = map[string]ScraperGenerator{
	"{property:}": objectPropertyExtractorFn,
	"{}->map":     objectExtractorFn,
}

func objectPropertyExtractorFn(config map[string]string) ([]*node.ScraperNode, error) {
	if _, ok := config["property"]; !ok {
		return nil, fmt.Errorf("ObjectPropertyExtractorFn: missing 'property' key in config map")
	}

	property := node.NewScraperNode(aggregator.NewAggregatorN(), NewJsAstNode(nodedata.AllAstNodeType, config["property"]))
	literalAst := NewJsAstNode(nodedata.LeafLiteral, "")
	literalAst.Extractors = append(literalAst.Extractors, extractor.Literal{config["property"]})
	literal := node.NewScraperNode(aggregator.NewAggregatorN(), literalAst)

	property.AddChild(&literal)

	return []*node.ScraperNode{&property}, nil
}

func objectExtractorFn(config map[string]string) ([]*node.ScraperNode, error) {

	obj := node.NewScraperNode(aggregator.NewAggregator1(), NewJsAstNode(nodedata.ObjectLiteral, ""))

	for key, _ := range config {
		literalAst := NewJsAstNode(nodedata.LeafLiteral, "")
		literalAst.Extractors = append(literalAst.Extractors, extractor.Literal{key})
		literal := node.NewScraperNode(aggregator.NewAggregatorN(), literalAst)

		property := node.NewScraperNode(aggregator.NewAggregatorN(), NewJsAstNode(nodedata.AllAstNodeType, key))

		property.AddChild(&literal)
		obj.AddChild(&property)
	}

	return []*node.ScraperNode{&obj}, nil
}

type generatorError struct {
	What string
}

func (e generatorError) Error() string {
	return fmt.Sprintf("Registry : %v", e.What)
}

// register an extractor generator into registry. The associated extractor can be then used in JSON config files
func AddExtractorGenerator(name string, e ExtractorGenerator) error {
	if _, ok := extractorRegistry[name]; ok {
		return generatorError{What: fmt.Sprintf("an extractor already has this name '%v'", name)}
	}

	extractorRegistry[name] = e

	return nil

}

// register a scraper into registry. Useful for common or complex scraper that needs to be duplicated
func AddScraperConfig(name string, s ScraperGenerator) error {
	if _, ok := scraperRegistry[name]; ok {
		return generatorError{What: fmt.Sprintf("a scraper config already has this name '%v'", name)}
	}

	scraperRegistry[name] = s

	return nil
}

func GetScraperConfig(name string, config map[string]string) ([]*node.ScraperNode, error) {
	if _, ok := scraperRegistry[name]; !ok {
		return nil, generatorError{What: fmt.Sprintf("no scraper config found for '%v'", name)}
	}

	if config == nil {
		config = make(map[string]string)
	}
	scrapers, err := scraperRegistry[name](config)
	if err != nil {
		return nil, err
	}

	length := len(scrapers)
	scrapersOut := make([]*node.ScraperNode, length, length)

	for i, scraper := range scrapers {
		s := scraper.Copy()
		scrapersOut[i] = s
	}

	return scrapersOut, nil
}
