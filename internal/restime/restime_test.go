package restime

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"math"
	"reflect"
	"testing"
)

type simpleAggregator struct{}

func (s simpleAggregator) GetResponseTimes(startTime, finishTime float64, filters map[string]interface{}) (*ResponseTimes, error) {
	rts := []float64{1, 1, 1, 1, 1, 2, 2, 2, 2, 2}
	rtsV := ResponseTimes(rts)
	return &rtsV, nil
}

func TestSimpleAggregatorStats(t *testing.T) {
	s := simpleAggregator{}
	responseTimes, err := s.GetResponseTimes(-1, -1, nil)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	mean, err := responseTimes.GetMean()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if math.Abs(mean-1.5) > 1e-4 {
		t.Fail()
	}

	p90, err := responseTimes.GetPercentile(90)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if math.Abs(p90-2) > 1e-4 {
		t.Fail()
	}
}

func TestInfluxDBRTA_GetResponseTimes(t *testing.T) {
	type fields struct {
		qAPI   api.QueryAPI
		ctx    context.Context
		cnF    context.CancelFunc
		org    string
		bucket string
	}
	type args struct {
		startTime  int64
		finishTime int64
		filters    map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ResponseTimes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfluxDBRTA{
				qAPI:   tt.fields.qAPI,
				ctx:    tt.fields.ctx,
				cnF:    tt.fields.cnF,
				org:    tt.fields.org,
				bucket: tt.fields.bucket,
			}
			got, err := i.GetResponseTimes(tt.args.startTime, tt.args.finishTime, tt.args.filters)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResponseTimes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResponseTimes() got = %v, want %v", got, tt.want)
			}
		})
	}
}