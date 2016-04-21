package extractor

import (
	"fmt"
	"regexp"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

// extract youtube video id from a youtube URL located in an attribute
type YoutubeId struct {
	Attr []byte
}

func (s YoutubeId) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	uri := node.Get(s.Attr)

	r := regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:v|embed)\/(?P<id>[\w\-]{11})(?:\?[^\'|"]*)?`)

	match := r.FindSubmatch(uri)
	if len(match) == 0 {
		return fmt.Errorf("no youtube item found in '%s'", string(uri))
	}

	for i, key := range r.SubexpNames()[1:] {
		i += 1
		agg.Aggregate(key, string(match[i]))
	}

	return nil
}

type YoutubeIdA struct {
	Attr atom.Atom
}

func (s YoutubeIdA) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	uri := node.GetAtom(s.Attr)

	r := regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:v|embed)\/(?P<id>[\w\-]{11})(?:\?[^\'|"]*)?`)

	match := r.FindSubmatch(uri)
	if len(match) == 0 {
		return fmt.Errorf("no youtube item found in '%s'", string(uri))
	}

	for i, key := range r.SubexpNames()[1:] {
		i += 1
		agg.Aggregate(key, string(match[i]))
	}

	return nil
}
func NewYoutubeId(config map[string]string) (Extractor, error) {
	if config["attr"] == "" {
		return nil, ExtractorInitError{What: "Missing attr key in config"}
	}

	if a := atom.Lookup([]byte(config["attr"])); a == 0 {
		return YoutubeId{Attr: []byte(config["attr"])}, nil
	} else {
		return YoutubeIdA{Attr: a}, nil
	}
}
