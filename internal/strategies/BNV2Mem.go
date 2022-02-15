package strategies

import (
	"fmt"
	"math"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/sla"
	"github.com/vahidmostofi/acfg/internal/workload"

	"github.com/thoas/go-funk"
)

type BNV2Mem struct {
	endpoints        []string
	resources        []string
	sla              *sla.SLA
	initialCPU       int64 // 1 CPU would be 1000
	initialMemory    int64 // 1 Gigabyte memory would be 1024
	maxCPUPerReplica int64
	maxMemPerReplica int64
	initialCPUDelta  int64
	initialMemDelta  int64
	minimumCPUValue  int64
	minimumMemValue  int64
	minimumCpuDelta  int64
	minimumMemDelta  int64

	cpuResource2Delta map[string]int64
	memResource2Delta map[string]int64
	cache             map[string]int64
	bestWhichMeets    int64
	iterationCount    int
	stage1Iterations  int
	stage2Iterations  int
	stage2IsStarted   bool
}

func NewBNV2Mem(endpoints []string, resources []string, initialCPU int64, initialMemory int64, maxCPUPerReplica int64, maxMemPerReplica int64, initialCPUDelta int64, initialMemDelta int64, minimumCPUValue int64, minimumMemValue int64, minimumCpuDelta int64, minimumMemDelta int64) (*BNV2Mem, error) {
	//log.Debugf("NewBNV2Mem: initalDelta: %d, initialCPU: %d, initialMemory: %d, maxCPUPerReplica: %d, minimumCPUValue: %d, minimumDelta: %d", initialDelta, initialCPU, initialMemory, maxCPUPerReplica, minimumCPUValue, minimumDelta)
	c := &BNV2Mem{
		endpoints:         endpoints,
		resources:         resources,
		initialCPU:        initialCPU,
		initialMemory:     initialMemory,
		maxCPUPerReplica:  maxCPUPerReplica,
		maxMemPerReplica:  maxMemPerReplica,
		initialCPUDelta:   initialCPUDelta,
		initialMemDelta:   initialMemDelta,
		minimumCPUValue:   minimumCPUValue,
		minimumMemValue:   minimumMemValue,
		minimumCpuDelta:   minimumCpuDelta,
		minimumMemDelta:   minimumMemDelta,
		cpuResource2Delta: make(map[string]int64),
		memResource2Delta: make(map[string]int64),
		cache:             make(map[string]int64),
		stage2IsStarted:   false,
	}

	for _, key := range c.resources {
		c.cpuResource2Delta[key] = c.initialCPUDelta
		c.memResource2Delta[key] = c.initialMemDelta
	}

	return c, nil
}

func (bnv2 *BNV2Mem) AddSLA(sla *sla.SLA) error {
	bnv2.sla = sla
	return nil
}

func (bnv2 *BNV2Mem) GetName() string {
	return "BNV2Mem"
}

func (bnv2 *BNV2Mem) GetInitialConfiguration(workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, error) {
	config := make(map[string]*configuration.Configuration)
	for _, resource := range bnv2.resources {
		config[resource] = &configuration.Configuration{}
		config[resource].ReplicaCount = int64Ptr(1)
		config[resource].CPU = int64Ptr(bnv2.initialCPU)
		config[resource].Memory = int64Ptr(bnv2.initialMemory)
		config[resource].ResourceType = configuration.ResourceTypeDeployment

		config[resource].UpdateEqualWithNewCpuMemValue(*config[resource].CPU, bnv2.maxCPUPerReplica, *config[resource].Memory, bnv2.maxMemPerReplica)
	}

	return config, nil
}

