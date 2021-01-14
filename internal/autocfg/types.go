package autocfg

import (
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vahidmostofi/acfg/internal/aggregators/restime"
	"github.com/vahidmostofi/acfg/internal/aggregators/sysstructureagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/utilizations"
	"github.com/vahidmostofi/acfg/internal/workload"
	"k8s.io/apimachinery/pkg/util/json"
)

type Configuration struct{
	ResourceType string
	ReplicaCount *int64
	CPU *int64
	Memory *int64
	EnvironmentValues map[string]string
}

func GetHash(c map[string]*Configuration, version string) (string,error){
	b, err := json.Marshal(c)
	if err != nil{
		return "", errors.Wrap(err, "cant convert configuration to json")
	}
	b = append(b, []byte(version)...)
	s := md5.Sum(b)
	h := fmt.Sprintf("%x", s)
	return h, err
}

type AggregatedData struct{
	ResponseTimes map[string]*restime.ResponseTimes				`yaml:"responseTimes"`
	CPUUtilizations map[string]*utilizations.CPUUtilizations	`yaml:"CPUUtilizations"`
	SystemStructure *sysstructureagg.SystemStructure			`yaml:"structure"`
	HappenedWorkload *workload.Workload							`yaml:"workload"`
}

type IterationInformation struct{
	Configuration map[string]*Configuration	`yaml:"configurations"`
	StartTime int64							`yaml:"startTime"`
	FinishTime int64						`yaml:"finishTime"`
	AggregatedData *AggregatedData			`yaml:"aggregatedData"`
}

type TestInformation struct{
	Name string								`yaml:"name"`
	VersionCode string						`yaml:"version"`
	AutoconfiguringApproach string 			`yaml:"autoConfigApproach"`
	Iterations []*IterationInformation 		`yaml:"iterations"`
	InputWorkload *workload.Workload 		`yaml:"workload"`
}