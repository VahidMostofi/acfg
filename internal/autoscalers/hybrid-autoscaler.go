package autoscalers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/workload"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func NewHybridAutoscaler(endpoints, resources []string, hpaThreshold int64, predefinedReplicasFilepath string, usecount int) (Agent, error) {
	h := &Hybrid{
		endpoints:        endpoints,
		resources:        resources,
		hpaThreshold:     hpaThreshold,
		cooldown:         0,
		usecount:         usecount,
		previousReplicas: make([]map[string]int, 0),
		lastReplicas:     make(map[string]int),
	}

	if h.usecount > 0 {
		b, err := ioutil.ReadFile(predefinedReplicasFilepath)
		if err != nil {
			return nil, err
		}

		h.predefinedReplicas = make([]ReplicasForWorkloadRange, 0)
		err = json.Unmarshal(b, &h.predefinedReplicas)
		if err != nil {
			return nil, err
		}

		for _, rwr := range h.predefinedReplicas {
			if rwr.Replicas == nil || len(rwr.Replicas) == 0 {
				return nil, errors.Errorf("there must a replica configuration for each workload range.")
			}
		}
	}

	return h, nil
}

type Range struct {
	High float32 `json:"high"`
	Low  float32 `json:"low"`
}

type ReplicasForWorkloadRange struct {
	WorkloadRange map[string]Range `json:"workload-range"`
	Replicas      map[string]int   `json:"replicas"`
}

type Hybrid struct {
	predefinedReplicas []ReplicasForWorkloadRange
	hpaThreshold       int64
	cooldown           int // this is also the max length of previousReplicas
	endpoints          []string
	resources          []string
	previousReplicas   []map[string]int
	sess               *session.Session
	svc                *dynamodb.DynamoDB
	usecount           int
	lastReplicas       map[string]int
}

func (h *Hybrid) GetName() string {
	return "hybrid"
}

// checkForConfigAvailability //TODO write how this works
func (h *Hybrid) checkForConfigAvailability(happendWorkload workload.Workload) (map[string]int, error) {
	// fmt.Println(happendWorkload.GetMapStringInt())
	return nil, nil
}

func (h *Hybrid) Evaluate(aggData *aggregators.AggregatedData) (map[string]int, error) {
	takenApproach := "threshold"
	replicas := make(map[string]int)

	//--- Start Gathering supporting data for string details of autoscaler actions
	supportingData := make(map[string]interface{})
	supportingData["cpu-mean"] = make(map[string]float64)
	supportingData["current-replica"] = make(map[string]int)

	for name, cpuU := range aggData.CPUUtilizations {
		currentReplicas := aggData.DeploymentInfos[name].Replica
		meanCPUUtilization, err := cpuU.GetMean()
		if err != nil {
			supportingData["error"] = err
			break
		}
		(supportingData["cpu-mean"].(map[string]float64))[name] = meanCPUUtilization
		(supportingData["current-replica"].(map[string]int))[name] = currentReplicas
	}
	supportingData["workload"] = aggData.HappenedWorkload.String()
	//--- Finish Gathering supporting data for string details of autoscaler actions

	// if there is config (map[string]int) (replicas) for the workload that was happening, we use it,
	// otherwise, we go for the threshold approach
	suggestedReplicaCounts, err := h.checkForConfigAvailability(*aggData.HappenedWorkload)
	if err != nil {
		panic(err)
	}
	if suggestedReplicaCounts == nil {
		// THE FORMULA: desiredReplicas = ceil[currentReplicas * ( currentMetricValue / desiredMetricValue )]
		for name, cpuU := range aggData.CPUUtilizations {
			currentReplicas := aggData.DeploymentInfos[name].Replica
			meanCPUUtilization, err := cpuU.GetMean()
			if err != nil {
				replicas[name] = h.lastReplicas[name]
				log.Info("lastReplicas are being use due an error in getting mean cpu utilzation for " + name + ".")
			} else {
				replicas[name] = int(math.Ceil(float64(currentReplicas) * (meanCPUUtilization / float64(h.hpaThreshold))))
			}
		}
		// this is commented, because we are using the same feature in custom-pod-autoscaler
		// h.addReplicaInfoToPreviousReplicaInfos(replicas)
		// for name := range replicas {
		// 	replicas[name] = int(math.Max(float64(replicas[name]), float64(getMaxPreviousReplicasForResource(name, h.previousReplicas))))
		// }
	} else {
		replicas = suggestedReplicaCounts
	}

	h.logScalingDecision(takenApproach, replicas, supportingData)
	for key, value := range replicas {
		h.lastReplicas[key] = value
	}
	return replicas, nil
}

func (h *Hybrid) logScalingDecision(approach string, replicas map[string]int, supportingData map[string]interface{}) error {
	type ScalingDecisionInfo struct {
		Timestamp      int64                  `json:"timestamp"`
		Approach       string                 `json:"approach"`
		Replicas       map[string]int         `json:"replicas"`
		SupportingData map[string]interface{} `json:"supports"`
	}
	sdi := ScalingDecisionInfo{Approach: approach, Replicas: replicas, SupportingData: supportingData}
	sdi.Timestamp = time.Now().Unix()

	if h.sess == nil || h.svc == nil {
		// Initialize a session that the SDK will use to load
		// credentials from the shared credentials file ~/.aws/credentials
		// and region from the shared configuration file ~/.aws/config.
		h.sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		h.svc = dynamodb.New(h.sess)
	}

	av, err := dynamodbattribute.MarshalMap(sdi)
	if err != nil {
		return err
	}

	tableName := "autoscaling-decisions"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = h.svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	return nil
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
