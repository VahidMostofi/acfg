package aggregators

import (
	deploymentinfoagg "github.com/vahidmostofi/acfg/internal/aggregators/deploymentInfoAggregator"
	"github.com/vahidmostofi/acfg/internal/aggregators/restime"
	"github.com/vahidmostofi/acfg/internal/aggregators/sysstructureagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/utilizations"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type AggregatedData struct {
	ResponseTimes    map[string]*restime.ResponseTimes        `yaml:"responseTimes"`
	CPUUtilizations  map[string]*utilizations.CPUUtilizations `yaml:"CPUUtilizations"`
	SystemStructure  *sysstructureagg.SystemStructure         `yaml:"structure"`
	HappenedWorkload *workload.Workload                       `yaml:"workload"`
	StartTime        *int64                                   `yaml:"startTime"`
	FinishTime       *int64                                   `yaml:"finishTime"`
	DeploymentInfos  map[string]*deploymentinfoagg.DeploymentInfo
}

func (ag *AggregatedData) GetMinMaxResourcesBasedOnCPUUtil(endpoint string) (string, string) {
	var maxValue float64 = 0
	var minValue float64 = 1000000
	var minName = ""
	var maxName = ""

	for _, resourceName := range ag.SystemStructure.GetEndpoints2Resources()[endpoint] {
		cpuu := ag.CPUUtilizations[resourceName]
		m, err := cpuu.GetMean()
		if err != nil {
			panic(err)
		}
		if minValue > m {
			minValue = m
			minName = resourceName
		}
		if maxValue < m {
			maxValue = m
			maxName = resourceName
		}
	}
	return minName, maxName
}
