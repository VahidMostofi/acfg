package restime

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/dataaccess"
)

// ResponseTimeAggregator uses some functionality to gather response times based on some functionality
// between startTime and finishTime. The filters might be used to add selection functionality.
// if startTime is less than 0, it should be replaced with 0; time.Unix()
// if finishTime is less than 0, it should be replaced with current time time.Unix()
// if filters is null, there operation should be done without any filtering.
// available filters: HTTP_METHOD, URI_REGEX
type ResponseTimeAggregator interface {
	GetResponseTimes(startTime, finishTime int64, filters map[string]interface{}) (*ResponseTimes, error)
	GetTimeStampedResponseTimes(startTime, finishTime int64, filters map[string]interface{}) ([]TimestampedResponseTime, error)
}

type TimestampedResponseTime struct {
	Timestamp    int64   `json:"ts"`
	ResponseTime float64 `json:"rs"`
}

// ResponseTimes is a named type for []float64 with helper functions
type ResponseTimes []float64

func (rts *ResponseTimes) String() string {
	m, err := rts.GetMean()
	if err != nil {
		panic(err)
	}
	p90, err := rts.GetPercentile(90)
	if err != nil {
		panic(err)
	}
	p95, err := rts.GetPercentile(95)
	if err != nil {
		panic(err)
	}
	p99, err := rts.GetPercentile(99)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("count: %d, mean: %fs, 90P: %fs, 95P: %fs, 99P: %fs", rts.GetCount(), m, p90, p95, p99)
}

// GetMean returns the average of response times
func (rts *ResponseTimes) GetMean() (float64, error) {
	mean, err := stats.Mean([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing mean failed on response times")
	}
	return mean, nil
}

// GetMedian returns the median of response times
func (rts *ResponseTimes) GetMedian() (float64, error) {
	med, err := stats.Mean([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing median failed on response times")
	}
	return med, nil
}

// GetStd returns the std of response times
func (rts *ResponseTimes) GetStd() (float64, error) {
	std, err := stats.StandardDeviation([]float64(*rts))
	if err != nil {
		return 0, errors.Wrap(err, "computing std failed on response times")
	}
	return std, nil
}

// GetPercentile returns the percentile of response times
func (rts *ResponseTimes) GetPercentile(p float64) (float64, error) {
	p, err := stats.Percentile([]float64(*rts), p)
	if err != nil {
		return 0, errors.Wrap(err, "computing percentile failed on response times")
	}
	return p, nil
}

// GetCount returns the number of response times
func (rts *ResponseTimes) GetCount() int {
	return len([]float64(*rts))
}

// InfluxDBRTA gets response times from influxdb
type InfluxDBRTA struct {
	qAPI   api.QueryAPI
	ctx    context.Context
	cnF    context.CancelFunc
	org    string
	bucket string
}

// NewInfluxDBRTA returns a new InfluxDBRTA
func NewInfluxDBRTA(url, token, organization, bucket string) (*InfluxDBRTA, error) {
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
	i := &InfluxDBRTA{ctx: ctx, cnF: cnF, org: organization, bucket: bucket}
	i.qAPI = dataaccess.GetNewClientAndQueryAPI(url, token, organization)

	return i, nil
}

func (i *InfluxDBRTA) buildQuery(startTime, finishTime int64, filters map[string]interface{}) string {
	query := `
	from(bucket: "$BUCKET_NAME")
		|> range(start: $START_TIME, stop: $FINISH_TIME)
		|> filter(fn: (r) => r["_measurement"] == "request_info" and r._field == "ust" $CONDITIONS)
		|> group()
		|> keep(columns: ["_time", "_value"])
	`
	query = strings.Replace(query, "$BUCKET_NAME", i.bucket, -1)
	query = strings.Replace(query, "$START_TIME", strconv.FormatInt(startTime, 10), -1)
	query = strings.Replace(query, "$FINISH_TIME", strconv.FormatInt(finishTime, 10), -1)

	conditions := ""
	if httpMethod, ok := filters[strings.ToLower("HTTP_METHOD")]; ok {
		if conditions == "" {
			conditions = " and "
		}
		conditions += "r.method == \"" + httpMethod.(string) + "\""
	}
	if uriRegex, ok := filters[strings.ToLower("URI_REGEX")]; ok {
		conditions += " and r.uri =~ /" + uriRegex.(string) + "/"
	}
	query = strings.Replace(query, "$CONDITIONS", conditions, -1)

	log.Debug("getting response times from influxdb with query:\n" + query)

	return query
}

func (i *InfluxDBRTA) GetTimeStampedResponseTimes(startTime, finishTime int64, filters map[string]interface{}) ([]TimestampedResponseTime, error) {
	if startTime >= finishTime {
		return nil, errors.Errorf("for getting GetCPUUtilizations(), startTime must be less than finishTime")
	}

	query := i.buildQuery(startTime, finishTime, filters)

	times, values, err := dataaccess.QuerySingleTable(i.qAPI, i.ctx, query, "_value")
	if err != nil {
		return nil, errors.Wrap(err, "error getting response times from influxdb using:\n"+query)
	}

	res := make([]TimestampedResponseTime, len(times))
	for i, _ := range times {
		res[i] = TimestampedResponseTime{Timestamp: times[i].Unix(), ResponseTime: values[i]}
	}

	return res, nil
}

// GetResponseTimes ....
func (i *InfluxDBRTA) GetResponseTimes(startTime, finishTime int64, filters map[string]interface{}) (*ResponseTimes, error) {
	if startTime >= finishTime {
		return nil, errors.Errorf("for getting GetCPUUtilizations(), startTime must be less than finishTime")
	}

	query := i.buildQuery(startTime, finishTime, filters)

	_, values, err := dataaccess.QuerySingleTable(i.qAPI, i.ctx, query, "_value")
	if err != nil {
		return nil, errors.Wrap(err, "error getting response times from influxdb using:\n"+query)
	}

	r := ResponseTimes(values)

	return &r, nil
}
