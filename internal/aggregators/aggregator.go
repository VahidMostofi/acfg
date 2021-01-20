package aggregators

import (
	"github.com/vahidmostofi/acfg/internal/aggregators/restime"
	"github.com/vahidmostofi/acfg/internal/aggregators/sysstructureagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/utilizations"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type AggregatedData struct{
	ResponseTimes map[string]*restime.ResponseTimes				`yaml:"responseTimes"`
	CPUUtilizations map[string]*utilizations.CPUUtilizations	`yaml:"CPUUtilizations"`
	SystemStructure *sysstructureagg.SystemStructure			`yaml:"structure"`
	HappenedWorkload *workload.Workload							`yaml:"workload"`
}
