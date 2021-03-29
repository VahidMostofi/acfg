package historyextractor

import (
	"encoding/json"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
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

	startTime := viper.GetInt64(constants.DumpStartTime)
	finishTime := viper.GetInt64(constants.DumpFinishTime)
	outputPath := viper.GetString(constants.DumpOutputPath)
	durationString := (time.Second * time.Duration(finishTime-startTime)).String()

	log.Infof("Dumping from %d to %d (%s), to %s", startTime, finishTime, durationString, outputPath)
	aggData, err := e.DumpDataWithTimestamp(startTime, finishTime)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(outputPath, aggData, 0777)
}
