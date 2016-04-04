package html

import (
	"strings"
	"testing"

	"github.com/jloup/scraper/node"
)

func _Display(t *testing.T, s *node.ScraperNode, depth int) {
	if s.Node == nil {
		return
	}

	htmlNode := s.Node.(*HtmlNode)
	t.Logf("addr %p %sValidators %v, Extractors %v Type %v Tag %v\n", s, strings.Repeat("\t", depth), htmlNode.Validators, htmlNode.Extractors, htmlNode.NodeType, htmlNode.TagString)

	for _, child := range s.Childs {
		_Display(t, child, depth+1)
	}
}

func TestSimpleConfig(t *testing.T) {

	scrapers, err := JSONFileToHtmlScraper("testdata/TestSimpleConfig/config.json")
	if err == nil {

		s := node.Wrap(scrapers)
		_Display(t, s, 0)

	} else {
		t.Fatalf("JSON scrapping failed %v", err)

	}
}

func TestBadConfig(t *testing.T) {
	_, err := JSONFileToHtmlScraper("testdata/TestBadConfig/config.json")
	if err == nil {
		t.Fatalf("JSON should have failed")
	} else {
		t.Logf("JSON parsing have logically failed '%v'", err)
	}
}

func Compare(o1, o2 *node.ScraperNode) bool {
	if o1 == o2 {
		return false
	}

	for i, _ := range o1.Childs {
		if Compare(o1.Childs[i], o2.Childs[i]) == false {
			return false
		}
	}

	return true
}

func TestConfigCopy(t *testing.T) {

	scrapers, err := JSONFileToHtmlScraper("testdata/TestConfigCopy/config.json")
	if err == nil {
		var copy *node.ScraperNode

		s := node.Wrap(scrapers)

		copy = s.Copy()

		t.Log("ORIGINAL:")
		_Display(t, s, 0)
		t.Log("COPY:")
		_Display(t, copy, 0)

		if Compare(copy, s) == false {
			t.Fatal("COPY FAILED")
		} else {
			t.Log("COPY OK")
		}

	} else {
		t.Fatalf("JSON scrapping failed '%v'", err)
	}
}

func TestConfigRegistry(t *testing.T) {
	scrapers, err := JSONFileToHtmlScraper("testdata/TestConfigRegistry/link.json")

	if err != nil {
		t.Fatalf("JSON scrapping failed '%v'", err)
	}

	err = AddScraperConfig("link", scrapers)
	if err != nil {
		t.Fatalf("cannot add Scraper config to registry '%v'", err)
	}

	scrapers, err = JSONFileToHtmlScraper("testdata/TestConfigRegistry/config.json")

	s := node.Wrap(scrapers)

	s.Init()

	if err != nil {
		_Display(t, s, 0)
		t.Fatalf("import of extern scraper config failed '%v'", err)
	} else {
		_Display(t, s, 0)
	}

}
