package html

import (
	"io"

	"github.com/jloup/scraper/html/nodedata"
	"github.com/jloup/scraper/node"
	html "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/net/html/charset"
)

const (
	initialNbAttr     = 5
	initialNbAtomAttr = 3
)

var autoCloseElements = [16]atom.Atom{
	atom.Br,
	atom.Img,
	atom.Link,
	atom.Meta,
	atom.Embed,
	atom.Source,
	atom.Input,
	atom.Area,
	atom.Base,
	atom.Col,
	atom.Command,
	atom.Hr,
	atom.Keygen,
	atom.Param,
	atom.Track,
	atom.Wbr,
}

func isAutoCloseElement(tag atom.Atom) bool {
	for i := 0; i < 16; i++ {
		if autoCloseElements[i] == tag {
			return true
		}
	}

	return false
}

// walk HTML tree and returns items scraped
func ScrapHTML(scrapers []*node.ScraperNode, r io.Reader, contentType string) ([]map[string]interface{}, error) {
	r, err := charset.NewReader(r, contentType)
	if err != nil {
		return nil, err
	}

	t := html.NewTokenizer(r)

	for _, s := range scrapers {
		s.Init()
	}

	arr, err := walkHTML(scrapers, t)

	return arr, err
}

func walkHTML(scrapers []*node.ScraperNode, t *html.Tokenizer) ([]map[string]interface{}, error) {
	depth := 0
	items := make([]map[string]interface{}, 0, 0)

	var input nodedata.NodeData
	var hasAttr, ok bool
	var a atom.Atom
	var key, val []byte
	var err error
	var tt html.TokenType
	var scraper *node.ScraperNode

	input.Attr = make([]nodedata.Attribute, initialNbAttr)
	input.AttrAtom = make([]nodedata.AtomAttribute, initialNbAtomAttr)

	for {
		tt = t.Next()

		if tt == html.ErrorToken {
			for _, scraper = range scrapers {
				scraper.ProcessNode(&nodedata.NodeData{}, &items, 0)
			}
			return items, nil
		}

		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			input.Type = html.ElementNode
			input.TextContent = nil
			for j := 0; j < len(input.Attr); j++ {
				if input.Attr[j].N != nil {
					input.Attr[j].N = nil
				} else {
					break
				}
			}
			for j := 0; j < len(input.AttrAtom); j++ {
				if input.AttrAtom[j].N != 0 {
					input.AttrAtom[j].N = 0
				} else {
					break
				}
			}
			input.TagString, hasAttr = t.TagName()
			input.TagAtom = atom.Lookup(input.TagString)

			if input.TagAtom != 0 && isAutoCloseElement(input.TagAtom) {
				tt = html.SelfClosingTagToken
			}

			if hasAttr == true {
				for {
					key, val, ok = t.TagAttr()
					switch a = atom.Lookup(key); a {
					case 0:
						input.Set(key, val)

					default:
						input.SetAtom(a, val)
					}

					if ok == false {
						break
					}
				}
			}

			for _, scraper = range scrapers {
				err = scraper.ProcessNode(&input, &items, depth)
				if err != nil {
					break
				}
			}

			if tt == html.StartTagToken {
				depth += 1
			}

			if err != nil {
				return nil, err
			}

		} else if tt == html.TextToken {
			input.Type = html.TextNode
			input.TagString = nil
			input.TextContent = t.Text()

			for _, scraper = range scrapers {
				err = scraper.ProcessNode(&input, &items, depth)
				if err != nil {
					break
				}
			}
			if err != nil {
				return nil, err
			}

		} else if tt == html.EndTagToken {
			input.TagString, _ = t.TagName()
			input.TagAtom = atom.Lookup(input.TagString)

			if input.TagAtom != 0 && !isAutoCloseElement(input.TagAtom) {
				depth -= 1
			}
		}
	}

}
