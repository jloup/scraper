package aggregator

// AggregatorArray
//
// AGGREGATE {station: Concorde} => [{station: Concorde, ArrayKey: []}]
//
// JOIN {metro: 1}     => [{station: Concorde, ArrayKey: [1]}]
//
// JOIN {metro: 4}     => [{station: Concorde, ArrayKey: [1, 4]}]
//
// JOIN {metro: 8}     => [{station: Concorde, ArrayKey: [1, 4, 8]}]

type AggregatorArray struct {
	ArrayKey string
	Store    map[string]interface{}
}

func NewAggregatorArray(key string) *AggregatorArray {
	a := AggregatorArray{ArrayKey: key}
	a.Reset()
	return &a
}

func (a *AggregatorArray) Aggregate(key string, value interface{}) {
	a.Store[key] = value
}

func (a *AggregatorArray) GetAggregate() []map[string]interface{} {

	length := len(a.Store)
	if length > 0 {
		s := make([]map[string]interface{}, 1, 1)
		s[0] = a.Store
		return s
	}

	return nil

}

func (a *AggregatorArray) Join(agg Aggregator) {
	s := agg.GetAggregate()
	for _, el := range s {
		for _, value := range el {
			if _, ok := a.Store[a.ArrayKey]; !ok {
				a.Store[a.ArrayKey] = make([]interface{}, 0)
			}
			arr := a.Store[a.ArrayKey].([]interface{})
			arr = append(arr, value)
			a.Store[a.ArrayKey] = arr
		}
	}
	agg.Reset()
}

func (a *AggregatorArray) Persist(store *[]map[string]interface{}) {
	if len(a.Store) > 0 {
		*store = append(*store, a.Store)
		a.Reset()
	}
}

func (a *AggregatorArray) Reset() {
	a.Store = nil
	a.Store = make(map[string]interface{})
}

func (a *AggregatorArray) Duplicate() Aggregator {
	dup := NewAggregatorArray(a.ArrayKey)
	return dup
}
