package autoscalers

import (
	"fmt"
	"github.com/vahidmostofi/acfg/internal/aggregators"
)

func NewHybridAutoscaler(endpoints, resources []string) (Agent, error){
	h := &Hybrid{
		endpoints: endpoints,
		resources: resources,
	}

	return h, nil
}

type Hybrid struct{
	endpoints []string
	resources []string
}

func (h *Hybrid) GetName() string{
	return "hybrid"
}

func (h *Hybrid) Evaluate(aggData *aggregators.AggregatedData) (map[string]int, error){
	for name, cpuu := range aggData.CPUUtilizations{
		m, e := cpuu.GetMean()
		if e != nil{
			panic(e)
		}
		fmt.Println(name, m)
	}
	return nil,nil
}


