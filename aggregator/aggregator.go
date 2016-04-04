//Package aggregator provides various strategies to organize scraped data
package aggregator

import "strings"

// Aggregator is an interface used by Scraper
//
//- Aggregate method should store in a key/value hash the data scrapped from one Scraper
//
//- Join method combine data stored in two Aggregators
//
//- GetAggregate method return data gathered with Aggregate and Join methods
//
//- Persist method persist GetAggregate result to an array
//
//- Reset method delete all data gathered with Aggregate and Join methods
//
//- Duplicate method return a clear Aggregator of the same type that target Aggregator
type Aggregator interface {
	Aggregate(key string, value interface{})
	Join(agg Aggregator)
	GetAggregate() []map[string]interface{}
	Persist(store *[]map[string]interface{})
	Reset()
	Duplicate() Aggregator
}

func NewAggregatorFromConfig(str string) Aggregator {

	if len(str) >= 3 && strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]") {
		str = strings.TrimPrefix(str, "[")
		str = strings.TrimSuffix(str, "]")
		return NewAggregatorArray(str)
	}

	switch str {
	case "1":
		return NewAggregator1()
	case "list":
		return NewAggregatorList()
	case "N":
		return NewAggregatorN()
	}

	return NewAggregatorN()
}