func (bnv2 *BNV2Mem) ConfigureNextStep(currentConfig map[string]*configuration.Configuration, workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, map[string]interface{}, bool, error) {
	isChanged := false
	extraInfo := make(map[string]interface{})
	bnv2.iterationCount++
	newConfig := make(map[string]*configuration.Configuration)

	initialCPUCount := make(map[string]int64)
	initialMemCount := make(map[string]int64)
	finalCPUCount := make(map[string]int64)
	finalMemCount := make(map[string]int64)

	for _, resource := range bnv2.resources {
		newConfig[resource] = currentConfig[resource].DeepCopy()
		initialCPUCount[resource] = *currentConfig[resource].CPU * *currentConfig[resource].ReplicaCount
		initialMemCount[resource] = *currentConfig[resource].Memory * *currentConfig[resource].ReplicaCount
		finalCPUCount[resource] = initialCPUCount[resource]
		finalMemCount[resource] = initialMemCount[resource]
	}

	allMeet := true
	for _, endpoint := range bnv2.endpoints {
		doMeet := true
		for _, condition := range getConditionsMatchingEndpoint(endpoint, bnv2.sla.Conditions) {
			log.Debugf("BNV2Mem.ConfigureNextStep(): condition %v for endpoint %s", condition, endpoint)
			if strings.Compare(condition.Type, sla.SLAConditionTypeResponseTime) == 0 {
				requiredValue := condition.GetComputeFunction()(*aggData.ResponseTimes[endpoint])
				if requiredValue > condition.Threshold*Safety {
					doMeet = false
					log.Debugf("BNV2Mem.ConfigureNextStep(): endpoint %s does'nt meet the SLA. %f > %f (%f)", endpoint, requiredValue, condition.Threshold, condition.Threshold*Safety)
					break
				}
				if len(*aggData.ResponseTimes[endpoint]) == 0 {
					doMeet = false
					log.Debugf("BNV2Mem.ConfigureNextStep(): endpoint %s does'nt meet the SLA. There are no recorded response times.", endpoint)
					break
				}
			}
		}

		if !doMeet { // this endpoint is not meeting the SLA
			allMeet = false
			_, resourceWithMaxPSI := aggData.GetMinMaxResourcesBasedOnPsi(endpoint)

			log.Infof("BNV2Mem.ConfigureNextStep(): for endpoint %s, %v", endpoint, resourceWithMaxPSI)

			increaseCpuValue := bnv2.cpuResource2Delta[resourceWithMaxPSI.ResourceName]
			increaseMemValue := bnv2.memResource2Delta[resourceWithMaxPSI.ResourceName]

			if resourceWithMaxPSI.ResourceType == aggregators.CPU {
				if increaseCpuValue > 0 {
					newCPUValue := finalCPUCount[resourceWithMaxPSI.ResourceName] + increaseCpuValue
					log.Infof("BNV2Mem.ConfigureNextStep(): adding %d CPU units to %s changing it from %d to %d", increaseCpuValue, resourceWithMaxPSI.ResourceName, finalCPUCount[resourceWithMaxPSI.ResourceName], newCPUValue)
					// Don't downgrade here. This is to resolve an endpoint not meeting an SLA
					finalCPUCount[resourceWithMaxPSI.ResourceName] = int64(math.Max(float64(newCPUValue), float64(finalCPUCount[resourceWithMaxPSI.ResourceName])))
					isChanged = true
				} else {
					return nil, extraInfo, false, fmt.Errorf("BNV2Mem.ConfigureNextStep(): the increaseCpuValue must be positive. Something is wrong")
				}

			} else if resourceWithMaxPSI.ResourceType == aggregators.Mem {
				if increaseMemValue > 0 {
					newMemValue := finalMemCount[resourceWithMaxPSI.ResourceName] + increaseMemValue
					log.Infof("BNV2Mem.ConfigureNextStep(): adding %d Mem units to %s changing it from %d to %d", increaseMemValue, resourceWithMaxPSI.ResourceName, finalMemCount[resourceWithMaxPSI.ResourceName], newMemValue)
					// Don't downgrade here. This is to resolve an endpoint not meeting an SLA
					finalMemCount[resourceWithMaxPSI.ResourceName] = int64(math.Max(float64(newMemValue), float64(finalMemCount[resourceWithMaxPSI.ResourceName])))
					isChanged = true
				} else {
					return nil, extraInfo, false, fmt.Errorf("BNV2Mem.ConfigureNextStep(): the increaseMemValue must be positive. Something is wrong")
				}
			}
		}
	}
	extraInfo["doMeet"] = allMeet
	if allMeet || bnv2.stage2IsStarted {
		log.Infof("BNV2Mem.ConfigureNextStep(): all endpoints are meeting the required SLAs, starting second stage.")
		bnv2.stage2IsStarted = true

		marginalRequests := make(map[string]float64)
		resourcesWhichAreMax := make(map[string]float64) // will be used as a set
		resource2endpoints := aggData.SystemStructure.GetResources2Endpoints()

		for _, endpoint := range bnv2.endpoints {
			for _, condition := range getConditionsMatchingEndpoint(endpoint, bnv2.sla.Conditions) {
				if strings.Compare(condition.Type, sla.SLAConditionTypeResponseTime) == 0 {
					actualValue := condition.GetComputeFunction()(*aggData.ResponseTimes[endpoint])
					if actualValue > 0.9*condition.Threshold {
						marginalRequests[endpoint] = actualValue
					}
				}
			}
			_, resourceWithMaxPsi := aggData.GetMinMaxResourcesBasedOnPsi(endpoint)

			resourcesWhichAreMax[resourceWithMaxPsi.ResourceName] = 0 // is being used as a set
		}
		backwardCandidates := make([]string, 0)
		for _, resource := range bnv2.resources {
			flag := true

			// don't prune a resource with the highest PSI (mem or cpu)
			if _, exists := resourcesWhichAreMax[resource]; exists {
				flag = false
			}

			// don't prune a resource if the resource cpu is at the min value
			if finalCPUCount[resource]-bnv2.minimumCPUValue < 5 {
				flag = false
			}
			// don't prune a resource if the resource mem is at the min value
			if finalMemCount[resource]-bnv2.minimumMemValue < 5 {
				flag = false
			}

			// don't prune a resource where a related endpoint SLA is close to violation
			for criticalRequest := range marginalRequests {
				if funk.ContainsString(resource2endpoints[resource], criticalRequest) {
					flag = false
				}
			}

			if flag {
				backwardCandidates = append(backwardCandidates, resource)
			}
		}

		log.Debugf("BNV2Mem.ConfigureNextStep(): backward candidates: %v", backwardCandidates)
		if len(backwardCandidates) == 0 {
			// fmt.Println("no candidate for moving backward")
		} else {
			// sort the backwards candidates by the candidate with the largest resource2Delta's
			sort.Slice(backwardCandidates, func(i int, j int) bool {
				// TODO fix this. add some real logic
				return (bnv2.cpuResource2Delta[backwardCandidates[i]] + bnv2.memResource2Delta[backwardCandidates[i]]) > (bnv2.cpuResource2Delta[backwardCandidates[j]] + bnv2.memResource2Delta[backwardCandidates[j]])
			})

			// Why?
			pruneCount := int(len(finalCPUCount) / 3)
			if pruneCount <= 0 {
				pruneCount = 1
			}

			// for every prune...
			for cIdx := 0; cIdx < pruneCount; cIdx++ {
				if cIdx == len(backwardCandidates) {
					break
				}

				serviceToDecrease := backwardCandidates[cIdx]
				// fmt.Println(minCPUUtil)
				log.Println("BNV2Mem.ConfigureNextStep():", serviceToDecrease, "is going to be pruned.")

				// Prune the resource with the smallest PSI
				cpuPSI, _ := aggData.CPUPSIUtilizations[serviceToDecrease].GetMean()
				memPSI, _ := aggData.MemPsiUtilizations[serviceToDecrease].GetMean()

				isChanged = true
				if cpuPSI > memPSI {
					//	prune mem
					err := pruneMem(bnv2, serviceToDecrease, finalMemCount)
					if err != nil {
						return nil, extraInfo, false, err
					}
				} else if cpuPSI < memPSI {
					err := pruneCPU(bnv2, serviceToDecrease, finalCPUCount)
					if err != nil {
						return nil, extraInfo, false, err
					}
				}
			}
		}
	} // end allMeet || bnv2.stage2IsStarted

	if !bnv2.stage2IsStarted {
		if bnv2.stage1Iterations == 0 { //need to know how many steps in stage one
			bnv2.stage1Iterations = bnv2.iterationCount
		}
	} else { // we are in stage 2

		log.Debugf("BNV2Mem.ConfigureNextStep(): iterations (1: %d, 2: %d)", bnv2.stage1Iterations, bnv2.stage2Iterations)
		if bnv2.stage2Iterations >= bnv2.stage1Iterations+1 { // +1 because the first iteration is not counted.
			return nil, extraInfo, false, nil
		}
		bnv2.stage2Iterations++ // tracking iterations in stage 2
	}

	if isChanged {
		for resourceName := range newConfig {
			newConfig[resourceName] = currentConfig[resourceName].DeepCopy()
			newConfig[resourceName].UpdateEqualWithNewCpuMemValue(finalCPUCount[resourceName], bnv2.maxCPUPerReplica, finalMemCount[resourceName], bnv2.maxMemPerReplica)
			log.Infof("BNV2Mem.ConfigureNextStep(): new configuration for %s: %s", resourceName, newConfig[resourceName].String())
		}
	}

	h := getConfigHash(newConfig)
	if _, exists := bnv2.cache[h]; exists {
		return nil, extraInfo, false, nil
	} else {
		bnv2.cache[h] = 0 // value doesn't matter?
	}

	return newConfig, extraInfo, isChanged, nil
}

