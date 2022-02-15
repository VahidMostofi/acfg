package aggregators

import (
	deploymentinfoagg "github.com/vahidmostofi/acfg/internal/aggregators/deploymentInfoAggregator"
	"github.com/vahidmostofi/acfg/internal/aggregators/restime"
	"github.com/vahidmostofi/acfg/internal/aggregators/sysstructureagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/utilizations"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type AggregatedData struct {
	ResponseTimes      map[string]*restime.ResponseTimes           `yaml:"responseTimes"`
	CPUUtilizations    map[string]*utilizations.CPUUtilizations    `yaml:"CPUUtilizations"`
	CPUPSIUtilizations map[string]*utilizations.CPUPsiUtilizations `yaml:"CPUPsiUtilizations"`
	MemUtilizations    map[string]*utilizations.MemUtilizations    `yaml:"MemUtilizations"`
	MemPsiUtilizations map[string]*utilizations.MemPsiUtilizations `yaml:"MemPsiUtilizations"`
	SystemStructure    *sysstructureagg.SystemStructure            `yaml:"structure"`
	HappenedWorkload   *workload.Workload                          `yaml:"workload"`
	StartTime          *int64                                      `yaml:"startTime"`
	FinishTime         *int64                                      `yaml:"finishTime"`
	DeploymentInfos    map[string]*deploymentinfoagg.DeploymentInfo
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

type ResourceType string

const (
	CPU ResourceType = "cpu"
	Mem              = "mem"
)

type MinMaxResource struct {
	ResourceType ResourceType
	ResourceName string
	Value        float64
}

func (ag *AggregatedData) GetMinMaxResourcesBasedOnPsi(endpoint string) (*MinMaxResource, *MinMaxResource) {
	minResource := &MinMaxResource{
		CPU,
		"",
		-1,
	}
	maxResource := &MinMaxResource{
		CPU,
		"",
		float64(^int64(0) >> 1),
	}
	for _, resourceName := range ag.SystemStructure.GetEndpoints2Resources()[endpoint] {
		cpuPSI, err := ag.CPUPSIUtilizations[resourceName].GetMean()
		cpuResource := MinMaxResource{
			CPU,
			resourceName,
			cpuPSI,
		}
		memPSI, err := ag.MemPsiUtilizations[resourceName].GetMean()
		memResource := MinMaxResource{
			Mem,
			resourceName,
			memPSI,
		}
		if err != nil {
			panic(err)
		}
		if cpuPSI > memPSI {
			//	CPU gt than mem
			if cpuPSI > maxResource.Value {
				//	new max value
				maxResource = &cpuResource
			}
			if memPSI < minResource.Value {
				//	new min value
				minResource = &memResource
			}
		} else {
			//	Mem gt than cpu
			if memPSI > maxResource.Value {
				//	new max value
				maxResource = &memResource
			}
			if cpuPSI < minResource.Value {
				//	new min value
				minResource = &cpuResource
			}
		}
	}

	return minResource, maxResource
}
