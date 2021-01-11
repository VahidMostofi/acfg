package dataaccess

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	fmt.Println(time.Now().Unix())
	qAPi := GetNewClientAndQueryAPI(os.Getenv("INFLUXDB_URL"), os.Getenv("INFLUXDB_TOKEN"), os.Getenv("INFLUXDB_ORG"))
//	q1 := `
//data_total = from(bucket: "general")
//  |> range(start: -2m, stop: 0m)
//  |> filter(fn: (r) => r["_measurement"] == "kubernetes_pod_container")
//  |> filter(fn: (r) => r["_field"] == "resource_limits_millicpu_units")
//  |> filter(fn: (r) => r["state"] == "running")
//  |> rename(columns: {pod_name: "podName"})
//  |> filter(fn: (r) => r["podName"] =~ /^auth*/)
//  |> keep(columns: ["_time","_value","podName"])
//
//data_usage = from(bucket: "general")
//  |> range(start: -2m, stop: 0m)
//  |> filter(fn: (r) => r["_measurement"] == "resource_usage")
//  |> filter(fn: (r) => r["_field"] == "cpu")
//  |> aggregateWindow(every: 10s, fn: mean)
//  |> keep(columns: ["_time","_value","podName"])
//
//joined = join(
//  tables: {d1: data_total, d2: data_usage},
//  on: ["_time","podName"], method: "inner"
//)
//  |> filter(fn: (r) => (exists r["_value_d1"]) and (exists r["_value_d2"]))
//  |> map(fn:(r) => ({ r with _value_d1: float(v: r._value_d1) }))
//  |> map(fn: (r) => ({ r with _value: (r._value_d2 / r._value_d1 )* 100.0 }))
//  |> group(columns: ["_time", "podName"], mode: "by")
//  |> group()
//  |> aggregateWindow(every: 10s, fn: mean)
//joined
//`

	times, values ,err := QuerySingleTable(qAPi, context.Background(), q1, "_value")
	if err != nil{
		t.Log(err)
		t.Fail()
	}

	minTime := times[0]
	maxTime := times[0]

	for _, time := range times{
		if time.UnixNano() < minTime.UnixNano(){
			minTime = time
		}
		if time.UnixNano() > maxTime.UnixNano(){
			maxTime = time
		}
	}

	fmt.Println(maxTime.Unix() - minTime.Unix(), len(values))
}
