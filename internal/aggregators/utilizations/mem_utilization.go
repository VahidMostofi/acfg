package utilizations

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/dataaccess"
	"strconv"
	"strings"
	"time"
)

// MemUtilizationAggregator uses some functionality to gather mem utilization values based on some functionality
// between startTime and finishTime. The filters might be used to add selection functionality.
// if startTime is less than 0, it should be replaced with 0; time.Unix()
// if finishTime is less than 0, it should be replaced with current time time.Unix()
// if filters is null, there operation should be done without any filtering.
// available filters: POD_NAME_REGEX
type MemUtilizationAggregator interface {
	GetMemUtilizations(startTime, finishTime int64, filters map[string]interface{}) (*MemUtilizations, error)
	GetMemPsiUtilizations(startTime, finishTime int64, filters map[string]interface{}) (*MemPsiUtilizations, error)
	GetMemUtilizationsWithTimestamp(startTime, finishTime int64, filters map[string]interface{}) ([]MemTimestampedUsage, error)
}

// MemTimestampedUsage ...
type MemTimestampedUsage struct {
	Timestamp      int64   `json:"ts"`
	MemUtilization float64 `json:"mem"`
}

// MemUtilizations is a named type for []float64 with helper functions
// values are between 0-100 (percentage utilization / allocation)
type MemUtilizations []float64

