package html

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/node"
	html "golang.org/x/net/html"
)

type configNode struct {
	Extern     string
	Tag        string
	Childs     []configNode
	Validators []map[string]string
	Extractors []map[string]string
	Type       string
}

type configRoot struct {
	Config []configNode
}

// create a scraper from a config in JSON format
func JSONToHtmlScraper(r io.Reader) ([]*node.ScraperNode, error) {
	var config configRoot

	dec := json.NewDecoder(r)

	err := dec.Decode(&config)
	if err != nil {
		return nil, err
	}

	scrapers := make([]*node.ScraperNode, 0, 0)

	for _, child := range config.Config {
		scraper, err := parseConfig(child)
		if err != nil {
			return nil, err
		}
		scrapers = append(scrapers, scraper...)
	}

	return scrapers, err
}

// create a scraper from a JSON config file
func JSONFileToHtmlScraper(filepath string) ([]*node.ScraperNode, error) {
	f, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	return JSONToHtmlScraper(f)
}

func resolveValidators(config configNode, htmlNode *HtmlNode) error {

	for _, validatorConfig := range config.Validators {
		if validatorRegistry[validatorConfig["name"]] != nil {
			generator := validatorRegistry[validatorConfig["name"]]
			if validator, err := generator(validatorConfig); err == nil {
				htmlNode.Validators = append(htmlNode.Validators, validator)
			} else {
				return err
			}

		} else {
			return generatorError{fmt.Sprintf("validator '%s' not found", validatorConfig["name"])}
		}
	}

	return nil
}

func resolveExtractors(config configNode, htmlNode *HtmlNode) error {

	for _, extractorConfig := range config.Extractors {
		if extractorRegistry[extractorConfig["name"]] != nil {
			generator := extractorRegistry[extractorConfig["name"]]
			if extractor, err := generator(extractorConfig); err == nil {
				htmlNode.Extractors = append(htmlNode.Extractors, extractor)
			} else {
				return err
			}

		} else {
			return generatorError{fmt.Sprintf("extractor '%s' not found", extractorConfig["name"])}
		}
	}

	return nil
}

func configToHtmlNode(config configNode) (node.ScraperNode, error) {
	var nodeType html.NodeType
	if config.Tag == "text" {
		nodeType = html.TextNode
	} else {
		nodeType = html.ElementNode
	}

	var tag []byte

	if config.Tag == "" {
		tag = nil
	} else {
		tag = []byte(config.Tag)
	}

	htmlNode := NewHtmlNode(tag, nodeType)

	agg := aggregator.NewAggregatorFromConfig(config.Type)

	if len(config.Validators) > 0 {
		if err := resolveValidators(config, htmlNode); err != nil {
			return node.ScraperNode{}, err
		}
	}

	if len(config.Extractors) > 0 {
		if err := resolveExtractors(config, htmlNode); err != nil {
			return node.ScraperNode{}, err
		}
	}

	return node.NewScraperNode(agg, htmlNode), nil
}

func parseConfig(config configNode) ([]*node.ScraperNode, error) {
	var scrapers []*node.ScraperNode
	var err error

	if config.Extern != "" {
		scrapers, err = getScraperConfig(config.Extern)
		if err != nil {
			return nil, err
		}

	} else {
		scrapers = make([]*node.ScraperNode, 1, 1)
		scraper, err := configToHtmlNode(config)
		if err != nil {
			return nil, err
		}

		scrapers[0] = &scraper

		for _, child := range config.Childs {
			newHtmlNode, err := parseConfig(child)
			if err != nil {
				return scrapers, err
			}
			scraper.AddChilds(newHtmlNode...)
		}
	}

	return scrapers, nil
}
