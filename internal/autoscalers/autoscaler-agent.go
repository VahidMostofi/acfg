package autoscalers

import (
	"github.com/vahidmostofi/acfg/internal/aggregators"
)


type Agent interface{
	GetName() string
	Evaluate(aggData *aggregators.AggregatedData) (map[string]int, error)
}