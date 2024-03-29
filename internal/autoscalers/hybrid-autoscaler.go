package autoscalers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
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
		fmt.Println("len of predefined replicas: ", len(h.predefinedReplicas))

		for _, rwr := range h.predefinedReplicas {
			if rwr.Replicas == nil || len(rwr.Replicas) == 0 {
				return nil, errors.Errorf("there must a replica configuration for each workload range.")
			}
		}

		for _, pdr := range h.predefinedReplicas {
			for key, value := range pdr.WorkloadRange {
				pdr.WorkloadRange[strings.ReplaceAll(key, "-", "")] = value
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

// func (h *Hybrid) getMaxOfAllLessOrEqual(happendWorkload workload.Workload) (map[string]int, int, error) {
// 	options := make([]map[string]int, 0)
// 	for i, pdr := range h.predefinedReplicas {
// 		if i >= h.usecount {
// 			continue
// 		}
// 		flag := true
// 		for endpoint, requestCount := range happendWorkload.GetMapStringInt() {
// 			endpoint = strings.ReplaceAll(endpoint, "-", "")
// 			rcf := float32(requestCount)
// 			log.Debugf("%f %f", rcf, math.Ceil(float64(pdr.WorkloadRange[endpoint].High)))
// 			if !(math.Ceil(float64(pdr.WorkloadRange[endpoint].High)) >= float64(rcf)+10) {
// 				flag = false
// 				break
// 			}
// 		}
// 		if flag {
// 			fmt.Println("found an option predefineds (just less than)", pdr.Replicas)
// 			options = append(options, pdr.Replicas)
// 		}
// 	}
// 	if len(options) > 0 {
// 		maxSum := 0
// 		maxIdx := 0
// 		for idx, replicas := range options {
// 			s := 0
// 			for _, v := range replicas {
// 				s += v
// 			}
// 			if s > maxSum {
// 				maxSum = s
// 				maxIdx = idx
// 			}
// 		}
// 		return options[maxIdx], maxSum, nil
// 	}
// 	return nil, -1, nil
// }

// checkForConfigAvailability //TODO write how this works
func (h *Hybrid) checkForConfigAvailability(happendWorkload workload.Workload) (map[string]int, error) {
	for i, pdr := range h.predefinedReplicas {
		if i >= h.usecount {
			continue
		}
		flag := true
		for endpoint, requestCount := range happendWorkload.GetMapStringInt() {
			endpoint = strings.ReplaceAll(endpoint, "-", "")
			rcf := float32(requestCount)
			log.Debugf("%f %f %f", pdr.WorkloadRange[endpoint].Low, rcf, pdr.WorkloadRange[endpoint].High)
			if !(pdr.WorkloadRange[endpoint].Low <= rcf && pdr.WorkloadRange[endpoint].High >= rcf) {
				flag = false
				break
			}
		}
		if flag {
			fmt.Println("found by predefineds", pdr.Replicas)
			return pdr.Replicas, nil
		}
	}
	return nil, nil
}

var violationCount int

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
			supportingData["error"] = err.Error()
			break
		}
		(supportingData["cpu-mean"].(map[string]float64))[name] = meanCPUUtilization
		(supportingData["current-replica"].(map[string]int))[name] = currentReplicas
	}
	supportingData["workload"] = aggData.HappenedWorkload.String()
	//--- Finish Gathering supporting data for string details of autoscaler actions

	// if there is config (map[string]int) (replicas) for the workload that was happening, we use it,
	// otherwise, we go for the threshold approach
	fmt.Println()
	vvv, err := aggData.ResponseTimes["login"].GetPercentile(95)
	if err == nil {
		if vvv > 0.25 {
			violationCount++
			fmt.Println("login VIOLATIOOONNNNN", violationCount)
		}
	}
	vvv, err = aggData.ResponseTimes["get-book"].GetPercentile(95)
	if err == nil {
		if vvv > 0.025 {
			violationCount++
			fmt.Println("get-book VIOLATIOOONNNNN", violationCount)
		}
	}
	vvv, err = aggData.ResponseTimes["edit-book"].GetPercentile(95)
	if err == nil {
		if vvv > 0.025 {
			violationCount++
			fmt.Println("edit-book VIOLATIOOONNNNN", violationCount)
		}
	}
	fmt.Println(aggData.HappenedWorkload)
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
		// this ia feture I've been looking into, I don't use it now.
		// maxOtherOptions, totalOfOtherOption, err := h.getMaxOfAllLessOrEqual(*aggData.HappenedWorkload)
		// if maxOtherOptions != nil && totalOfOtherOption > 0 {
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	hpaAllocatedValue := 0
		// 	for _, v := range replicas {
		// 		hpaAllocatedValue += v
		// 	}
		// 	log.Debugf("total required by HPA is %d", hpaAllocatedValue)
		// 	fmt.Printf("total required by HPA is %d \n", hpaAllocatedValue)

		// 	if hpaAllocatedValue > totalOfOtherOption {
		// 		replicas = maxOtherOptions
		// 		log.Debugf("we can do better by max of the predefineds %d %v.", totalOfOtherOption, maxOtherOptions)
		// 		fmt.Printf("we can do better by max of the predefineds %d %v.\n", totalOfOtherOption, maxOtherOptions)
		// 		supportingData["with-less-than-predefineds-max"] = replicas
		// 		takenApproach = "max-of-predefineds"
		// 	}
		// }

		// this is commented, because we are using the same feature in custom-pod-autoscaler
		// h.addReplicaInfoToPreviousReplicaInfos(replicas)
		// for name := range replicas {
		// 	replicas[name] = int(math.Max(float64(replicas[name]), float64(getMaxPreviousReplicasForResource(name, h.previousReplicas))))
		// }
	} else {
		takenApproach = "predefined"
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
