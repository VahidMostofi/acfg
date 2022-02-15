package strategies

import (
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/sla"
	"github.com/vahidmostofi/acfg/internal/workload"
	"math"
	"strings"
)

type BNV1Mem struct {
	endpoints     []string
	resources     []string
	initialCPU    int64 // 1 CPU would be 1000
	initialMemory int64 // 1 Gigabyte memory would be 1024
	cpuDelta      int64
	memDelta      int64
	sla           *sla.SLA
}

func NewBNV1Mem(cpuDelta int64, memDelta int64, endpoints []string, resources []string, initialCPU, initialMemory int64) (*BNV1Mem, error) {
	c := &BNV1Mem{
		endpoints:     endpoints,
		resources:     resources,
		initialCPU:    initialCPU,
		initialMemory: initialMemory,
		cpuDelta:      cpuDelta,
		memDelta:      memDelta,
	}

	return c, nil
}

func (bnv1 *BNV1Mem) AddSLA(sla *sla.SLA) error {
	bnv1.sla = sla
	return nil
}

func (bnv1 *BNV1Mem) GetName() string {
	return "BNV1Mem"
}

func (bnv1 *BNV1Mem) GetInitialConfiguration(workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, error) {
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

func (bnv1 *BNV1Mem) ConfigureNextStep(currentConfig map[string]*configuration.Configuration, workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, map[string]interface{}, bool, error) {
	isChanged := false
	extraInfo := make(map[string]interface{})
	newConfig := make(map[string]*configuration.Configuration)

	initialCPUCount := make(map[string]int64)
	initialMemCount := make(map[string]int64)
	finalCPUCount := make(map[string]int64)
	finalMemCount := make(map[string]int64)

	for _, resource := range bnv1.resources {
		newConfig[resource] = currentConfig[resource].DeepCopy()
		initialCPUCount[resource] = *currentConfig[resource].CPU * *currentConfig[resource].ReplicaCount
		initialMemCount[resource] = *currentConfig[resource].Memory * *currentConfig[resource].ReplicaCount
		finalCPUCount[resource] = initialCPUCount[resource]
		finalMemCount[resource] = initialMemCount[resource]
	}

	for _, endpoint := range bnv1.endpoints {

		// does this endpoint meet all the conditions?
		doMeet := true
		for _, condition := range getConditionsMatchingEndpoint(endpoint, bnv1.sla.Conditions) {
			if strings.Compare(condition.Type, sla.SLAConditionTypeResponseTime) == 0 {
				requiredValue := condition.GetComputeFunction()(*aggData.ResponseTimes[endpoint])
				if requiredValue > condition.Threshold {
					doMeet = false
					log.Debugf("BNV1Mem.ConfigureNextStep(): endpoint %s does'nt meet the SLA. %f > %f", endpoint, requiredValue, condition.Threshold)
					break
				}
				if len(*aggData.ResponseTimes[endpoint]) == 0 {
					doMeet = false
					log.Debugf("BNV1Mem.ConfigureNextStep(): endpoint %s does'nt meet the SLA. There are no recorded response times.", endpoint)
					break
				}
			}
		}

		// doesn't meet the sla so we add cpu or mem shares to the bottleneck
		if !doMeet {

			_, resourceWithMaxPSI := aggData.GetMinMaxResourcesBasedOnPsi(endpoint)

			log.Infof("BNV1Mem.ConfigureNextStep(): for endpoint %s, %v", endpoint, resourceWithMaxPSI)
			// TODO there is a bug here where each endpoint SLA violation can add to the new resource allocation
			if resourceWithMaxPSI.ResourceType == aggregators.CPU {
				increaseValue := bnv1.cpuDelta
				newCPUValue := finalCPUCount[resourceWithMaxPSI.ResourceName] + increaseValue
				log.Infof("BNV1Mem.ConfigureNextStep(): adding %d CPU units to %s changing it from %d to %d", bnv1.cpuDelta, resourceWithMaxPSI.ResourceName, finalCPUCount[resourceWithMaxPSI.ResourceName], newCPUValue)
				// Dont downgrade
				finalCPUCount[resourceWithMaxPSI.ResourceName] = int64(math.Max(float64(newCPUValue), float64(finalCPUCount[resourceWithMaxPSI.ResourceName])))
				isChanged = true
			} else if resourceWithMaxPSI.ResourceType == aggregators.Mem {
				increaseValue := bnv1.memDelta
				newMemValue := finalMemCount[resourceWithMaxPSI.ResourceName] + increaseValue
				log.Infof("BNV1Mem.ConfigureNextStep(): adding %d Mem units to %s changing it from %d to %d", bnv1.memDelta, resourceWithMaxPSI.ResourceName, finalMemCount[resourceWithMaxPSI.ResourceName], newMemValue)
				// Dont downgrade
				finalMemCount[resourceWithMaxPSI.ResourceName] = int64(math.Max(float64(newMemValue), float64(finalMemCount[resourceWithMaxPSI.ResourceName])))
				isChanged = true
			}
		}
	}

	if isChanged {
		for resourceName := range newConfig {
			// TODO add max to cli args
			newConfig[resourceName].UpdateEqualWithNewCpuMemValue(finalCPUCount[resourceName], 100000, finalMemCount[resourceName], 100000)
			log.Infof("BNV1Mem.ConfigureNextStep(): new configuration for %s: %s", resourceName, newConfig[resourceName].String())
		}
	}

	return newConfig, extraInfo, isChanged, nil
}
