package dataaccess

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// GetNewClientAndQueryAPI returns new ...
func GetNewClientAndQueryAPI(url, token string) (influxdb2.Client, error) {
	client := influxdb2.NewClient(url, token)

	return client, nil
}
