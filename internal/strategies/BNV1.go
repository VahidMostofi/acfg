package strategies

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/sla"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type BNV1 struct {
	endpoints     []string
	resources     []string
	initialCPU    int64 // 1 CPU would be 1000
	initialMemory int64 // 1 Gigabyte memory would be 1024
	delta         int64
	sla           *sla.SLA
}

func NewBNV1(delta int64, endpoints []string, resources []string, initialCPU, initialMemory int64) (*BNV1, error) {
	c := &BNV1{
		endpoints:     endpoints,
		resources:     resources,
		initialCPU:    initialCPU,
		initialMemory: initialMemory,
		delta:         delta,
	}

	return c, nil
}

func (bnv1 *BNV1) AddSLA(sla *sla.SLA) error {
	bnv1.sla = sla
	return nil
}

func (bnv1 *BNV1) GetName() string {
	return "BNV1"
}

func (bnv1 *BNV1) GetInitialConfiguration(workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, error) {
	config := make(map[string]*configuration.Configuration)
	for _, resource := range bnv1.resources {
		config[resource] = &configuration.Configuration{}
		config[resource].ReplicaCount = int64Ptr(1)
		config[resource].CPU = int64Ptr(bnv1.initialCPU)
		config[resource].Memory = int64Ptr(bnv1.initialMemory)
		config[resource].ResourceType = configuration.ResourceTypeDeployment
	}

	return config, nil
}

func (bnv1 *BNV1) ConfigureNextStep(currentConfig map[string]*configuration.Configuration, workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, bool, error) {
	isChanged := false
	newConfig := make(map[string]*configuration.Configuration)

	initialCPUCount := make(map[string]int64)
	finalCPUCount := make(map[string]int64)

	for _, resource := range bnv1.resources {
		newConfig[resource] = currentConfig[resource].DeepCopy()
		initialCPUCount[resource] = *currentConfig[resource].CPU * *currentConfig[resource].ReplicaCount
		finalCPUCount[resource] = initialCPUCount[resource]
	}

	for _, endpoint := range bnv1.endpoints {

		// does this endpoint meet all the conditions?
		doMeet := true
		for _, condition := range getConditionsMatchingEndpoint(endpoint, bnv1.sla.Conditions) {
			if strings.Compare(condition.Type, sla.SLAConditionTypeResponseTime) == 0 {
				requiredValue := condition.GetComputeFunction()(*aggData.ResponseTimes[endpoint])
				if requiredValue > condition.Threshold {
					doMeet = false
					log.Debugf("BNV1.ConfigureNextStep(): endpoint %s does'nt meet the SLA. %f > %f", endpoint, requiredValue, condition.Threshold)
					break
				}
			}
		}

		// doesn't meet the sla so we add bnv1.delta CPU units to the bottle neck. (the resource in the path the most CPU utilization).
		if !doMeet {
			_, resourceWithMaxCPUUtil := aggData.GetMinMaxResourcesBasedOnCPUUtil(endpoint)
			log.Infof("BNV1.ConfigureNextStep(): for endpoint %s, %s has the most CPU utillization.", endpoint, resourceWithMaxCPUUtil)

			increaseValue := bnv1.delta

			newCPUValue := finalCPUCount[resourceWithMaxCPUUtil] + increaseValue
			log.Infof("BNV1.ConfigureNextStep(): adding %d CPU units to %s changing it from %d to %d", bnv1.delta, resourceWithMaxCPUUtil, finalCPUCount[resourceWithMaxCPUUtil], newCPUValue)
			finalCPUCount[resourceWithMaxCPUUtil] = newCPUValue

			// there is a change to we mare isChanged as true.
			isChanged = true
		}
	}

	if isChanged {
		for resourceName := range newConfig {
			newConfig[resourceName].UpdateEqualWithNewCPUValue(finalCPUCount[resourceName], 1000)
			log.Infof("BNV1.ConfigureNextStep(): new configuration for %s: %s", resourceName, newConfig[resourceName].String())
		}
	}

	return newConfig, isChanged, nil
}

func getConditionsMatchingEndpoint(endpoint string, conditions []sla.Condition) []sla.Condition {
	res := make([]sla.Condition, 0)
	for _, c := range conditions {
		if strings.Compare(c.EndpointName, endpoint) == 0 {
			res = append(res, c)
		}
	}
	return res
}
