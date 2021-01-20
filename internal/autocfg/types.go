package autocfg

import (
	"crypto/md5"
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/workload"
	"k8s.io/apimachinery/pkg/util/json"
	"strconv"
	"strings"
)

//TODO add more args to this hash function. probably make it work with ...
func GetHash(c map[string]*configuration.Configuration, version string) (string,error){
	b, err := json.Marshal(c)
	if err != nil{
		return "", errors.Wrap(err, "cant convert configuration to json")
	}
	b = append(b, []byte(version)...)
	log.Debugf("hashing with %s", string(b))
	s := md5.Sum(b)
	h := fmt.Sprintf("%x", s)
	return h, err
}

type IterationInformation struct{
	Configuration map[string]*configuration.Configuration	`yaml:"configurations"`
	StartTime int64							`yaml:"startTime"`
	FinishTime int64						`yaml:"finishTime"`
	AggregatedData *aggregators.AggregatedData			`yaml:"aggregatedData"`
}

type TestInformation struct{
	Name string								`yaml:"name"`
	VersionCode string						`yaml:"version"`
	AutoconfiguringApproach string 			`yaml:"autoConfigApproach"`
	Iterations []*IterationInformation 		`yaml:"iterations"`
	InputWorkload *workload.Workload 		`yaml:"workload"`
	AllSettings	map[string]interface{}		`yaml:"allSettings"`
}

type SLA struct{
	Conditions []Condition
}

type Condition struct{
	Type string					`yaml:"type"`
	EndpointName string			`yaml:"endpointName"`
	Threshold float64			`yaml:"threshold"`
	ComputeFnName string		`yaml:"computeFunctionName"`
}

func (c *Condition) GetComputeFunction() func([]float64) float64{
	if strings.ToLower(c.ComputeFnName) == "mean"{
		return func(values []float64) float64{
			m, err := stats.Mean(values)
			if err != nil{
				panic(err)
			}
			return m
		}
	} else if strings.Contains(strings.ToLower(c.ComputeFnName), "percentile_"){
		return func(values []float64) float64 {
			percent, err := strconv.ParseFloat(strings.Replace(c.ComputeFnName,"percentile_", "",1), 64 )
			if err != nil{
				panic(errors.Wrapf(err,"cant parse %s to get ComputeFn function for condition of SLA. Acceptable example is percentile_90", c.ComputeFnName))
			}
			v, err := stats.Percentile(values, percent)
			return v
		}
	}
	panic(errors.New("unknown computeFnName for SLA condition: " + c.ComputeFnName))
}