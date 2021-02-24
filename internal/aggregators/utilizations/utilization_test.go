package utilizations

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestInfluxDBCPUUA_GetCPUUtilizations(t *testing.T) {
	i,err := NewInfluxDBCPUUA(os.Getenv("INFLUXDB_URL"), os.Getenv("INFLUXDB_TOKEN"), os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))
	if err != nil{
		t.Log(err)
		t.Fail()
		return
	}
	values, err := i.GetCPUUtilizations(time.Now().Add(-3 * time.Minute).Unix(), time.Now().Add(-1 * time.Minute).Unix(), map[string]interface{}{"POD_NAME_REGEX":"^auth-*"})
	if err != nil{
		t.Log(err)
		t.Fail()
		return
	}
	fmt.Println(values)
	mean, err := values.GetMean()
	if err != nil{
		t.Log(err)
		t.Fail()
		return
	}
	p90, err := values.GetPercentile(90)
	if err != nil{
		t.Log(err)
		t.Fail()
		return
	}
	fmt.Println(mean, p90)
}
