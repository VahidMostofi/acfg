package strategies

import (
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/sla"
	"github.com/vahidmostofi/acfg/internal/workload"
	"os"
	"strconv"
	"time"
)

type MemPSIThreshold_2 struct {
	endpoints            []string
	resources            []string
	initialCPU           int64 // 1 CPU would be 1000
	initialMemory        int64 // 1 Gigabyte memory would be 1024
	utilizationThreshold float64
	utilizationIndicator string // mean
	csvWriter            *csv.Writer
	currentIteration     int64
	prevMems             map[string][]int64
	prevCpus             map[string][]int64
}

func NewMemPSIThreshold_2(utilizationIndicator string, utilizationThreshold float64, endpoints []string, resources []string, initialCPU, initialMemory int64) (*MemPSIThreshold_2, error) {
	f, err := os.Create(fmt.Sprintf("/home/evan/Documents/5th-Year/Research/acfg-results/output-%d.csv", time.Now().Unix()))
	//defer func(f *os.File) {
	//	err := f.Close()
	//	if err != nil {
	//		log.Fatalln("failed to close file", err)
	//	}
	//}(f)

	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(f)
	//defer w.Flush()

	c := &MemPSIThreshold_2{
		endpoints:            endpoints,
		resources:            resources,
		initialCPU:           initialCPU,
		initialMemory:        initialMemory,
		utilizationThreshold: utilizationThreshold,
		utilizationIndicator: utilizationIndicator,
		csvWriter:            w,
		currentIteration:     1,
		prevMems:             make(map[string][]int64),
		prevCpus:             make(map[string][]int64),
	}

	return c, nil
}

func (ct *MemPSIThreshold_2) AddSLA(sla *sla.SLA) error {
	return nil
}

func (ct *MemPSIThreshold_2) GetName() string {
	return "MemPSIThreshold_2"
}

func (ct *MemPSIThreshold_2) GetInitialConfiguration(workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, error) {
	config := make(map[string]*configuration.Configuration)

	if err := ct.csvWriter.Write([]string{"iteration", "resource", "cpu_limit", "mem_limit", "response_times_login", "response_times_get_book", "response_times_edit_book", "cpu_psi", "mem_psi"}); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	ct.csvWriter.Flush()

	for _, resource := range ct.resources {
		config[resource] = &configuration.Configuration{}
		config[resource].ReplicaCount = int64Ptr(1)
		config[resource].CPU = int64Ptr(ct.initialCPU)
		config[resource].Memory = int64Ptr(ct.initialMemory)
		config[resource].ResourceType = "Deployment"
		log.Infof("%s.GetInitialConfiguration(): initial config for %s: %v", ct.GetName(), resource, config[resource])
		if err := ct.csvWriter.Write([]string{"0", resource, strconv.FormatInt(*config[resource].CPU, 10), strconv.FormatInt(*config[resource].Memory, 10), "", "", "", ""}); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		ct.prevMems[resource] = []int64{}
		ct.prevCpus[resource] = []int64{}
	}

	return config, nil
}

