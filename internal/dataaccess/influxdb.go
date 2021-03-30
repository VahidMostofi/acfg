package dataaccess

import (
	"context"
	"crypto/tls"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/pkg/errors"
)

// GetNewQueryAPI returns new QueryAPI
func GetNewClientAndQueryAPI(url, token, organization string) api.QueryAPI {
	client := getNewClient(url, token)
	queryAPI := client.QueryAPI(organization)
	return queryAPI
}

func QuerySingleTableInt64(queryAPI api.QueryAPI, ctx context.Context, q, valueName string) ([]time.Time, []int64, error) {
	res, err := queryAPI.Query(ctx, q)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error running query: "+q)
	}

	times := make([]time.Time, 0)
	values := make([]int64, 0)

	for res.Next() {
		t := res.Record().Time()
		v := res.Record().ValueByKey(valueName)
		if v != nil {
			times = append(times, t)
			values = append(values, v.(int64))
		}
	}

	return times, values, nil
}

func QuerySingleTable(queryAPI api.QueryAPI, ctx context.Context, q, valueName string) ([]time.Time, []float64, error) {
	res, err := queryAPI.Query(ctx, q)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error running query: "+q)
	}

	times := make([]time.Time, 0)
	values := make([]float64, 0)

	for res.Next() {
		t := res.Record().Time()
		v := res.Record().ValueByKey(valueName)
		if v != nil {
			times = append(times, t)
			values = append(values, v.(float64))
		}
	}

	return times, values, nil
}

// getNewClient returns new Client
func getNewClient(url, token string) influxdb2.Client {
	options := influxdb2.DefaultOptions()
	options.SetTLSConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	options.SetHTTPRequestTimeout(13600)
	return influxdb2.NewClientWithOptions(url, token, options)
}
