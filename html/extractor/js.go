//Package extractor implements sraping functions
package extractor

import (
	"bytes"
	"fmt"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/html/nodedata"
	"github.com/jloup/scraper/js"
	"github.com/jloup/scraper/node"
)

// extract a specified attribute from a HTML tag
type Js struct {
	ScrapConfig string
	scrapers    []*node.ScraperNode
}

func (j Js) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {

	results, err := js.ScrapJS(j.scrapers, bytes.NewReader(node.TextContent))
	if err != nil {
		return err
	}

	jsAgg := aggregator.NewAggregatorN()

	for _, result := range results {
		for key, value := range result {
			jsAgg.Aggregate(key, value)
		}

		agg.Join(jsAgg)
	}

	return nil
}

func NewJs(config map[string]string) (Extractor, error) {
	if config["config"] == "" {
		return nil, ExtractorInitError{What: "Missing 'config' key in config"}
	}

	j := Js{ScrapConfig: config["config"]}

	scrapers, err := js.GetScraperConfig(j.ScrapConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot find js scrap config %v", err)
	}

	j.scrapers = scrapers

	return j, nil
}
