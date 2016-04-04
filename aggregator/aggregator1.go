package aggregator

// Aggregator1 arrange all data in a unique hash
//
// AGGREGATE {name: Jean} => [{name: Jean}]
//
// JOIN {age: 14}     => [{name: Jean, age: 14}]
//
// JOIN {height: 178} => [{name: Jean, age: 14, height: 178}]
//
// JOIN {age: 15}     => [{name: Jean, age: 15, height: 178}]
type Aggregator1 struct {
	Store map[string]interface{}
}

func NewAggregator1() *Aggregator1 {
	a := Aggregator1{}
	a.Reset()
	return &a
}

func (a *Aggregator1) Aggregate(key string, value interface{}) {
	a.Store[key] = value
}

func (a *Aggregator1) GetAggregate() []map[string]interface{} {

	if len(a.Store) > 0 {
		s := make([]map[string]interface{}, 1, 1)
		s[0] = a.Store
		return s
	}

	return nil

}

func (a *Aggregator1) Join(agg Aggregator) {
	s := agg.GetAggregate()
	for _, el := range s {
		for key, value := range el {
			a.Aggregate(key, value)
		}
	}
	agg.Reset()
}

func (a *Aggregator1) Persist(store *[]map[string]interface{}) {
	if len(a.Store) > 0 {
		*store = append(*store, a.Store)
		a.Reset()
	}
}

func (a *Aggregator1) Reset() {
	a.Store = nil
	a.Store = make(map[string]interface{})
}

func (a *Aggregator1) Duplicate() Aggregator {
	dup := NewAggregator1()
	return dup
}
