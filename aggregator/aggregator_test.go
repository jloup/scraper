package aggregator

import (
	"testing"
)

func TestAggregatorBase(t *testing.T) {
	parent := NewAggregatorN()
	child := NewAggregatorN()
	child2 := NewAggregatorN()

	parent.Aggregate("id", "main")

	child.Aggregate("a", "first link")
	child.Aggregate("class", "link")

	child2.Aggregate("a", "second link")
	child2.Aggregate("class", "link")

	parent.Join(child)
	parent.Join(child2)

	s := make([]map[string]interface{}, 0, 0)

	parent.Persist(&s)

	t.Log(s)
}

func TestAggregator1Base(t *testing.T) {
	root := NewAggregatorN()
	parent := NewAggregator1()
	child := NewAggregatorN()
	child2 := NewAggregatorN()

	root.Aggregate("page", "root")

	parent.Aggregate("id", "561457545")

	child.Aggregate("prix", "920euros")

	child2.Aggregate("date", "23 juillet")

	parent.Join(child)
	parent.Join(child2)

	root.Join(parent)
	root.Join(parent)

	s := make([]map[string]interface{}, 0, 0)

	root.Persist(&s)

	t.Log(s)
}
