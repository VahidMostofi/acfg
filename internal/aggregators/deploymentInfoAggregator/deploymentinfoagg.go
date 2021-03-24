package deploymentinfoagg

import (
	"context"
	"strconv"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/vahidmostofi/acfg/internal/dataaccess"
)

type DeploymentInfo struct {
	Name    string
	Replica int
}

type DeploymentInfoAggregator interface {
	GetDeploymentInfo(deploymentName string, startTime, finishTime int64, filters map[string]interface{}) (*DeploymentInfo, error)
}

// InfluxDBDIA gets deployment info from from influxdb
type InfluxDBDIA struct {
	qAPI   api.QueryAPI
	ctx    context.Context
	cnF    context.CancelFunc
	org    string
	bucket string
}

// NewInfluxDBDIA returns a new InfluxDBDIA
func NewInfluxDBDIA(url, token, organization, bucket string) (*InfluxDBDIA, error) {
	if len(strings.Trim(url, " ")) == 0 {
		return nil, errors.Errorf("the argument %s cant be empty string", "url")
	}
	if len(strings.Trim(token, " ")) == 0 {
		return nil, errors.Errorf("the argument %s cant be empty string", "token")
	}
	if len(strings.Trim(organization, " ")) == 0 {
		return nil, errors.Errorf("the argument %s cant be empty string", "organization")
	}
	if len(strings.Trim(bucket, " ")) == 0 {
		return nil, errors.Errorf("the argument %s cant be empty string", "bucket")
	}
	ctx, cnF := context.WithCancel(context.Background())
	i := &InfluxDBDIA{ctx: ctx, cnF: cnF, org: organization, bucket: bucket}
	i.qAPI = dataaccess.GetNewClientAndQueryAPI(url, token, organization)

	return i, nil
}

func (i *InfluxDBDIA) GetDeploymentInfo(deploymentName string, startTime, finishTime int64, filters map[string]interface{}) (*DeploymentInfo, error) {
	if startTime >= finishTime {
		return nil, errors.Errorf("for getting GetCPUUtilizations(), startTime must be less than finishTime")
	}
	query := `
from(bucket: "$BUCKET_NAME")
|> range(start: $START_TIME, stop: $FINISH_TIME)
  |> filter(fn: (r) => r["_measurement"] == "kubernetes_deployment")
  |> filter(fn: (r) => r["_field"] == "replicas_available")
  |> filter(fn: (r) => r["deployment_name"] == "$DEPLOYMENT_NAME")
  |> last()
  |> yield(name: "_value")
`
	query = strings.Replace(query, "$BUCKET_NAME", i.bucket, -1)

	query = strings.Replace(query, "$START_TIME", strconv.FormatInt(startTime, 10), -1)
	query = strings.Replace(query, "$FINISH_TIME", strconv.FormatInt(finishTime, 10), -1)
	query = strings.Replace(query, "$DEPLOYMENT_NAME", deploymentName, -1)

	log.Debug("getting deployment info from influxdb with query:\n" + query)

	_, values, err := dataaccess.QuerySingleTableInt64(i.qAPI, i.ctx, query, "_value")
	if err != nil {
		return nil, errors.Wrap(err, "error getting response times from influxdb using:\n"+query)
	}

	return &DeploymentInfo{Replica: int(values[0])}, nil
}

// NewDeploymenInfoAggregator ...
// available kinds: influxdb
// for influxdb, it uses the url, token, organization and bucket pass them in args which is map[string]interface{}
func NewDeploymentInfoAggregator(kind string, args map[string]interface{}) (DeploymentInfoAggregator, error) {
	u := &InfluxDBDIA{}
	var err error
	if kind == "influxdb" {
		u, err = NewInfluxDBDIA(
			args["url"].(string),
			args["token"].(string),
			args["organization"].(string),
			args["bucket"].(string),
		)

		if err != nil {
			return nil, errors.Wrap(err, "cant create InfluxDBRTA")
		}

	} else {
		return nil, errors.Errorf("unknown kind: %s", kind)
	}

	return u, nil
}
