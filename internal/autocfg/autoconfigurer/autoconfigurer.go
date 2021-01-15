package autoconfigurer

import (
	"github.com/vahidmostofi/acfg/internal/autocfg"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type AutoConfigurationAgent interface{
	GetName() string
	GetInitialConfiguration(workload *workload.Workload, aggData *autocfg.AggregatedData) (map[string]*autocfg.Configuration, error)
	ConfigureNextStep(currentConfig map[string]*autocfg.Configuration, workload *workload.Workload, aggData *autocfg.AggregatedData) (map[string]*autocfg.Configuration, bool, error)
}

func CheckCondition(data *autocfg.AggregatedData, condition autocfg.Condition) (bool, error){
	if condition.Type == "ResponseTime"{
		value := condition.ComputeFn(*data.ResponseTimes[condition.EndpointName])
		if value <= condition.Threshold{
			return true, nil
		}
	}
	return false, nil
}

func int64Ptr(i int64) *int64 { return &i }