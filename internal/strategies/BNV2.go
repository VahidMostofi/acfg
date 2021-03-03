package strategies

import (
	"crypto/md5"
	"encoding/json"
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

type BNV2 struct {
	endpoints        []string
	resources        []string
	initialCPU       int64 // 1 CPU would be 1000
	initialMemory    int64 // 1 Gigabyte memory would be 1024
	initialDelta     int64
	sla              *sla.SLA
	minimumDelta     int64
	minimumCPUValue  int64
	maxCPUPerReplica int64

	resource2Delta   map[string]int64
	cache            map[string]int64
	bestWhichMeets   int64
	iterationCount   int
	stage1Iterations int
	stage2Iterations int
	stage2IsStarted  bool
}

func NewBNV2(initialDelta int64, endpoints []string, resources []string, initialCPU, initialMemory, maxCPUPerReplica, minimumCPUValue, minimumDelta int64) (*BNV2, error) {
	log.Debugf("NewBNV2: initalDelta: %d, initialCPU: %d, initialMemory: %d, maxCPUPerReplica: %d, minimumCPUValue: %d, minimumDelta: %d", initialDelta, initialCPU, initialMemory, maxCPUPerReplica, minimumCPUValue, minimumDelta)
	c := &BNV2{
		endpoints:        endpoints,
		resources:        resources,
		initialCPU:       initialCPU,
		initialMemory:    initialMemory,
		initialDelta:     initialDelta,
		maxCPUPerReplica: maxCPUPerReplica,
		minimumDelta:     minimumDelta,
		minimumCPUValue:  minimumCPUValue,
		resource2Delta:   make(map[string]int64),
		cache:            make(map[string]int64),
		stage2IsStarted:  false,
	}

	for _, key := range c.resources {
		c.resource2Delta[key] = c.initialDelta
	}

	return c, nil
}

func round1(value float64) float64 {
	return math.Round(value*10) / 10
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func (bnv2 *BNV2) AddSLA(sla *sla.SLA) error {
	bnv2.sla = sla
	return nil
}

func (bnv2 *BNV2) GetName() string {
	return "BNV2"
}

func (bnv2 *BNV2) GetInitialConfiguration(workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, error) {
	config := make(map[string]*configuration.Configuration)
	for _, resource := range bnv2.resources {
		config[resource] = &configuration.Configuration{}
		config[resource].ReplicaCount = int64Ptr(1)
		config[resource].CPU = int64Ptr(bnv2.initialCPU)
		config[resource].Memory = int64Ptr(bnv2.initialMemory)
		config[resource].ResourceType = configuration.ResourceTypeDeployment

		config[resource].UpdateEqualWithNewCPUValue(*config[resource].CPU, bnv2.maxCPUPerReplica)
	}

	return config, nil
}

func (bnv2 *BNV2) ConfigureNextStep(currentConfig map[string]*configuration.Configuration, workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, bool, error) {
	isChanged := false
	bnv2.iterationCount++
	newConfig := make(map[string]*configuration.Configuration)

	initialCPUCount := make(map[string]int64)
	finalCPUCount := make(map[string]int64)

	for _, resource := range bnv2.resources {
		newConfig[resource] = currentConfig[resource].DeepCopy()
		initialCPUCount[resource] = *currentConfig[resource].CPU * *currentConfig[resource].ReplicaCount
		finalCPUCount[resource] = initialCPUCount[resource]
	}

	allMeet := true
	for _, endpoint := range bnv2.endpoints {
		doMeet := true
		for _, condition := range getConditionsMatchingEndpoint(endpoint, bnv2.sla.Conditions) {
			log.Debugf("BNV2.ConfigureNextStep(): condition %v for endpoint %s", condition, endpoint)
			if strings.Compare(condition.Type, sla.SLAConditionTypeResponseTime) == 0 {
				requiredValue := condition.GetComputeFunction()(*aggData.ResponseTimes[endpoint])
				if requiredValue > condition.Threshold {
					doMeet = false
					log.Debugf("BNV2.ConfigureNextStep(): endpoint %s does'nt meet the SLA. %f > %f", endpoint, requiredValue, condition.Threshold)
					break
				}
			}
		}
		if !doMeet { // this endpoint is not meeting the SLA
			allMeet = false
			_, resourceWithMaxCPUUtil := aggData.GetMinMaxResourcesBasedOnCPUUtil(endpoint)
			log.Infof("BNV2.ConfigureNextStep(): for endpoint %s, %s has the most CPU utillization.", endpoint, resourceWithMaxCPUUtil)

			increaseValue := bnv2.resource2Delta[resourceWithMaxCPUUtil]
			if increaseValue > 0 {
				newCPUValue := finalCPUCount[resourceWithMaxCPUUtil] + increaseValue
				log.Infof("BNV2.ConfigureNextStep(): adding %d CPU units to %s changing it from %d to %d", increaseValue, resourceWithMaxCPUUtil, finalCPUCount[resourceWithMaxCPUUtil], newCPUValue)
				finalCPUCount[resourceWithMaxCPUUtil] = newCPUValue
				isChanged = true
			} else {
				return nil, false, fmt.Errorf("BNV2.ConfigureNextStep(): the increaseValue must be positive. Something is wrong.")
			}
		}
	}
	if allMeet || bnv2.stage2IsStarted {
		log.Infof("BNV2.ConfigureNextStep(): all endpoints are meeting the required SLAs, starting second stage.")
		bnv2.stage2IsStarted = true

		marginalRequests := make(map[string]float64)
		resourcesWhichAreMax := make(map[string]float64) // will be used as a set
		resource2ednpoints := aggData.SystemStructure.GetResources2Endpoints()
		// resource2EndpointWithMaxResponseTime := make(map[string]float64)

		for _, endpoint := range bnv2.endpoints {
			for _, condition := range getConditionsMatchingEndpoint(endpoint, bnv2.sla.Conditions) {
				if strings.Compare(condition.Type, sla.SLAConditionTypeResponseTime) == 0 {
					actualValue := condition.GetComputeFunction()(*aggData.ResponseTimes[endpoint])
					if actualValue > 0.9*condition.Threshold {
						marginalRequests[endpoint] = actualValue
					}
				}
			}
			_, resourceWithMaxCPUUtil := aggData.GetMinMaxResourcesBasedOnCPUUtil(endpoint)
			resourcesWhichAreMax[resourceWithMaxCPUUtil] = 0 // is being used as a set
		}
		backwardCandidates := make([]string, 0)
		for _, resouce := range bnv2.resources {
			flag := true

			if _, exists := resourcesWhichAreMax[resouce]; exists {
				flag = false
			}

			if finalCPUCount[resouce]-bnv2.minimumCPUValue < 5 {
				flag = false
			}

			for criticalRequest := range marginalRequests {
				if funk.ContainsString(resource2ednpoints[resouce], criticalRequest) {
					flag = false
				}
			}

			if flag {
				backwardCandidates = append(backwardCandidates, resouce)
			}
		}

		log.Debugf("BNV2.ConfigureNextStep(): backward candidates: %v", backwardCandidates)
		if len(backwardCandidates) == 0 {
			// fmt.Println("no candidate for moving backward")
		} else {
			sort.Slice(backwardCandidates, func(i int, j int) bool {
				return bnv2.resource2Delta[backwardCandidates[i]] > bnv2.resource2Delta[backwardCandidates[j]]
			})
			pruneCount := int(len(finalCPUCount) / 3)
			if pruneCount <= 0 {
				pruneCount = 1
			}
			for cIdx := 0; cIdx < pruneCount; cIdx++ {
				if cIdx == len(backwardCandidates) {
					break
				}

				serviceToDecrease := backwardCandidates[cIdx]
				// fmt.Println(minCPUUtil)
				log.Println("BNV2.ConfigureNextStep():", serviceToDecrease, "is going to be pruned.")
				bnv2.resource2Delta[serviceToDecrease] /= 2
				bnv2.resource2Delta[serviceToDecrease] = int64(math.Max(float64(bnv2.resource2Delta[serviceToDecrease]), float64(bnv2.minimumCPUValue)))
				decreaseValue := bnv2.resource2Delta[serviceToDecrease]
				if decreaseValue > 0 {
					// log.Println("BNV2.ConfigureNextStep():", serviceToDecrease, "is part of", requestName, "stepSize for path(request)", requestName, "is", bnv2.resource2Delta[serviceToDecrease])
					prev := finalCPUCount[serviceToDecrease]
					finalCPUCount[serviceToDecrease] -= decreaseValue
					finalCPUCount[serviceToDecrease] = int64(math.Max(float64(bnv2.minimumCPUValue), float64(finalCPUCount[serviceToDecrease])))
					log.Println("BNV2.ConfigureNextStep():", serviceToDecrease, "updating total CPU count from", prev, "to", finalCPUCount[serviceToDecrease])
					isChanged = true
				} else {
					return nil, false, fmt.Errorf("BNV2.ConfigureNextStep(): the increaseValue must be positive")
				}
			}
		}
	} // end allMeet || bnv2.stage2IsStarted

	var totalPrev int64 = 0
	var totalNew int64 = 0
	for service := range initialCPUCount {
		finalCPUCount[service] = minInt64(finalCPUCount[service], initialCPUCount[service]+bnv2.initialDelta)
		finalCPUCount[service] = maxInt64(finalCPUCount[service], initialCPUCount[service]-bnv2.initialDelta/2)
		totalNew += finalCPUCount[service]
		totalPrev += initialCPUCount[service]
	}

	if !bnv2.stage2IsStarted {
		if bnv2.stage1Iterations == 0 { //need to know how many steps in stage one
			bnv2.stage1Iterations = bnv2.iterationCount
		}
	} else { // we are in stage 2

		log.Debugf("BNV2.ConfigureNextStep(): iterations (1: %d, 2: %d)", bnv2.stage1Iterations, bnv2.stage2Iterations)
		if bnv2.stage2Iterations >= bnv2.stage1Iterations+1 { // +1 because the first iteration is not counted.
			return nil, false, nil
		}
		bnv2.stage2Iterations++ // tracking iterations in stage 2
	}

	if totalPrev < bnv2.bestWhichMeets {
		bnv2.bestWhichMeets = totalPrev
	}

	if isChanged {
		for resourceName := range newConfig {
			newConfig[resourceName] = currentConfig[resourceName].DeepCopy()
			newConfig[resourceName].UpdateEqualWithNewCPUValue(finalCPUCount[resourceName], bnv2.maxCPUPerReplica)
			log.Infof("BNV2.ConfigureNextStep(): new configuration for %s: %s", resourceName, newConfig[resourceName].String())
		}
	}
	h := getConfigHash(newConfig)
	if _, exists := bnv2.cache[h]; exists {
		return nil, false, nil
	} else {
		bnv2.cache[h] = totalNew
	}

	return newConfig, isChanged, nil
}

func getConfigHash(in map[string]*configuration.Configuration) string {
	b, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	h := md5.Sum(b)
	return fmt.Sprintf("%x", h)
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
