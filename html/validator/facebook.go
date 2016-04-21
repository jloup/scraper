package validator

import (
	"bytes"
	"net/url"
	"regexp"
	"strings"

	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

const facebookHost = "facebook.com"
const facebookVideoIframeSrcRegex = `^https?:\/\/(?:www.)?facebook.com\/plugins\/video\.php`

var facebookHostB = []byte(facebookHost)

type FacebookVideoIframeFast struct {
}

func (f FacebookVideoIframeFast) Validate(node *nodedata.NodeData) bool {
	return bytes.Contains(node.GetAtom(atom.Src), facebookHostB)
}

func NewFacebookVideoIframeFast(config map[string]string) (Validator, error) {
	return FacebookVideoIframeFast{}, nil
}

type FacebookVideoIframe struct {
}

func (f FacebookVideoIframe) Validate(node *nodedata.NodeData) bool {
	url, err := url.Parse(string(node.GetAtom(atom.Src)))
	if err != nil {
		return false
	}

	if strings.Contains(url.Host, facebookHost) == false {
		return false
	}

	ok, _ := regexp.Match(facebookVideoIframeSrcRegex, node.GetAtom(atom.Src))

	return ok
}

func NewFacebookVideoIframe(config map[string]string) (Validator, error) {
	return FacebookVideoIframe{}, nil
}
