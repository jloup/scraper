package aggregator

import "strconv"

// AggregatorList index data that share the same key
//
// AGGREGATE {station: Concorde} => [{station: Concorde}]
//
// JOIN {metro: 1}     => [{station: Concorde, metro: 1}]
//
// JOIN {metro: 4}     => [{station: Concorde, metro0: 1, metro1: 4}]
//
// JOIN {metro: 8}     => [{station: Concorde, metro0: 1, metro1: 4, metro2: 8}]
type AggregatorList struct {
	Store   map[string]interface{}
	Indexes map[string]int
}

func NewAggregatorList() *AggregatorList {
	a := AggregatorList{}
	a.Reset()
	return &a
}

func (a *AggregatorList) Aggregate(key string, value interface{}) {
	if _, exists := a.Indexes[key]; exists {

		if a.Indexes[key] == 0 {
			a.Store[key+"0"] = a.Store[key]
			delete(a.Store, key)
		}

		a.Indexes[key] += 1
		key = key + strconv.Itoa(a.Indexes[key])

	} else {
		a.Indexes[key] = 0
	}

	a.Store[key] = value
}

func (a *AggregatorList) GetAggregate() []map[string]interface{} {

	length := len(a.Store)
	if length > 0 {
		s := make([]map[string]interface{}, 1, 1)
		s[0] = a.Store
		return s
	}

	return nil

}

func (a *AggregatorList) Join(agg Aggregator) {
	s := agg.GetAggregate()
	for _, el := range s {
		for key, value := range el {
			a.Aggregate(key, value)
		}
	}
	agg.Reset()
}

func (a *AggregatorList) Persist(store *[]map[string]interface{}) {
	if len(a.Store) > 0 {
		*store = append(*store, a.Store)
		a.Reset()
	}
}

func (a *AggregatorList) Reset() {
	a.Store = nil
	a.Store = make(map[string]interface{})
	a.Indexes = make(map[string]int)
}

func (a *AggregatorList) Duplicate() Aggregator {
	dup := NewAggregatorList()
	return dup
}
