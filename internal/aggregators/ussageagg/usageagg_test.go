package ussageagg

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
	"os"
	"testing"
	"time"
)

func TestUsageAggregator_GetAggregatedCPUUtilizations(t *testing.T) {
	viper.Set(constants.EndpointsAggregatorArgsURL, os.Getenv("INFLUXDB_URL"))
	viper.Set(constants.EndpointsAggregatorArgsToken, os.Getenv("INFLUXDB_TOKEN"))
	viper.Set(constants.EndpointsAggregatorArgsOrganization, os.Getenv("INFLUXDB_ORG"))
	viper.Set(constants.EndpointsAggregatorArgsBucket, os.Getenv("INFLUXDB_BUCKET"))

	resourceFilters := map[string]map[string]interface{}{
		"auth": {"POD_NAME_REGEX":"^auth-*"},
		"gateway": {"POD_NAME_REGEX":"^gateway-*"},
		"books": {"POD_NAME_REGEX":"^books-*"},
	}
	viper.Set(constants.ResourceFilters, resourceFilters)
	i,err := NewUsageAggregator("influxdb")
	if err != nil{
		t.Log(err)
		t.Fail()
		return
	}

	res, err := i.GetAggregatedCPUUtilizations(time.Now().Add(-3 * time.Minute).Unix(), time.Now().Add(-1 * time.Minute).Unix())
	if err != nil{
		t.Log(err)
		t.Fail()
		return
	}

	if len(res) != 3{
		t.Log("len(res) must be 3")
		t.Fail()
		return
	}

	for r, v := range res{
		m, err := v.GetMean()
		if err != nil{
			t.Log(err)
			t.Fail()
			return
		}
		p90, err :=v.GetPercentile(90)
		if err != nil{
			t.Log(err)
			t.Fail()
			return
		}
		fmt.Println(r, m,p90)
	}


}