func (ct *MemPSIThreshold_2) ConfigureNextStep(currentConfig map[string]*configuration.Configuration, workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, map[string]interface{}, bool, error) {

	isChanged := false
	newConfig := make(map[string]*configuration.Configuration)

	/*
		threshold = 10
		thresholdRange = 2
		if use > 12 scale up
		if use < 8 scale down
		else end
	*/
	thresholdRange := 8.0

	for _, resource := range ct.resources {
		newConfig[resource] = currentConfig[resource].DeepCopy()

		//isChanged = true

		var whatToCompareMem float64
		var whatToCompareCpu float64
		var err error
		if ct.utilizationIndicator == "mean" {
			whatToCompareMem, err = aggData.MemPsiUtilizations[resource].GetMean()
			if err != nil {
				return nil, make(map[string]interface{}), false, errors.Wrapf(err, "error while computing mean of mem psi utilizations for %s.", resource)
			}

			whatToCompareCpu, err = aggData.CPUPSIUtilizations[resource].GetMean()
			if err != nil {
				return nil, make(map[string]interface{}), false, errors.Wrapf(err, "error while computing mean of cpu psi utilizations for %s.", resource)
			}
		}

		var newMem *int64

		if whatToCompareMem > (ct.utilizationThreshold + thresholdRange) {
			newMem = int64Ptr(*newConfig[resource].Memory + 500)
			log.Infof("%s.ConfigureNextStep() mem psi utilization for %s is %f is more than %f changing mem... from %d to %d", ct.GetName(), resource, whatToCompareMem, ct.utilizationThreshold, *newConfig[resource].Memory, *newMem)
		} else if whatToCompareMem < (ct.utilizationThreshold - thresholdRange) {
			if *newConfig[resource].Memory > 500 {
				newMem = int64Ptr(*newConfig[resource].Memory - 500)
				log.Infof("%s.ConfigureNextStep() mem psi utilization for %s is %f is less than %f changing mem... from %d to %d", ct.GetName(), resource, whatToCompareMem, ct.utilizationThreshold, *newConfig[resource].Memory, *newMem)
			}
		}

		// make change
		if newMem != nil {
			if len(ct.prevMems[resource]) > 2 && *newMem == ct.prevMems[resource][len(ct.prevMems[resource])-2] {
				log.Infof("%s Possible oscilation detected", resource)
			} else {
				newConfig[resource].Memory = newMem
				isChanged = true
			}
		}

		if whatToCompareCpu > (ct.utilizationThreshold + thresholdRange) {
			newCpu := int64Ptr(*newConfig[resource].CPU + 500)
			log.Infof("%s.ConfigureNextStep() cpu psi utilization for %s is %f is more than %f changing cpu... from %d to %d", ct.GetName(), resource, whatToCompareCpu, ct.utilizationThreshold, *newConfig[resource].CPU, *newCpu)
			newConfig[resource].CPU = newCpu
			isChanged = true
		} else if whatToCompareCpu < (ct.utilizationThreshold - thresholdRange) {
			if *newConfig[resource].CPU > 500 {
				newCpu := int64Ptr(*newConfig[resource].CPU - 500)
				log.Infof("%s.ConfigureNextStep() cpu psi utilization for %s is %f is less than %f changing cpu... from %d to %d", ct.GetName(), resource, whatToCompareCpu, ct.utilizationThreshold, *newConfig[resource].CPU, *newCpu)
				newConfig[resource].CPU = newCpu
				isChanged = true
			}
		}

		ct.prevMems[resource] = append(ct.prevMems[resource], *newConfig[resource].Memory)
		ct.prevCpus[resource] = append(ct.prevCpus[resource], *newConfig[resource].CPU)

		loginMean, _ := aggData.ResponseTimes["login"].GetMean()
		getBookMean, _ := aggData.ResponseTimes["get-book"].GetMean()
		editBookMean, _ := aggData.ResponseTimes["edit-book"].GetMean()
		if err := ct.csvWriter.Write([]string{strconv.FormatInt(ct.currentIteration, 10), resource, strconv.FormatInt(*newConfig[resource].CPU, 10), strconv.FormatInt(*newConfig[resource].Memory, 10), fmt.Sprintf("%f", loginMean), fmt.Sprintf("%f", getBookMean), fmt.Sprintf("%f", editBookMean), strconv.FormatInt(*newConfig[resource].Memory, 10), fmt.Sprintf("%f", whatToCompareCpu), fmt.Sprintf("%f", whatToCompareMem)}); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		ct.csvWriter.Flush()
		//if isChanged {
		//	//log.Infof("%s.ConfigureNextStep() Mem PSI utilization for %s is %f is less than %f not changing replica from %d", ct.GetName(), resource, whatToCompare, ct.utilizationThreshold, *newConfig[resource].ReplicaCount)
		//	//log.Infof("%s.ConfigureNextStep() CPU utilization for %s is %f is less than %f not changing replica from %d", ct.GetName(), resource, whatToCompare, ct.utilizationThreshold, *newConfig[resource].ReplicaCount)
		//}
	}
	ct.currentIteration += 1

	return newConfig, make(map[string]interface{}), isChanged, nil

}
