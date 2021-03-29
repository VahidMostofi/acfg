package historyextractor

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/vahidmostofi/acfg/internal/factory"
)

func ExtractHistory() {
	var err error
	e, err := factory.NewEnsembleAggregator(factory.EnsembleAggregatorArgs{
		WithEndpointsAggregator:  true,
		WithWorkloadAggregator:   true,
		WithUsageAggregator:      true,
		WithDeploymentAggregator: true,
	})

	if err != nil {
		panic(err)
	}

	finishTime := time.Now().Unix()
	startTime := finishTime - 1*int64(time.Hour.Seconds())
	aggData, err := e.AggregateData(startTime, finishTime)
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(aggData)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("/home/vahid/Desktop/test_extract_24hours.json", b, 0777)
}

func DumpHistory() {
	var err error
	e, err := factory.NewEnsembleAggregator(factory.EnsembleAggregatorArgs{
		WithEndpointsAggregator:  true,
		WithWorkloadAggregator:   true,
		WithUsageAggregator:      true,
		WithDeploymentAggregator: true,
	})

	if err != nil {
		panic(err)
	}

	finishTime := time.Now().Unix()
	startTime := finishTime - 24*int64(time.Hour.Seconds())
	aggData, err := e.DumpDataWithTimestamp(startTime, finishTime)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("/home/vahid/Desktop/test_extract_24hours.json", aggData, 0777)
}
