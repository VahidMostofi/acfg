package ussageagg

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/aggregators/utilizations"
	"github.com/vahidmostofi/acfg/internal/constants"
)

type Resource struct {
	Name string
}

type UsageAggregator struct {
	cpuUsageAggregator utilizations.CPUUtilizationAggregator
	resourceFilters    map[string]map[string]interface{}
}

// NewUsageAggregator ...
// uses resource filters to select which pods CPU usage to consider and compute. Current filters: POD_NAME_REGEX
// available kinds: influxdb
// for influxdb, it uses the url, token, org and bucket provided in config
func NewUsageAggregator(kind string, args map[string]interface{}, resourceFilters map[string]map[string]interface{}) (*UsageAggregator, error) {
	u := &UsageAggregator{}
	var err error
	if kind == "influxdb" {
		u.cpuUsageAggregator, err = utilizations.NewInfluxDBCPUUA(
			viper.GetString(constants.EndpointsAggregatorArgsURL),
			viper.GetString(constants.EndpointsAggregatorArgsToken),
			viper.GetString(constants.EndpointsAggregatorArgsOrganization),
			viper.GetString(constants.EndpointsAggregatorArgsBucket))

		if err != nil {
			return nil, errors.Wrap(err, "cant create InfluxDBCPUUA")
		}

		u.resourceFilters = resourceFilters

	} else {
		return nil, errors.Errorf("unknown kind: %s", kind)
	}

	return u, nil
}

func (u *UsageAggregator) GetListOfResourcesBeingTracked() []*Resource {
	resources := make([]*Resource, 0)
	for resourceName := range u.resourceFilters {
		resources = append(resources, &Resource{resourceName})
	}

	return resources
}

func (u *UsageAggregator) GetAggregatedCPUUtilizations(startTime, finishTime int64) (map[string]*utilizations.CPUUtilizations, error) {
	result := make(map[string]*utilizations.CPUUtilizations)

	for _, r := range u.GetListOfResourcesBeingTracked() {
		_cpuUtil, err := u.cpuUsageAggregator.GetCPUUtilizations(startTime, finishTime, u.resourceFilters[r.Name])
		if err != nil {
			return nil, errors.Wrap(err, "error while getting CPU utilization for "+r.Name)
		}
		result[r.Name] = _cpuUtil
	}

	return result, nil
}

func (u *UsageAggregator) GetAggregatedCPUUtilizationsWithTimestamped(startTime, finishTime int64) (map[string][]utilizations.TimestampedUsage, error) {
	result := make(map[string][]utilizations.TimestampedUsage)

	for _, r := range u.GetListOfResourcesBeingTracked() {
		_cpuUtil, err := u.cpuUsageAggregator.GetCPUUtilizationsWithTimestamp(startTime, finishTime, u.resourceFilters[r.Name])
		if err != nil {
			return nil, errors.Wrap(err, "error while getting CPU utilization for "+r.Name)
		}
		result[r.Name] = _cpuUtil
	}

	return result, nil
}
