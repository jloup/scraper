package aggregator

// AggregatorN is a basic aggregator which combines its own data with each Join
//
// AGGREGATE {container: main} => []
//
// JOIN {href: golang.org} => [{container: main, href: golang.org}]
//
// JOIN {href: google.fr} => [{container: main, href: golang.org}, {container: main, href: google.fr}]
//
type AggregatorN struct {
	Store []map[string]interface{}
	Data  map[string]interface{}
}

func NewAggregatorN() *AggregatorN {
	a := AggregatorN{Data: make(map[string]interface{})}
	a.Reset()
	return &a
}

func (a *AggregatorN) GetAggregate() []map[string]interface{} {

	if len(a.Store) == 0 {

		if len(a.Data) > 0 {
			s := make([]map[string]interface{}, 1, 1)
			s[0] = a.Data
			return s
		}
		return nil

	} else {

		if len(a.Data) > 0 {

			s := make([]map[string]interface{}, len(a.Store), len(a.Store))
			for i, _ := range s {
				s[i] = a.Store[i]
				for key, value := range a.Data {
					s[i][key] = value
				}
			}
			return s
		}

		return a.Store
	}
}

func (a *AggregatorN) Aggregate(key string, value interface{}) {
	a.Data[key] = value
}

func (a *AggregatorN) Join(agg Aggregator) {
	s := agg.GetAggregate()

	a.Store = append(a.Store, s...)
	agg.Reset()
}

func (a *AggregatorN) Persist(store *[]map[string]interface{}) {
	agg := a.GetAggregate()
	if len(agg) > 0 {
		*store = append(*store, agg...)
		a.Reset()
	}

}

func (a *AggregatorN) Reset() {
	a.Store = nil
	a.Data = nil
	a.Store = make([]map[string]interface{}, 0, 0)
	a.Data = make(map[string]interface{})
}

func (a *AggregatorN) Duplicate() Aggregator {
	dup := NewAggregatorN()
	return dup
}
