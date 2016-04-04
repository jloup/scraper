package nodedata

import (
	"bytes"

	html "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// HTML tag attributes
type Attribute struct {
	N []byte
	C []byte
}

type AtomAttribute struct {
	N atom.Atom
	C []byte
}

// contains all information about a HTML tag
type NodeData struct {
	Type        html.NodeType
	TagString   []byte
	TagAtom     atom.Atom
	Attr        []Attribute
	AttrAtom    []AtomAttribute
	TextContent []byte
}

func (n *NodeData) Get(name []byte) []byte {
	for i := 0; i < len(n.Attr); i++ {
		if n.Attr[i].N != nil && bytes.Compare(name, n.Attr[i].N) == 0 {
			return n.Attr[i].C
		}
	}
	return nil
}

func (n *NodeData) Set(name []byte, v []byte) {
	for i := 0; i < len(n.Attr); i++ {
		if nil == n.Attr[i].N {
			n.Attr[i].N = name
			n.Attr[i].C = v
			return
		}
	}
	n.Attr = append(n.Attr, Attribute{name, v})
}

func (n *NodeData) GetAtom(name atom.Atom) []byte {
	for i := 0; i < len(n.AttrAtom); i++ {
		if name == n.AttrAtom[i].N {
			return n.AttrAtom[i].C
		}
	}
	return nil
}

func (n *NodeData) SetAtom(name atom.Atom, v []byte) {
	for i := 0; i < len(n.AttrAtom); i++ {
		if 0x0 == n.AttrAtom[i].N {
			n.AttrAtom[i].N = name
			n.AttrAtom[i].C = v
			return
		}
	}
	n.AttrAtom = append(n.AttrAtom, AtomAttribute{name, v})
}
