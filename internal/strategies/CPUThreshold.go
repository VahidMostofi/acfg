package strategies

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type CPUThreshold struct{
	endpoints []string
	resources []string
	initialCPU int64 // 1 CPU would be 1000
	initialMemory int64 // 1 Gigabyte memory would be 1024
	utilizationThreshold float64
	utilizationIndicator string // mean
}

func NewCPUThreshold(utilizationIndicator string, utilizationThreshold float64, endpoints []string, resources []string, initialCPU, initialMemory int64) (*CPUThreshold, error){
	c := &CPUThreshold{
		endpoints: endpoints,
		resources: resources,
		initialCPU: initialCPU,
		initialMemory: initialMemory,
		utilizationThreshold: utilizationThreshold,
		utilizationIndicator: utilizationIndicator,
	}

	return c, nil
}

func (ct *CPUThreshold) GetName() string{
	return "CPUThreshold"
}

func (ct *CPUThreshold) GetInitialConfiguration(workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, error){
	config := make(map[string]*configuration.Configuration)
	for _,resource := range ct.resources{
		config[resource] = &configuration.Configuration{}
		config[resource].ReplicaCount = int64Ptr(1)
		config[resource].CPU = int64Ptr(1000)
		config[resource].Memory = int64Ptr(512)
		config[resource].ResourceType = "Deployment"
	}

	return config, nil
}

func (ct *CPUThreshold) ConfigureNextStep(currentConfig map[string]*configuration.Configuration, workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, bool, error){
	isChanged := false
	var err error
	newConfig := make(map[string]*configuration.Configuration)


	for _, resource := range ct.resources{
		newConfig[resource] = currentConfig[resource].DeepCopy()

		var whatToCompare float64
		if ct.utilizationIndicator == "mean"{
			whatToCompare, err = aggData.CPUUtilizations[resource].GetMean()
			if err != nil{
				return nil,false, errors.Wrapf(err, "error while computing mean of CPU utilizations for %s.", resource)
			}
		}
		if whatToCompare > ct.utilizationThreshold{
			newCount := int64Ptr(*newConfig[resource].ReplicaCount + 1)
			log.Infof(ct.GetName(), "ConfigureNextStep() CPU utilization for %s is %f is more than %f changing replica from %d to %d", resource, whatToCompare, ct.utilizationThreshold, *newConfig[resource].ReplicaCount, *newCount)
			newConfig[resource].ReplicaCount = newCount
			isChanged = true
		}
	}

	return newConfig, isChanged, nil
}