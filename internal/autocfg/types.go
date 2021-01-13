package autocfg

import (
	"github.com/vahidmostofi/acfg/internal/aggregators/restime"
	"github.com/vahidmostofi/acfg/internal/aggregators/sysstructureagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/utilizations"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type Configuration struct{
	ResourceType string
	ReplicaCount *int64
	CPU *int64
	Memory *int64
	EnvironmentValues map[string]string
}

func (c *Configuration) GetHash(version string) []byte{
	panic("not implemneted yet")
	return nil
}

type AggregatedData struct{
	ResponseTimes map[string]*restime.ResponseTimes
	CPUUtilizations map[string]*utilizations.CPUUtilizations
	SystemStructure *sysstructureagg.SystemStructure
	HappenedWorkload *workload.Workload
}

type IterationInformation struct{
	Configuration map[string]*Configuration
	StartTime int64
	FinishTime int64
	AggregatedData *AggregatedData
}

type TestInformation struct{
	VersionCode string
	AutoconfiguringApproach string
	Iterations []*IterationInformation
	InputWorkload *workload.Workload
}