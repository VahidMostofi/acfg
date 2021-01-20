package loadgenerator

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/aggregators/workloadagg"
	"github.com/vahidmostofi/acfg/internal/constants"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestK6LocalLoadGenerator_Start(t *testing.T) {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("ACFG")

	log.SetLevel(log.DebugLevel)
	endpointsFilters := map[string]map[string]interface{}{
		"login": {"uri_regex":"login*", "http_method":"POST"},
		"get-book": {"uri_regex":"books*", "http_method":"GET"},
		"edit-book": {"uri_regex":"books*", "http_method":"PUT"},
	}

	wa, err := newWorkloadAggregator()
	if err != nil{
		panic(err)
	}
	// ---------------------
	var k = &K6LocalLoadGenerator{}
	f, err := os.Open("/home/vahid/Desktop/temp.js")
	if err != nil{
		panic(err)
	}
	d, err := ioutil.ReadAll(f)
	k.Data = d
	// ---------------------
	for i := 0;i < 4; i++{
		err = k.Start(nil, nil)
		if err != nil{
			panic(err)
			t.Fail()
			return
		}
		// ---------------------
		startTime := time.Now().Unix()
		time.Sleep(30 * time.Second)
		finishTime := time.Now().Unix()
		// ---------------------
		err = k.Stop()
		if err != nil{
			panic(err)
			t.Fail()
			return
		}
		// ---------------------
		time.Sleep(12 * time.Second)
		// ---------------------
		w, err := wa.GetWorkload(startTime, finishTime, endpointsFilters)
		if err != nil{
			panic(err)
		}
		fmt.Println(w.String())
		// ---------------------
		fb,_ := k.GetFeedback()
		fmt.Println(fb["data"])
	}
}

func newWorkloadAggregator()(workloadagg.WorkloadAggregator, error){
	if viper.GetString(constants.WorkloadAggregatorType) == "influxdb"{
		url := viper.GetString(constants.WorkloadAggregatorArgsURL)
		token := viper.GetString(constants.WorkloadAggregatorArgsToken)
		organization := viper.GetString(constants.WorkloadAggregatorArgsOrganization)
		bucket := viper.GetString(constants.WorkloadAggregatorArgsBucket)
		wg, err := workloadagg.NewInfluxDBWA(url, token, organization, bucket)
		if err != nil{
			return nil, errors.Wrapf(err, "error while creating workload aggregator with %s %s %s %s", url, "some token", organization, bucket)
		}
		return wg, err
	}
	return nil, errors.New("unknown worker aggregator type: " + viper.GetString(constants.WorkloadAggregatorType))
}

func parseMapMapInterface(in map[string]interface{}) (map[string]map[string]interface{},error){
	res := make(map[string]map[string]interface{})
	for key,value := range in{
		v, ok := value.(map[string]interface{})
		if !ok{
			return nil, errors.New(fmt.Sprintf("cant convert %s to %s", reflect.TypeOf(v), "map[string]interface{}"))
		}
		res[key] = v
	}
	return res, nil
}