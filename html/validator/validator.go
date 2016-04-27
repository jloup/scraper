//Package validator implements HTML tag validation functions
package validator

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/jloup/scraper/html/nodedata"
	"golang.org/x/net/html/atom"
)

var emptyString = []byte("")

type ValidatorInitError struct {
	What string
}

func (e ValidatorInitError) Error() string {
	return fmt.Sprintf("Validator initialization error: %s", e.What)
}

// validator tests if a node is valid through a validate function
type Validator interface {
	Validate(node *nodedata.NodeData) bool
}

// check that a node has an attribute
type Exists struct {
	Attr []byte
}

func (v Exists) Validate(node *nodedata.NodeData) bool {
	return node.Get(v.Attr) != nil
}

// check that a node has an attribute - atom mode
type ExistsA struct {
	Attr atom.Atom
}

func (v ExistsA) Validate(node *nodedata.NodeData) bool {
	return node.GetAtom(v.Attr) != nil
}

func NewExists(config map[string]string) (Validator, error) {
	if config["attr"] == "" {
		return nil, ValidatorInitError{What: "Missing attr key in config"}
	}

	if a := atom.Lookup([]byte(config["attr"])); a == 0 {
		return Exists{Attr: []byte(config["attr"])}, nil
	} else {
		return ExistsA{Attr: a}, nil
	}
}

// check that node's attribute is equal to a specific value. Assume that attribute exists
type AttrEquals struct {
	Attr  []byte
	Value []byte
}

func (v AttrEquals) Validate(node *nodedata.NodeData) bool {
	return bytes.Compare(node.Get(v.Attr), v.Value) == 0
}

type AttrEqualsA struct {
	Attr  atom.Atom
	Value []byte
}

func (v AttrEqualsA) Validate(node *nodedata.NodeData) bool {
	return bytes.Compare(node.GetAtom(v.Attr), v.Value) == 0
}

func NewAttrEquals(config map[string]string) (Validator, error) {
	if config["attr"] == "" || config["value"] == "" {
		return nil, ValidatorInitError{What: "Missing 'attr' or 'value' key in config"}
	}
	if a := atom.Lookup([]byte(config["attr"])); a == 0 {
		return AttrEquals{Attr: []byte(config["attr"]), Value: []byte(config["value"])}, nil
	} else {
		return AttrEqualsA{Attr: a, Value: []byte(config["value"])}, nil
	}
}

// check that node's attribute strings.Contains a specific value. Assume that attribute exists
type AttrContains struct {
	Attr  []byte
	Value []byte
}

func (v AttrContains) Validate(node *nodedata.NodeData) bool {
	return bytes.Contains(node.Get(v.Attr), v.Value)
}

type AttrContainsA struct {
	Attr  atom.Atom
	Value []byte
}

func (v AttrContainsA) Validate(node *nodedata.NodeData) bool {
	return bytes.Contains(node.GetAtom(v.Attr), v.Value)
}

func NewAttrContains(config map[string]string) (Validator, error) {
	if config["attr"] == "" || config["value"] == "" {
		return nil, ValidatorInitError{What: "Missing 'attr' or 'value' key in config"}
	}

	if a := atom.Lookup([]byte(config["attr"])); a == 0 {
		return AttrContains{Attr: []byte(config["attr"]), Value: []byte(config["value"])}, nil
	} else {
		return AttrContainsA{Attr: a, Value: []byte(config["value"])}, nil
	}
}

// validate a specified attribute from a HTML tag with a Regexp expression
type Regexp struct {
	R        string
	Attr     []byte
	AtomAttr atom.Atom
}

func (r Regexp) Validate(node *nodedata.NodeData) bool {

	var haystack []byte
	if len(r.Attr) > 0 {
		haystack = node.Get(r.Attr)
	} else {
		haystack = node.GetAtom(r.AtomAttr)
	}

	ok, err := regexp.Match(r.R, haystack)
	if err != nil {
		panic(err)
	}

	return ok
}

func NewRegexp(config map[string]string) (Validator, error) {
	if config["attr"] == "" {
		return nil, ValidatorInitError{What: "Missing attr key in config"}
	}

	if config["regexp"] == "" {
		return nil, ValidatorInitError{What: "Missing attr regexp in config"}
	}

	if a := atom.Lookup([]byte(config["attr"])); a == 0 {
		return Regexp{Attr: []byte(config["attr"]), R: config["regexp"]}, nil
	} else {
		return Regexp{AtomAttr: a, R: config["regexp"]}, nil
	}
}
