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
	"strconv"
)

type Configuration struct{
	ResourceType string
	ReplicaCount *int64
	CPU *int64
	Memory *int64
	EnvironmentValues map[string]string
}

func (c *Configuration) DeepCopy() *Configuration{
	c2 := &Configuration{
		ResourceType: c.ResourceType,
		ReplicaCount: c.ReplicaCount,
		CPU: c.CPU,
		Memory: c.Memory,
		EnvironmentValues: make(map[string]string),
	}
	for key,value := range c.EnvironmentValues{
		c2.EnvironmentValues[key] = value
	}

	return c
}

func (c *Configuration) GetCPUStringForK8s() string{
	s := strconv.FormatInt(*c.CPU, 10) + "m"
	return s
}

func (c *Configuration) GetMemoryStringForK8s() string{
	s := strconv.FormatInt(*c.Memory, 10) + "Mi"
	return s
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

type SLA struct{
	Conditions []Condition
}

type Condition struct{
	Type string
	EndpointName string
	Threshold float64
	ComputeFn func([]float64) float64
}