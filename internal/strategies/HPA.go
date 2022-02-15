package strategies

import (
	"math"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/sla"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type HPA struct {
	counter       int
	threshold     float64
	resources     []string
	initialCPU    int64 // 1 CPU would be 1000
	initialMemory int64 // 1 Gigabyte memory would be 1024
}

func NewHPA(resources []string, threshold float64, initialCPU, initialMemory int64) (*HPA, error) {
	c := &HPA{threshold: threshold, resources: resources, initialCPU: initialCPU, initialMemory: initialMemory}
	c.counter = 1
	return c, nil
}

func (HPA *HPA) AddSLA(sla *sla.SLA) error {
	return nil
}

func (hpa *HPA) GetName() string {
	return "HPA-" + strconv.Itoa(int(hpa.threshold))
}

func (hpa *HPA) GetInitialConfiguration(workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, error) {
	config := make(map[string]*configuration.Configuration)
	for _, resource := range hpa.resources {
		config[resource] = &configuration.Configuration{}
		config[resource].ReplicaCount = int64Ptr(1)
		config[resource].CPU = int64Ptr(hpa.initialCPU)
		config[resource].Memory = int64Ptr(hpa.initialMemory)
		config[resource].ResourceType = configuration.ResourceTypeDeployment
	}

	return config, nil
}

func (hpa *HPA) ConfigureNextStep(currentConfig map[string]*configuration.Configuration, workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, map[string]interface{}, bool, error) {
	isChanged := false
	extraInfo := make(map[string]interface{})
	newConfig := make(map[string]*configuration.Configuration)
	finalReplicas := make(map[string]int)

	for _, resource := range hpa.resources {
		newConfig[resource] = currentConfig[resource].DeepCopy()

		currentReplicas := aggData.DeploymentInfos[resource].Replica

		meanCPUUtilization, err := aggData.CPUUtilizations[resource].GetMean()
		if err != nil {
			log.Info("lastReplicas are being use due an error in getting mean cpu utilzation for " + resource + ".")
		} else {
			finalReplicas[resource] = int(math.Ceil(float64(currentReplicas) * (meanCPUUtilization / hpa.threshold)))
			newConfig[resource].ReplicaCount = int64Ptr(int64(finalReplicas[resource]))
			if currentReplicas != finalReplicas[resource] {
				isChanged = true
			}
		}
	}
	hpa.counter++
	if hpa.counter >= 3 {
		return newConfig, extraInfo, false, nil
	}

	if isChanged {
		for resourceName := range newConfig {
			log.Infof("HPA.ConfigureNextStep(): new configuration for %s: %s", resourceName, newConfig[resourceName].String())
		}
	}

	return newConfig, extraInfo, isChanged, nil
}
