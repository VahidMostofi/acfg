package autocfg

import (
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/workload"
	"k8s.io/apimachinery/pkg/util/json"
)

func GetHash(c map[string]*configuration.Configuration, version string) (string,error){
	b, err := json.Marshal(c)
	if err != nil{
		return "", errors.Wrap(err, "cant convert configuration to json")
	}
	b = append(b, []byte(version)...)
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
}

type SLA struct{
	Conditions []Condition
}

type Condition struct{
	Type string
	EndpointName string
	Threshold float64
	ComputeFn func([]float64) float64
}