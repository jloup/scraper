package extractor

import (
	"fmt"
	"net/url"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

const facebookVideoHrefRegex = `^https?:\/\/(?:www.)?facebook.com\/(?P<entity>[^\/]+)\/videos(?:\/vb\.[0-9]+)?\/(?P<id>[0-9]+)`

// extract facebook iframe video id
type FacebookVideo struct {
	Attr     []byte
	AtomAttr atom.Atom
}

func (f FacebookVideo) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	var uri *url.URL
	var err error

	if len(f.Attr) > 0 {
		uri, err = url.Parse(string(node.Get(f.Attr)))
	} else {
		uri, err = url.Parse(string(node.GetAtom(f.AtomAttr)))
	}

	if err != nil {
		return err
	}

	videoHref, ok := uri.Query()["href"]
	if ok {
		agg.Aggregate("href", videoHref[0])
	}

	if !tryRegex(facebookVideoHrefRegex, videoHref[0], agg) {
		return fmt.Errorf("no facebook item has been found in '%s'", uri.String())
	}

	return nil
}

func NewFacebookVideo(config map[string]string) (Extractor, error) {
	if config["attr"] == "" {
		return nil, ExtractorInitError{What: "Missing attr key in config"}
	}

	if a := atom.Lookup([]byte(config["attr"])); a == 0 {
		return FacebookVideo{Attr: []byte(config["attr"])}, nil
	} else {
		return FacebookVideo{AtomAttr: a}, nil
	}

}
