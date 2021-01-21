package workloadagg

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/dataaccess"
	"github.com/vahidmostofi/acfg/internal/workload"
	"math"
	"strconv"
	"strings"
)

// WorkloadAggregator somehow returns the workload which was running between specific times
type WorkloadAggregator interface{
	GetWorkload(startTime, finishTime int64, endpointFilters map[string]map[string]interface{}) (*workload.Workload,error)
}

// InfluxDBWA ...
type InfluxDBWA struct{
	qAPI api.QueryAPI
	ctx context.Context
	cnF context.CancelFunc
	org string
	bucket string
}

// NewInfluxDBWA returns a new InfluxDBWA
func NewInfluxDBWA(url, token, organization, bucket string) (*InfluxDBWA,error){
	if len(strings.Trim(url, " ")) == 0{
		return nil, errors.Errorf("the argument %s cant be empty string", "url")
	}
	if len(strings.Trim(token, " ")) == 0{
		return nil, errors.Errorf("the argument %s cant be empty string", "token")
	}
	if len(strings.Trim(organization, " ")) == 0{
		return nil, errors.Errorf("the argument %s cant be empty string", "organization")
	}
	if len(strings.Trim(bucket, " ")) == 0{
		return nil, errors.Errorf("the argument %s cant be empty string", "bucket")
	}

	ctx, cnF := context.WithCancel(context.Background())
	i := &InfluxDBWA{ctx: ctx, cnF: cnF, org: organization, bucket: bucket}
	i.qAPI = dataaccess.GetNewClientAndQueryAPI(url, token, organization)

	return i, nil
}

// GetWorkload ...
func (i *InfluxDBWA) GetWorkload(startTime, finishTime int64, endpointFilters map[string]map[string]interface{}) (*workload.Workload,error){
	fmt.Println(startTime, finishTime, finishTime- startTime)
	if startTime >= finishTime {
		return nil, errors.Errorf("for getting GetCPUUtilizations(), startTime must be less than finishTime")
	}

	w := make(map[string]string)

	for endpointName, filters := range endpointFilters{
		query := `
from(bucket: "$BUCKET_NAME")
  |> range(start: $START_TIME, stop: $FINISH_TIME)
  |> filter(fn: (r) => r["_measurement"] == "request_info" and r["_field"] == "ust")
  |> filter(fn: (r) => r["method"] == "\"$HTTP_METHOD\"" and r["uri"] =~ /$URI_REGEX/)
  |> group()
`
		query = strings.Replace(query, "$BUCKET_NAME", i.bucket,-1 )
		query = strings.Replace(query, "$START_TIME", strconv.FormatInt(startTime, 10),-1 )
		query = strings.Replace(query, "$FINISH_TIME", strconv.FormatInt(finishTime, 10),-1 )

		httpMethod, ok := filters[strings.ToLower("HTTP_METHOD")]
		if !ok{
			return nil, errors.Errorf("HTTP_METHOD must be provided as filter for GetTargetWorkload %s", endpointName)
		}
		query = strings.Replace(query, "$HTTP_METHOD", httpMethod.(string), -1)

		uriRegex, ok := filters[strings.ToLower("URI_REGEX")]
		if !ok{
			return nil, errors.Errorf("URI_REGEX must be provided as filter for GetTargetWorkload %s", endpointName)
		}
		query = strings.Replace(query, "$URI_REGEX", uriRegex.(string), -1)

		log.Debug("getting request count for " + endpointName + " from influxdb with query:\n" + query)

		_, values, err := dataaccess.QuerySingleTable(i.qAPI, i.ctx, query, "_value")
		if err != nil{
			return nil, errors.Wrap(err, "error getting response times from influxdb using:\n" + query)
		}

		c := len(values)
		w[endpointName] = strconv.FormatInt(int64(math.Round(float64(c) / float64(finishTime - startTime))), 10)
	}

	wc := workload.Workload(w)
	return &wc, nil
}