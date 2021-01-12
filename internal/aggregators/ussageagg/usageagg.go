package ussageagg

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/aggregators/utilizations"
	"github.com/vahidmostofi/acfg/internal/constants"
)

type Resource struct{
	Name string
}

type UsageAggregator struct{
	cpuUsageAggregator utilizations.CPUUtilizationAggregator
	resourceFilters    map[string]map[string]interface{}
}

// NewUsageAggregator ...
// uses resource filters to select which pods CPU usage to consider and compute. Current filters: POD_NAME_REGEX
// available kinds: influxdb
// for influxdb, it uses the url, token, org and bucket provided in config
func NewUsageAggregator(kind string)(*UsageAggregator, error){
	u := &UsageAggregator{}
	var err error
	if kind == "influxdb"{
		u.cpuUsageAggregator,err = utilizations.NewInfluxDBCPUUA(
			viper.GetString(constants.CONFIG_INFLUXDB_URL),
			viper.GetString(constants.CONFIG_INFLUXDB_TOKEN),
			viper.GetString(constants.CONFIG_INFLUXDB_ORG),
			viper.GetString(constants.CONFIG_INFLUXDB_BUCKET))

		if err != nil{
			return nil, errors.Wrap(err, "cant create InfluxDBCPUUA")
		}

		temp := viper.Get(constants.CONFIG_RESOURCE_FILTERS)
		tempConverted, ok := temp.(map[string]map[string]interface{})
		if !ok {
			return nil, errors.Errorf("cant find resource filters in configs using: %s with type map[string]map[string]interface{}", constants.CONFIG_RESOURCE_FILTERS)
		}
		u.resourceFilters = tempConverted

	}else{
		return nil, errors.Errorf("unknown kind: %s", kind)
	}

	return u, nil
}

func (u *UsageAggregator) GetListOfResourcesBeingTracked() []*Resource {
	resources := make([]*Resource,0)
	for resourceName := range u.resourceFilters{
		resources = append(resources, &Resource{resourceName})
	}

	return resources
}

func (u *UsageAggregator) GetAggregatedCPUUtilizations(startTime, finishTime int64) (map[string]*utilizations.CPUUtilizations, error) {
	result := make(map[string]*utilizations.CPUUtilizations)

	for _, r := range u.GetListOfResourcesBeingTracked(){
		_cpuUtil, err := u.cpuUsageAggregator.GetCPUUtilizations(startTime, finishTime, u.resourceFilters[r.Name])
		if err != nil {
			return nil, errors.Wrap(err, "error while getting CPU utilization for " + r.Name)
		}
		result[r.Name] = _cpuUtil
	}

	return result, nil
}