func pruneCPU(bnv2 *BNV2Mem, serviceToDecrease string, finalCPUCount map[string]int64) error {
	//	prune cpu
	bnv2.cpuResource2Delta[serviceToDecrease] /= 2
	bnv2.cpuResource2Delta[serviceToDecrease] = int64(math.Max(float64(bnv2.cpuResource2Delta[serviceToDecrease]), float64(bnv2.minimumCpuDelta)))
	cpuDecreaseValue := bnv2.cpuResource2Delta[serviceToDecrease]
	if cpuDecreaseValue > 0 {
		prev := finalCPUCount[serviceToDecrease]
		finalCPUCount[serviceToDecrease] -= cpuDecreaseValue
		finalCPUCount[serviceToDecrease] = int64(math.Max(float64(bnv2.minimumCPUValue), float64(finalCPUCount[serviceToDecrease])))
		log.Println("BNV2Mem.ConfigureNextStep():", serviceToDecrease, "updating total CPU count from", prev, "to", finalCPUCount[serviceToDecrease])
	} else {
		return fmt.Errorf("BNV2Mem.ConfigureNextStep(): the cpuDecreaseValue must be positive")
	}
	return nil
}

func pruneMem(bnv2 *BNV2Mem, serviceToDecrease string, finalMemCount map[string]int64) error {
	bnv2.memResource2Delta[serviceToDecrease] /= 2
	bnv2.memResource2Delta[serviceToDecrease] = int64(math.Max(float64(bnv2.memResource2Delta[serviceToDecrease]), float64(bnv2.minimumMemDelta)))
	memDecreaseValue := bnv2.memResource2Delta[serviceToDecrease]
	if memDecreaseValue > 0 {
		prev := finalMemCount[serviceToDecrease]
		finalMemCount[serviceToDecrease] -= memDecreaseValue
		finalMemCount[serviceToDecrease] = int64(math.Max(float64(bnv2.minimumMemValue), float64(finalMemCount[serviceToDecrease])))
		log.Println("BNV2Mem.ConfigureNextStep():", serviceToDecrease, "updating total Mem count from", prev, "to", finalMemCount[serviceToDecrease])
	} else {
		return fmt.Errorf("BNV2Mem.ConfigureNextStep(): the memDecreaseValue must be positive")
	}
	return nil
}
