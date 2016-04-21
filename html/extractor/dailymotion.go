package extractor

import (
	"fmt"
	"regexp"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

var dailymotionRegex = `(?:https?:\/\/)?(?:www\.)?dailymotion\.com\/embed\/video\/(?P<id>[\w\-]+)`

// extract youtube video id from a youtube URL located in an attribute
type DailymotionVideo struct {
	Attr     []byte
	AtomAttr atom.Atom
}

func (d DailymotionVideo) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	var uri []byte

	if len(d.Attr) > 0 {
		uri = node.Get(d.Attr)
	} else {
		uri = node.GetAtom(d.AtomAttr)
	}

	r := regexp.MustCompile(dailymotionRegex)

	match := r.FindSubmatch(uri)
	if len(match) == 0 {
		return fmt.Errorf("no dailymotion item found in '%s'", string(uri))
	}

	//resu := make(map[string]string)
	for i, key := range r.SubexpNames()[1:] {
		i += 1
		//resu[key] = string(match[i])
		agg.Aggregate(key, string(match[i]))
	}

	return nil
}

func NewDailymotionVideo(config map[string]string) (Extractor, error) {
	if config["attr"] == "" {
		return nil, ExtractorInitError{What: "Missing attr key in config"}
	}

	if a := atom.Lookup([]byte(config["attr"])); a == 0 {
		return DailymotionVideo{Attr: []byte(config["attr"])}, nil
	} else {
		return DailymotionVideo{AtomAttr: a}, nil
	}
}
