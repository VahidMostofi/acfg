package endpointsagg

import (
	"github.com/pkg/errors"
	"github.com/vahidmostofi/acfg/internal/aggregators/restime"
)

type Endpoint struct {
	Name string
}

type EndpointsAggregator struct{
	responseTimeAggregator restime.ResponseTimeAggregator
	endpointFilters        map[string]map[string]interface{}
}

// NewEndpointsAggregator ...
// uses endpoint filters to select which endpoints' response times to gather. Current filters: HTTP_METHOD, URI_REGEX
// available kinds: influxdb
// for influxdb, it uses the url, token, organization and bucket pass them in args which is map[string]interface{}
func NewEndpointsAggregator(kind string, args map[string]interface{}, endpointFilters map[string]map[string]interface{})(*EndpointsAggregator, error){
	u := &EndpointsAggregator{}
	var err error
	if kind == "influxdb"{
		u.responseTimeAggregator, err = restime.NewInfluxDBRTA(
			args["url"].(string),
			args["token"].(string),
			args["organization"].(string),
			args["bucket"].(string),
		)

		if err != nil{
			return nil, errors.Wrap(err, "cant create InfluxDBRTA")
		}

		u.endpointFilters = endpointFilters

	}else{
		return nil, errors.Errorf("unknown kind: %s", kind)
	}

	return u, nil
}

func (e *EndpointsAggregator) GetListOfEndpointsBeingTracked() []*Endpoint {
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
