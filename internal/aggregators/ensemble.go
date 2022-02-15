package aggregators

import (
	"encoding/json"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	deploymentinfoagg "github.com/vahidmostofi/acfg/internal/aggregators/deploymentInfoAggregator"
	"github.com/vahidmostofi/acfg/internal/aggregators/endpointsagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/restime"
	"github.com/vahidmostofi/acfg/internal/aggregators/ussageagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/utilizations"
	"github.com/vahidmostofi/acfg/internal/aggregators/workloadagg"
	"github.com/vahidmostofi/acfg/internal/constants"
)

type Ensemble struct {
	Endpoints      *endpointsagg.EndpointsAggregator
	Workload       workloadagg.WorkloadAggregator
	Usage          *ussageagg.UsageAggregator
	DeploymentInfo deploymentinfoagg.DeploymentInfoAggregator
}

// DumpDataWithTimestamp dumps the history data from startTime to finishTime.
// the dumped data would have timestamp for each item stored.
// ResponseTimes, CPUUtilizations, DeploymentInfos
func (e *Ensemble) DumpDataWithTimestamp(startTime, finishTime int64) ([]byte, error) {
	var err error
	type Dump struct {
		StartTime        int64
		FinishTime       int64
		ResponseTimes    map[string][]restime.TimestampedResponseTime
		UsageUtilization map[string][]utilizations.CPUTimestampedUsage
		DeploymentInfo   map[string][]deploymentinfoagg.TimestampedDeploymentInfo
	}

	d := Dump{StartTime: startTime, FinishTime: finishTime}
	d.ResponseTimes, err = e.Endpoints.GetEndpointsResponseTimesWithTimestamp(startTime, finishTime)
	if err != nil {
		panic(err)
	}

	if viper.GetBool(constants.DumpWithCPUInfo) {
		d.UsageUtilization, err = e.Usage.GetAggregatedCPUUtilizationsWithTimestamped(startTime, finishTime)
		if err != nil {
			panic(err)
		}
	}

	d.DeploymentInfo, err = e.DeploymentInfo.GetAllDeploymentsInfoTimestamped(startTime, finishTime, make(map[string]interface{}))
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return b, nil
}

func (e *Ensemble) AggregateData(startTime, finishTime int64) (*AggregatedData, error) {
	log.Debugf("aggregating data")
	var err error

	// response times
	ag := &AggregatedData{StartTime: &startTime, FinishTime: &finishTime}

	if e.Endpoints != nil {
		// Response Times
		ag.ResponseTimes, err = e.Endpoints.GetEndpointsResponseTimes(startTime, finishTime)
		if err != nil {
			return nil, errors.Wrap(err, "error while getting response times")
		}
	}

	if e.Workload != nil {
		// Workload
		ag.HappenedWorkload, err = e.Workload.GetWorkload(startTime, finishTime)
		if err != nil {
			return nil, errors.Wrap(err, "error while getting the workload that happened")
		}
	}

	// TODO handle system structure
	// System Structure
	// ag.SystemStructure = a.systemStructure

	if e.Usage != nil {
		// Resource Utilization
		// CPU
		ag.CPUUtilizations, err = e.Usage.GetAggregatedCPUUtilizations(startTime, finishTime)
		if err != nil {
			return nil, errors.Wrap(err, "error while getting the CPU utilizations")
		}

		ag.MemUtilizations, err = e.Usage.GetAggregatedMemUtilizations(startTime, finishTime)
		if err != nil {
			return nil, errors.Wrap(err, "error while getting the mem utilizations")
		}
	}

	// Memory //TODO
	if e.DeploymentInfo != nil {
		ag.DeploymentInfos, err = e.DeploymentInfo.GetAllDeploymentsInfo(startTime, finishTime, make(map[string]interface{}))
		if err != nil {
			return nil, errors.Wrap(err, "error while getting deployment infos")
		}
	}

	return ag, nil
}
