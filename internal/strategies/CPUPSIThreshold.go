package strategies

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/sla"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type CPUPSIThreshold struct {
	endpoints            []string
	resources            []string
	initialCPU           int64 // 1 CPU would be 1000
	initialMemory        int64 // 1 Gigabyte memory would be 1024
	utilizationThreshold float64
	utilizationIndicator string // mean
}

func NewCPUPSIThreshold(utilizationIndicator string, utilizationThreshold float64, endpoints []string, resources []string, initialCPU, initialMemory int64) (*CPUPSIThreshold, error) {
	c := &CPUPSIThreshold{
		endpoints:            endpoints,
		resources:            resources,
		initialCPU:           initialCPU,
		initialMemory:        initialMemory,
		utilizationThreshold: utilizationThreshold,
		utilizationIndicator: utilizationIndicator,
	}

	return c, nil
}

func (ct *CPUPSIThreshold) AddSLA(sla *sla.SLA) error {
	return nil
}

func (ct *CPUPSIThreshold) GetName() string {
	return "CPUPSIThreshold"
}

func (ct *CPUPSIThreshold) GetInitialConfiguration(workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, error) {
	config := make(map[string]*configuration.Configuration)
	for _, resource := range ct.resources {
		config[resource] = &configuration.Configuration{}
		config[resource].ReplicaCount = int64Ptr(1)
		config[resource].CPU = int64Ptr(ct.initialCPU)
		config[resource].Memory = int64Ptr(ct.initialMemory)
		config[resource].ResourceType = "Deployment"
		log.Infof("%s.GetInitialConfiguration(): initial config for %s: %v", ct.GetName(), resource, config[resource])
	}
	return config, nil
}

func (ct *CPUPSIThreshold) ConfigureNextStep(currentConfig map[string]*configuration.Configuration, workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, map[string]interface{}, bool, error) {
	isChanged := false
	//var err error
	newConfig := make(map[string]*configuration.Configuration)

	for _, resource := range ct.resources {
		newConfig[resource] = currentConfig[resource].DeepCopy()
		//isChanged = true

		var whatToCompare float64
		var err error
		if ct.utilizationIndicator == "mean" {
			// TODO
			whatToCompare, err = aggData.MemPsiUtilizations[resource].GetMean()
			if err != nil {
				return nil, make(map[string]interface{}), false, errors.Wrapf(err, "error while computing mean of mem psi utilizations for %s.", resource)
			}
		}
		if whatToCompare > ct.utilizationThreshold {
			//newCount := int64Ptr(*newConfig[resource].ReplicaCount + 1)
			log.Infof("%s.ConfigureNextStep() mem psi utilization for %s is %f is more than %f changing replica allocation... from %d to %d", ct.GetName(), resource, whatToCompare, ct.utilizationThreshold, *newConfig[resource].ReplicaCount, 0)
			//newConfig[resource].ReplicaCount = newCount
			//newConfig[resource].CPU = int64Ptr(1000)
			newConfig[resource].Memory = int64Ptr(*newConfig[resource].Memory + (1 * 1024))
			isChanged = true
		} else {
			log.Infof("%s.ConfigureNextStep() CPU utilization for %s is %f is less than %f not changing replica from %d", ct.GetName(), resource, whatToCompare, ct.utilizationThreshold, *newConfig[resource].ReplicaCount)
			//log.Infof("%s.ConfigureNextStep() CPU utilization for %s is %f is less than %f not changing replica from %d", ct.GetName(), resource, whatToCompare, ct.utilizationThreshold, *newConfig[resource].ReplicaCount)
		}
	}

	return newConfig, make(map[string]interface{}), isChanged, nil
}