// GetMean returns the average of mem utilizations
func (rts *MemUtilizations) GetMean() (float64, error) {
	mean, err := stats.Mean([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing mean failed on mem utilizations")
	}
	return mean, nil
}

// GetMedian returns the median of mem utilizations
func (rts *MemUtilizations) GetMedian() (float64, error) {
	med, err := stats.Mean([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing median failed on mem utilizations")
	}
	return med, nil
}

// GetStd returns the std of mem utilizations
func (rts *MemUtilizations) GetStd() (float64, error) {
	std, err := stats.StandardDeviation([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing std failed on mem utilizations")
	}
	return std, nil
}

// GetPercentile returns the percentile of mem utilizations
func (rts *MemUtilizations) GetPercentile(p float64) (float64, error) {
	p, err := stats.Percentile([]float64(*rts), p)
	if err != nil {
		return 0, errors.Wrap(err, "computing percentile failed on mem utilizations")
	}
	return p, nil
}

// MemPsiUtilizations region Mem PSI
type MemPsiUtilizations []float64

// GetMean returns the average of mem utilizations
func (rts *MemPsiUtilizations) GetMean() (float64, error) {
	mean, err := stats.Mean([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing mean failed on mem psi utilizations")
	}
	return mean, nil
}

// GetMedian returns the median of mem utilizations
func (rts *MemPsiUtilizations) GetMedian() (float64, error) {
	med, err := stats.Mean([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing median failed on mem psi utilizations")
	}
	return med, nil
}

// GetStd returns the std of mem utilizations
func (rts *MemPsiUtilizations) GetStd() (float64, error) {
	std, err := stats.StandardDeviation([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing std failed on mem psi utilizations")
	}
	return std, nil
}

// GetPercentile returns the percentile of mem utilizations
func (rts *MemPsiUtilizations) GetPercentile(p float64) (float64, error) {
	p, err := stats.Percentile([]float64(*rts), p)
	if err != nil {
		return 0, errors.Wrap(err, "computing percentile failed on mem psi utilizations")
	}
	return p, nil
}

// InfluxDBMemUA gets mem utilizations from influxdb
type InfluxDBMemUA struct {
	qAPI   api.QueryAPI
	ctx    context.Context
	cnF    context.CancelFunc
	org    string
	bucket string
}

// NewInfluxDBMemUA returns a new NewInfluxDBMemUA
func NewInfluxDBMemUA(url, token, organization, bucket string) (*InfluxDBMemUA, error) {
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
	i := &InfluxDBMemUA{ctx: ctx, cnF: cnF, org: organization, bucket: bucket}
	i.qAPI = dataaccess.GetNewClientAndQueryAPI(url, token, organization)

	return i, nil
}

func (i *InfluxDBMemUA) buildQuery(startTime, finishTime int64, filters map[string]interface{}) string {
	query := `
	data_total = from(bucket: "$BUCKET_NAME")
	 |> range(start: $START_TIME, stop: $FINISH_TIME)
	 |> filter(fn: (r) => r["_measurement"] == "kubernetes_pod_container")
	 |> filter(fn: (r) => r["_field"] == "resource_limits_memory_bytes")
	 |> filter(fn: (r) => r["state"] == "running")
	 |> rename(columns: {pod_name: "podName"})
	 |> filter(fn: (r) => r["podName"] =~ /$POD_NAME_REGEX/)
	 |> keep(columns: ["_time","_value","podName"])
	
	data_usage = from(bucket: "$BUCKET_NAME")
	 |> range(start: $START_TIME, stop: $FINISH_TIME)
	 |> filter(fn: (r) => r["_measurement"] == "resource_usage")
	 |> filter(fn: (r) => r["_field"] == "memory")
	 |> aggregateWindow(every: 10s, fn: mean)
	 |> keep(columns: ["_time","_value","podName"])
	
	joined = join(
	 tables: {d1: data_total, d2: data_usage},
	 on: ["_time","podName"], method: "inner"
	)
	 |> filter(fn: (r) => (exists r["_value_d1"]) and (exists r["_value_d2"]))
	 |> map(fn:(r) => ({ r with _value_d1: float(v: r._value_d1) }))
     |> map(fn: (r) => ({ r with _value: ((r._value_d2 / 1000.0) / r._value_d1 )* 100.0 }))
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
		panic(errors.Errorf("need POD_NAME_REGEX in filters to get mem utilizations"))
	}

	log.Debug("getting mem utilizations from influxdb with query:\n" + query)

	return query
}

// GetMemUtilizations ....
func (i *InfluxDBMemUA) GetMemUtilizations(startTime, finishTime int64, filters map[string]interface{}) (*MemUtilizations, error) {
	if startTime >= finishTime {
		return nil, errors.Errorf("for getting GetMemUtilizations(), startTime must be less than finishTime")
	}

	query := i.buildQuery(startTime, finishTime, filters)
	_, values, err := dataaccess.QuerySingleTable(i.qAPI, i.ctx, query, "_value")
	if err != nil {
		return nil, errors.Wrap(err, "error getting mem utilizations from influxdb using:\n"+query)
	}

	r := MemUtilizations(values)

	return &r, nil
}

// GetMemPsiUtilizations ....
func (i *InfluxDBMemUA) GetMemPsiUtilizations(startTime, finishTime int64, filters map[string]interface{}) (*MemPsiUtilizations, error) {
	// Not supported
	log.Warning("InfluxDBMemUA does not support mem PSI")
	return nil, nil
}

// GetMemUtilizationsWithTimestamp ....
func (i *InfluxDBMemUA) GetMemUtilizationsWithTimestamp(startTime, finishTime int64, filters map[string]interface{}) ([]MemTimestampedUsage, error) {
	if startTime >= finishTime {
		return nil, errors.Errorf("for getting GetMemUtilizationsWithTimestamp(), startTime must be less than finishTime")
	}

	query := i.buildQuery(startTime, finishTime, filters)

	times, values, err := dataaccess.QuerySingleTable(i.qAPI, i.ctx, query, "_value")
	if err != nil {
		return nil, errors.Wrap(err, "error getting mem utilizations from influxdb using:\n"+query)
	}

	r := make([]MemTimestampedUsage, len(times))
	for i, _ := range times {
		r[i] = MemTimestampedUsage{Timestamp: times[i].Unix(), MemUtilization: values[i]}
	}

	return r, nil
}

type PromDBMemUA struct {
	api    v1.API
	ctx    context.Context
	cnF    context.CancelFunc
	org    string
	bucket string
}

func NewPromDBMemUA(url, token, organization, bucket string) (*PromDBMemUA, error) {
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
	i := &PromDBMemUA{ctx: ctx, cnF: cnF, org: organization, bucket: bucket}
	i.api = dataaccess.GetNewPromClientAndQueryAPI(url)

	return i, nil
}

// GetMemUtilizations ....
func (i *PromDBMemUA) GetMemUtilizations(startTime, finishTime int64, filters map[string]interface{}) (*MemUtilizations, error) {
	if startTime >= finishTime {
		return nil, errors.Errorf("for getting GetMemUtilizations(), startTime must be less than finishTime")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v1.Range{
		Start: time.Unix(startTime, 0),
		End:   time.Unix(finishTime, 0),
		Step:  time.Second * 10,
	}

	query := `sum(node_namespace_pod_container:container_memory_working_set_bytes{pod=~"$POD_NAME_REGEX"}) / sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_requests{pod=~"$POD_NAME_REGEX"}) * 100`
	if podNameRegex, ok := filters[strings.ToLower("POD_NAME_REGEX")]; ok {
		query = strings.Replace(query, "$POD_NAME_REGEX", podNameRegex.(string), -1)
	} else {
		panic(errors.Errorf("need POD_NAME_REGEX in filters to get mem utilizations"))
	}

	log.Debug("getting Mem utilizations from prom with query:\n" + query)

	result, warnings, err := i.api.QueryRange(ctx, query, r)
	if err != nil {
		log.Error("Error querying Prometheus: %v\n", err)
	}
	if len(warnings) > 0 {
		log.Error("Warnings: %v\n", warnings)
	}

	resultMat, _ := result.(model.Matrix)

	var resultSlice []float64
	if len(resultMat) == 0 {
		log.Error("No Mem ute values. assuming 100")
		return (*MemUtilizations)(&[]float64{100}), nil
	}
	for _, kv := range resultMat[0].Values {
		resultSlice = append(resultSlice, float64(kv.Value))
	}

	return (*MemUtilizations)(&resultSlice), nil
}

// GetMemUtilizationsWithTimestamp ....
func (i *PromDBMemUA) GetMemUtilizationsWithTimestamp(startTime, finishTime int64, filters map[string]interface{}) ([]MemTimestampedUsage, error) {
	if startTime >= finishTime {
		return nil, errors.Errorf("for getting GetMemUtilizationsWithTimestamp(), startTime must be less than finishTime")
	}
	log.Warning("PromDBMemUA does not support mem timestamps")
	return nil, nil
}

// GetMemPsiUtilizations ....
func (i *PromDBMemUA) GetMemPsiUtilizations(startTime, finishTime int64, filters map[string]interface{}) (*MemPsiUtilizations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v1.Range{
		Start: time.Unix(startTime, 0),
		End:   time.Unix(finishTime, 0),
		Step:  time.Second * 10,
	}

	query := `cgroup_monitor_sc_monitored_mem_psi{type="some",window="10s",job=~"$POD_NAME_REGEX"}`

	if podNameRegex, ok := filters[strings.ToLower("POD_NAME_REGEX")]; ok {
		query = strings.Replace(query, "$POD_NAME_REGEX", podNameRegex.(string), -1)
	} else {
		panic(errors.Errorf("need POD_NAME_REGEX in filters to get Mem psi utilizations"))
	}

	log.Debug("getting mem psi utilizations from prom with query:\n" + query)

	result, warnings, err := i.api.QueryRange(ctx, query, r)
	if err != nil {
		log.Error("Error querying Prometheus: %v\n", err)
	}
	if len(warnings) > 0 {
		log.Error("Warnings: %v\n", warnings)
	}
	resultMat, _ := result.(model.Matrix)

	var resultSlice []float64
	for _, kv := range resultMat[0].Values {
		resultSlice = append(resultSlice, float64(kv.Value))
	}

	return (*MemPsiUtilizations)(&resultSlice), nil
}
