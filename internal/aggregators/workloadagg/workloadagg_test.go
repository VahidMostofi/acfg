package workloadagg

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"os"
	"testing"
	"time"
)

func TestInfluxDBWA_GetWorkload(t *testing.T) {
	viper.Set(constants.EndpointsAggregatorArgsURL, os.Getenv("INFLUXDB_URL"))
	viper.Set(constants.EndpointsAggregatorArgsToken, os.Getenv("INFLUXDB_TOKEN"))
	viper.Set(constants.EndpointsAggregatorArgsOrganization, os.Getenv("INFLUXDB_ORG"))
	viper.Set(constants.EndpointsAggregatorArgsBucket, os.Getenv("INFLUXDB_BUCKET"))

	endpointsFilters := map[string]map[string]interface{}{
		"login": {"URI_REGEX":"login*", "HTTP_METHOD":"POST"},
		"get-book": {"URI_REGEX":"books*", "HTTP_METHOD":"GET"},
		"edit-book": {"URI_REGEX":"books*", "HTTP_METHOD":"PUT"},
	}

	wag,err := NewInfluxDBWA()
	if err != nil{
		t.Log(err)
		t.Fail()
		return
	}

	w, err := wag.GetWorkload(time.Now().Add(-3 *time.Minute).Unix(), time.Now().Add(-1 * time.Minute).Unix(), endpointsFilters)
	if err != nil{
		t.Log(err)
		t.Fail()
		return
	}

	for e,v := range w.GetMap(){
		fmt.Println(e,v)
	}
}
