package validator

import (
	"bytes"
	"regexp"

	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

const vimeoHost = "vimeo.com"

var vimeoHostB = []byte(vimeoHost)
var vimeoRegex = `(?:https?:\/\/)?(?:www\.)?(?:player\.)?vimeo\.com\/video\/`

type VimeoIframeFast struct {
}

func (v VimeoIframeFast) Validate(node *nodedata.NodeData) bool {
	return bytes.Contains(node.GetAtom(atom.Src), vimeoHostB)
}

func NewVimeoIframeFast(config map[string]string) (Validator, error) {
	return VimeoIframeFast{}, nil
}

type VimeoIframe struct {
}

func (d VimeoIframe) Validate(node *nodedata.NodeData) bool {
	ok, _ := regexp.Match(vimeoRegex, node.GetAtom(atom.Src))
	return ok
}

func NewVimeoIframe(config map[string]string) (Validator, error) {
	return VimeoIframe{}, nil
}
