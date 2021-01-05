package restime

import (
	"fmt"
	"math"
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
