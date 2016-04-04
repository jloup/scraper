package validator

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

const soundcloudHost = "soundcloud.com"

var soundcloudHostB = []byte(soundcloudHost)

type SCIframeFast struct {
}

func (v SCIframeFast) Validate(node *nodedata.NodeData) bool {
	return bytes.Contains(node.GetAtom(atom.Src), soundcloudHostB)
}

func NewSCIframeFast(config map[string]string) (Validator, error) {
	return SCIframeFast{}, nil
}

type SCIframe struct {
}

func (v SCIframe) Validate(node *nodedata.NodeData) bool {
	url, err := url.Parse(string(node.GetAtom(atom.Src)))
	if err != nil {
		return false
	}

	return strings.Contains(url.Host, soundcloudHost)
}

func NewSCIframe(config map[string]string) (Validator, error) {
	return SCIframe{}, nil
}

type SCObjectFast struct {
}

func (v SCObjectFast) Validate(node *nodedata.NodeData) bool {
	return bytes.Contains(node.GetAtom(atom.Value), soundcloudHostB)
}

func NewSCObjectFast(config map[string]string) (Validator, error) {
	return SCObjectFast{}, nil
}

type SCObject struct {
}

func (v SCObject) Validate(node *nodedata.NodeData) bool {
	url, err := url.Parse(string(node.GetAtom(atom.Value)))
	if err != nil {
		return false
	}

	return strings.Contains(url.Host, soundcloudHost)
}

func NewSCObject(config map[string]string) (Validator, error) {

	return SCObject{}, nil
}
