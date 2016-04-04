package extractor

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/jloup/scraper/aggregator"
	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

const (
	soundcloudRegex    = `https?:\/\/(?:www.)?api.soundcloud.com\/(?P<type>playlists|tracks|users)\/(?P<id>[0-9]+)`
	soundcloudUrlRegex = `(?P<url>https?:\/\/(?:www.)?(?:api\.)?soundcloud.com\/.+)`
)

// extract soundcloud song type and id from a soundcloud URL located in an attribute
type SoundcloudR struct {
	Attr []byte
}

func tryRegex(rgx, s string, agg aggregator.Aggregator) bool {
	r := regexp.MustCompile(rgx)
	match := r.FindStringSubmatch(s)
	if len(match) == 0 {
		return false
	}

	for i, key := range r.SubexpNames()[1:] {
		i += 1
		agg.Aggregate(key, match[i])
	}

	return true
}

func (s SoundcloudR) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	uri, err := url.Parse(string(node.Get(s.Attr)))
	if err != nil {
		return err
	}

	l := uri.Query()["url"][0]

	lurl, err := url.Parse(l)
	if err != nil {
		return err
	}

	if r, ok := lurl.Query()["secret_token"]; ok {
		agg.Aggregate("secretToken", r[0])
	}

	if r, ok := uri.Query()["secret_token"]; ok {
		agg.Aggregate("secretToken", r[0])
	}

	if tryRegex(soundcloudRegex, l, agg) {
		return nil
	}

	if !tryRegex(soundcloudUrlRegex, l, agg) {
		return fmt.Errorf("no soundcloud item has been found in '%s'", uri.String())
	}

	return nil
}

type SoundcloudRA struct {
	Attr atom.Atom
}

func (s SoundcloudRA) Extract(node *nodedata.NodeData, agg aggregator.Aggregator) error {
	uri, err := url.Parse(string(node.GetAtom(s.Attr)))
	if err != nil {
		return err
	}

	l := uri.Query()["url"][0]

	lurl, err := url.Parse(l)
	if err != nil {
		return err
	}

	if r, ok := lurl.Query()["secret_token"]; ok {
		agg.Aggregate("secretToken", r[0])
	}

	if r, ok := uri.Query()["secret_token"]; ok {
		agg.Aggregate("secretToken", r[0])
	}

	if tryRegex(soundcloudRegex, l, agg) {
		return nil
	}

	if !tryRegex(soundcloudUrlRegex, l, agg) {
		return fmt.Errorf("no soundcloud item has been found in '%s'", uri.String())
	}

	return nil
}

func NewSoundcloudR(config map[string]string) (Extractor, error) {
	if config["attr"] == "" {
		return nil, ExtractorInitError{What: "Missing attr key in config"}
	}

	if a := atom.Lookup([]byte(config["attr"])); a == 0 {
		return SoundcloudR{Attr: []byte(config["attr"])}, nil
	} else {
		return SoundcloudRA{Attr: a}, nil
	}

}
