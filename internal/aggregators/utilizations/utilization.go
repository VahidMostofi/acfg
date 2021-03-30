package utilizations

import (
	"context"
	"strconv"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/dataaccess"
)

// CPUUtilizationAggregator uses some functionality to gather CPU utilization values based on some functionality
// between startTime and finishTime. The filters might be used to add selection functionality.
// if startTime is less than 0, it should be replaced with 0; time.Unix()
// if finishTime is less than 0, it should be replaced with current time time.Unix()
// if filters is null, there operation should be done without any filtering.
// available filters: POD_NAME_REGEX
type CPUUtilizationAggregator interface {
	GetCPUUtilizations(startTime, finishTime int64, filters map[string]interface{}) (*CPUUtilizations, error)
	GetCPUUtilizationsWithTimestamp(startTime, finishTime int64, filters map[string]interface{}) ([]TimestampedUsage, error)
}

type TimestampedUsage struct {
	Timestamp      int64   `json:"ts"`
	CPUUtilization float64 `json:"cpu"`
}

// CPUUtilizations is a named type for []float64 with helper functions
// values are between 0-100 (they can be more than 100 due to something that happens in linux)
type CPUUtilizations []float64

// GetMean returns the average of CPU utilizations
func (rts *CPUUtilizations) GetMean() (float64, error) {
	mean, err := stats.Mean([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing mean failed on CPU utilizations")
	}
	return mean, nil
}

// GetMedian returns the median of CPU utilizations
func (rts *CPUUtilizations) GetMedian() (float64, error) {
	med, err := stats.Mean([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing median failed on CPU utilizations")
	}
	return med, nil
}

// GetStd returns the std of CPU utilizations
func (rts *CPUUtilizations) GetStd() (float64, error) {
	std, err := stats.StandardDeviation([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing std failed on CPU utilizations")
	}
	return std, nil
}

// GetPercentile returns the percentile of CPU utilizations
func (rts *CPUUtilizations) GetPercentile(p float64) (float64, error) {
	p, err := stats.Percentile([]float64(*rts), p)
	if err != nil {
		return 0, errors.Wrap(err, "computing percentile failed on CPU utilizations")
	}
	return p, nil
}

// InfluxDBRTA gets CPU utilizations from influxdb
type InfluxDBCPUUA struct {
	qAPI   api.QueryAPI
	ctx    context.Context
	cnF    context.CancelFunc
	org    string
	bucket string
}

// NewInfluxDBCPUUA returns a new InfluxDBCPUUA
func NewInfluxDBCPUUA(url, token, organization, bucket string) (*InfluxDBCPUUA, error) {
	if len(strings.Trim(url, " ")) == 0 {
		return nil, errors.Errorf("the argument %s cant be empty string", "url")
	}
	if len(strings.Trim(token, " ")) == 0 {
		return nil, errors.Errorf("the argument %s cant be empty string", "token")
	}
	if len(strings.Trim(organization, " ")) == 0 {
		return nil, errors.Errorf("the argument %s cant be empty string", "organization")
	}
	if len(strings.Trim(bucket, " ")) == 0 {
		return nil, errors.Errorf("the argument %s cant be empty string", "bucket")
	}
	ctx, cnF := context.WithCancel(context.Background())
	i := &InfluxDBCPUUA{ctx: ctx, cnF: cnF, org: organization, bucket: bucket}
	i.qAPI = dataaccess.GetNewClientAndQueryAPI(url, token, organization)

	return i, nil
}

func (i *InfluxDBCPUUA) buildQuery(startTime, finishTime int64, filters map[string]interface{}) string {
	query := `
	data_total = from(bucket: "$BUCKET_NAME")
	 |> range(start: $START_TIME, stop: $FINISH_TIME)
	 |> filter(fn: (r) => r["_measurement"] == "kubernetes_pod_container")
	 |> filter(fn: (r) => r["_field"] == "resource_limits_millicpu_units")
	 |> filter(fn: (r) => r["state"] == "running")
	 |> rename(columns: {pod_name: "podName"})
	 |> filter(fn: (r) => r["podName"] =~ /$POD_NAME_REGEX/)
	 |> keep(columns: ["_time","_value","podName"])
	
	data_usage = from(bucket: "$BUCKET_NAME")
	 |> range(start: $START_TIME, stop: $FINISH_TIME)
	 |> filter(fn: (r) => r["_measurement"] == "resource_usage")
	 |> filter(fn: (r) => r["_field"] == "cpu")
	 |> aggregateWindow(every: 10s, fn: mean)
	 |> keep(columns: ["_time","_value","podName"])
	
	joined = join(
	 tables: {d1: data_total, d2: data_usage},
	 on: ["_time","podName"], method: "inner"
	)
	 |> filter(fn: (r) => (exists r["_value_d1"]) and (exists r["_value_d2"]))
	 |> map(fn:(r) => ({ r with _value_d1: float(v: r._value_d1) }))
	 |> map(fn: (r) => ({ r with _value: (r._value_d2 / r._value_d1 )* 100.0 }))
	 |> group(columns: ["_time", "podName"], mode: "by")
	 |> group()
	 |> aggregateWindow(every: 10s, fn: mean)
	joined
	`
	query = strings.Replace(query, "$BUCKET_NAME", i.bucket, -1)
	query = strings.Replace(query, "$START_TIME", strconv.FormatInt(startTime, 10), -1)
	query = strings.Replace(query, "$FINISH_TIME", strconv.FormatInt(finishTime, 10), -1)

	if podNameRegex, ok := filters[strings.ToLower("POD_NAME_REGEX")]; ok {
		query = strings.Replace(query, "$POD_NAME_REGEX", podNameRegex.(string), -1)
	} else {
		panic(errors.Errorf("need POD_NAME_REGEX in filters to get CPU utilizations"))
	}

	log.Debug("getting CPU utilizations from influxdb with query:\n" + query)

	return query
}

// GetCPUUtilizations ....
func (i *InfluxDBCPUUA) GetCPUUtilizations(startTime, finishTime int64, filters map[string]interface{}) (*CPUUtilizations, error) {
	if startTime >= finishTime {
		return nil, errors.Errorf("for getting GetCPUUtilizations(), startTime must be less than finishTime")
	}

	query := i.buildQuery(startTime, finishTime, filters)
	_, values, err := dataaccess.QuerySingleTable(i.qAPI, i.ctx, query, "_value")
	if err != nil {
		return nil, errors.Wrap(err, "error getting CPU utilizations from influxdb using:\n"+query)
	}

	r := CPUUtilizations(values)

	return &r, nil
}

// GetCPUUtilizations ....
func (i *InfluxDBCPUUA) GetCPUUtilizationsWithTimestamp(startTime, finishTime int64, filters map[string]interface{}) ([]TimestampedUsage, error) {
	if startTime >= finishTime {
		return nil, errors.Errorf("for getting GetCPUUtilizations(), startTime must be less than finishTime")
	}

	query := i.buildQuery(startTime, finishTime, filters)

	times, values, err := dataaccess.QuerySingleTable(i.qAPI, i.ctx, query, "_value")
	if err != nil {
		return nil, errors.Wrap(err, "error getting CPU utilizations from influxdb using:\n"+query)
	}

	r := make([]TimestampedUsage, len(times))
	for i, _ := range times {
		r[i] = TimestampedUsage{Timestamp: times[i].Unix(), CPUUtilization: values[i]}
	}

	return r, nil
}
