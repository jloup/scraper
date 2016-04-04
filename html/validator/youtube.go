package validator

import (
	"bytes"
	"regexp"

	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

const youtubeHost = "youtube.com"

var youtubeHostB = []byte(youtubeHost)
var youtubeRegex = `(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:v|embed)\/(?:[\w\-]{11})(?:\?[^'|"]*)?`

type YTIframeFast struct {
}

func (v YTIframeFast) Validate(node *nodedata.NodeData) bool {
	return bytes.Contains(node.GetAtom(atom.Src), youtubeHostB)
}

func NewYTIframeFast(config map[string]string) (Validator, error) {

	return YTIframeFast{}, nil
}

type YTIframe struct {
}

func (v YTIframe) Validate(node *nodedata.NodeData) bool {

	ok, _ := regexp.Match(youtubeRegex, node.GetAtom(atom.Src))
	return ok
}

func NewYTIframe(config map[string]string) (Validator, error) {

	return YTIframe{}, nil
}

type YTObjectFast struct {
	Key []byte
}

func (v YTObjectFast) Validate(node *nodedata.NodeData) bool {
	if node.Get(v.Key) == nil {
		return false
	}
	return bytes.Contains(node.Get(v.Key), youtubeHostB)
}

type YTObjectFastA struct {
	Key atom.Atom
}

func (v YTObjectFastA) Validate(node *nodedata.NodeData) bool {
	if bytes.Compare(node.GetAtom(v.Key), emptyString) == 0 {
		return false
	}
	return bytes.Contains(node.GetAtom(v.Key), youtubeHostB)
}

func NewYTObjectFast(config map[string]string) (Validator, error) {
	if a := atom.Lookup([]byte(config["key"])); a == 0 {
		return YTObjectFast{Key: []byte(config["key"])}, nil
	} else {
		return YTObjectFastA{Key: a}, nil
	}
}

type YTObject struct {
	Key []byte
}

func (v YTObject) Validate(node *nodedata.NodeData) bool {
	if node.Get(v.Key) == nil {
		return false
	}
	ok, _ := regexp.Match(youtubeRegex, node.Get(v.Key))
	return ok
}

type YTObjectA struct {
	Key atom.Atom
}

func (v YTObjectA) Validate(node *nodedata.NodeData) bool {
	if bytes.Compare(node.GetAtom(v.Key), emptyString) == 0 {
		return false
	}
	ok, _ := regexp.Match(youtubeRegex, node.GetAtom(v.Key))
	return ok
}

func NewYTObject(config map[string]string) (Validator, error) {
	if a := atom.Lookup([]byte(config["key"])); a == 0 {

		return YTObject{Key: []byte(config["key"])}, nil
	} else {
		return YTObjectA{Key: a}, nil
	}
}
