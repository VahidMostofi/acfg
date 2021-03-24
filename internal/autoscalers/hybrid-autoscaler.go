package autoscalers

import (
	"fmt"
	"math"

	"github.com/vahidmostofi/acfg/internal/aggregators"
)

func NewHybridAutoscaler(endpoints, resources []string, hpaThreshold int64) (Agent, error) {
	h := &Hybrid{
		endpoints:        endpoints,
		resources:        resources,
		hpaThreshold:     hpaThreshold,
		cooldown:         0,
		previousReplicas: make([]map[string]int, 0),
	}

	return h, nil
}

type Hybrid struct {
	hpaThreshold     int64
	cooldown         int // this is also the max length of previousReplicas
	endpoints        []string
	resources        []string
	previousReplicas []map[string]int
}

func (h *Hybrid) GetName() string {
	return "hybrid"
}

func (h *Hybrid) Evaluate(aggData *aggregators.AggregatedData) (map[string]int, error) {
	// desiredReplicas = ceil[currentReplicas * ( currentMetricValue / desiredMetricValue )]
	replicas := make(map[string]int)

	for name, cpuU := range aggData.CPUUtilizations {
		currentReplicas := aggData.DeploymentInfos[name].Replica
		meanCPUUtilization, err := cpuU.GetMean()
		if err != nil {
			panic(err) //TODO maybe we should try again.
		}
		replicas[name] = int(math.Ceil(float64(currentReplicas) * (meanCPUUtilization / float64(h.hpaThreshold))))
	}
	h.addReplicaInfoToPreviousReplicaInfos(replicas)

	for name := range replicas {
		replicas[name] = int(math.Max(float64(replicas[name]), float64(getMaxPreviousReplicasForResource(name, h.previousReplicas))))
	}

	return replicas, nil
}

func (h *Hybrid) addReplicaInfoToPreviousReplicaInfos(newReplicaInfo map[string]int) {
	if h.cooldown == 0 {
		return
	}
	h.previousReplicas = append(h.previousReplicas, newReplicaInfo)
	if len(h.previousReplicas) > h.cooldown {
		h.previousReplicas = h.previousReplicas[1:h.cooldown]
	}
	if len(h.previousReplicas) > h.cooldown {
		panic(fmt.Sprintf("len(h.previousReplicas) > h.cooldown, %d, %d", len(h.previousReplicas), h.cooldown))
	}
}

func getMaxPreviousReplicasForResource(resource string, previousConfigs []map[string]int) int {
	maxValue := -1
	for _, m := range previousConfigs {
		if m[resource] > maxValue {
			maxValue = m[resource]
		}
	}
	return maxValue
}
