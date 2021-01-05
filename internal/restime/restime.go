package restime

import (
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
)

// ResponseTimeAggregator uses some functionality to gather response times based on some functionality
// between startTime and finishTime. The filters might be used to add selection functionality.
// if startTime is less than 0, it should be replaced with 0
// if finishTime is less than 0, it should be replaced with current time
// if filters is null, there operation should be done without any filtering.
type ResponseTimeAggregator interface {
	GetResponseTimes(startTime, finishTime float64, filters map[string]interface{}) (*ResponseTimes, error)
}

// ResponseTimes is a named type for []float64 with helper functions
type ResponseTimes []float64

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
