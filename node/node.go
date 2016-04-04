package node

import "github.com/jloup/scraper/aggregator"

type Node interface {
	InputNode(data interface{})
	ValidateNode() bool
	ExtractNode(aggregator.Aggregator) error
	Copy() Node
}

type ScraperNode struct {
	Node Node

	Aggregator aggregator.Aggregator

	Parent *ScraperNode
	Childs []*ScraperNode

	Depth   int
	hasData bool
}

func NewScraperNode(agg aggregator.Aggregator, n Node) ScraperNode {
	return ScraperNode{Node: n, Aggregator: agg, Parent: nil, Depth: -1, hasData: false}
}

func (s *ScraperNode) Init() {
	s.resetDepth()
	s.Aggregator.Reset()

	for _, child := range s.Childs {
		child.Init()
	}
}

func (s ScraperNode) Copy() *ScraperNode {
	var n Node = nil

	if s.Node != nil {
		n = s.Node.Copy()
	}
	scraper := ScraperNode{
		Node:       n,
		Aggregator: s.Aggregator.Duplicate(),
		Parent:     nil,
	}

	for _, child := range s.Childs {
		newChild := child.Copy()
		scraper.AddChild(newChild)
	}

	return &scraper
}

func (s *ScraperNode) AddChild(child *ScraperNode) {
	child.Parent = s
	s.Childs = append(s.Childs, child)
}

func (s *ScraperNode) AddChilds(childs ...*ScraperNode) {
	for _, child := range childs {
		child.Parent = s
	}
	s.Childs = append(s.Childs, childs...)
}

func (s *ScraperNode) resetDepth() {
	s.Depth = -1
	s.hasData = false
}

func (s *ScraperNode) join(agg aggregator.Aggregator) {
	s.hasData = true
	s.Aggregator.Join(agg)
}

func (s *ScraperNode) checkDepth(depth int, store *[]map[string]interface{}) {
	for _, child := range s.Childs {
		child.checkDepth(depth, store)
	}

	if depth <= s.Depth {
		if s.Parent != nil {

			if s.hasData == true || s.Childs == nil {
				s.Parent.join(s.Aggregator)
			}
		} else {
			s.Aggregator.Persist(store)
		}

		s.resetDepth()

	}
}

func (s *ScraperNode) End(store *[]map[string]interface{}) {
	s.checkDepth(-1, store)
}

func (s *ScraperNode) ProcessNode(n interface{}, store *[]map[string]interface{}, depth int) error {
	var err error

	s.checkDepth(depth, store)

	if s.Depth == -1 && s.Node != nil {

		s.Node.InputNode(n)

		if s.Node.ValidateNode() {
			if err = s.Node.ExtractNode(s.Aggregator); err != nil {
				return err
			}

			s.Depth = depth
		}

	} else {
		if s.Node == nil && s.Depth == -1 {
			s.Depth = depth
		}

		for _, child := range s.Childs {
			err = child.ProcessNode(n, store, depth)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

// wrap an array of scrapers in a top level scraper
func Wrap(scrapers []*ScraperNode) *ScraperNode {
	root := NewScraperNode(aggregator.NewAggregatorN(), nil)

	for _, scraper := range scrapers {
		root.AddChild(scraper)
	}

	return &root
}
