package validator

import (
	"bytes"
	"regexp"

	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

const dailymotionHost = "dailymotion.com"

var dailymotionHostB = []byte(dailymotionHost)
var dailymotionRegex = `(?:https?:\/\/)?(?:www\.)?dailymotion\.com\/embed\/video\/`

type DailymotionIframeFast struct {
}

func (d DailymotionIframeFast) Validate(node *nodedata.NodeData) bool {
	return bytes.Contains(node.GetAtom(atom.Src), dailymotionHostB)
}

func NewDailymotionIframeFast(config map[string]string) (Validator, error) {
	return DailymotionIframeFast{}, nil
}

type DailymotionIframe struct {
}

func (d DailymotionIframe) Validate(node *nodedata.NodeData) bool {
	ok, _ := regexp.Match(dailymotionRegex, node.GetAtom(atom.Src))
	return ok
}

func NewDailymotionIframe(config map[string]string) (Validator, error) {
	return DailymotionIframe{}, nil
}
