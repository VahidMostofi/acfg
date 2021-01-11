package endpointsagg

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/restime"
)

type Endpoint struct {
	Name string
}

type EndpointsAggregator struct{
	responseTimeAggregator restime.ResponseTimeAggregator
	endpointFilters map[string]map[string]interface{}
}

// NewEndpointsAggregator ...
// uses endpoint filters to select which endpoints' response times to gather. Current filters: HTTP_METHOD, URI_REGEX
// available kinds: influxdb
// for influxdb, it uses the url, token, org and bucket provided in config
func NewEndpointsAggregator(kind string)(*EndpointsAggregator, error){
	u := &EndpointsAggregator{}
	var err error
	if kind == "influxdb"{
		u.responseTimeAggregator, err = restime.NewInfluxDBRTA(
			viper.GetString(constants.CONFIG_INFLUXDB_URL),
			viper.GetString(constants.CONFIG_INFLUXDB_TOKEN),
			viper.GetString(constants.CONFIG_INFLUXDB_ORG),
			viper.GetString(constants.CONFIG_INFLUXDB_BUCKET))

		if err != nil{
			return nil, errors.Wrap(err, "cant create InfluxDBRTA")
		}

		temp := viper.Get(constants.CONFIG_ENDPOINTS_FILTERS)
		tempConverted, ok := temp.(map[string]map[string]interface{})
		if !ok {
			return nil, errors.Errorf("cant find endpoints filters in configs using: %s with type map[string]map[string]interface{}", constants.CONFIG_ENDPOINTS_FILTERS)
		}
		u.endpointFilters = tempConverted

	}else{
		return nil, errors.Errorf("unknown kind: %s", kind)
	}

	return u, nil
}

func (e *EndpointsAggregator) GetListOfEndpointsBeingTracked() []*Endpoint{
	endpoints := make([]*Endpoint,0)
	for resourceName := range e.endpointFilters{
		endpoints = append(endpoints, &Endpoint{resourceName})
	}

	return endpoints
}

func (e *EndpointsAggregator) GetEndpointsResponseTimes(startTime, finishTime int64) (map[string]*restime.ResponseTimes, error) {
	result := make(map[string]*restime.ResponseTimes)

	for _, r := range e.GetListOfEndpointsBeingTracked(){
		_responseTimes, err := e.responseTimeAggregator.GetResponseTimes(startTime, finishTime, e.endpointFilters[r.Name])
		if err != nil {
			return nil, errors.Wrap(err, "error while getting response times for " + r.Name)
		}
		result[r.Name] = _responseTimes
	}

	return result, nil
}
