package strategies

import (
	"encoding/csv"
	"fmt"
	"github.com/d4l3k/go-bayesopt"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/sla"
	"github.com/vahidmostofi/acfg/internal/workload"
	"math"
	"os"
	"strconv"
	"time"
)

type MemPSIThreshold struct {
	endpoints            []string
	resources            []string
	initialCPU           int64 // 1 CPU would be 1000
	initialMemory        int64 // 1 Gigabyte memory would be 1024
	utilizationThreshold float64
	utilizationIndicator string // mean
	csvWriter            *csv.Writer
	currentIteration     int64
	memBoChans           map[string]*BoConfig
	cpuBoChans           map[string]*BoConfig
}

func NewMemPSIThreshold(utilizationIndicator string, utilizationThreshold float64, endpoints []string, resources []string, initialCPU, initialMemory int64) (*MemPSIThreshold, error) {
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

	c := &MemPSIThreshold{
		endpoints:            endpoints,
		resources:            resources,
		initialCPU:           initialCPU,
		initialMemory:        initialMemory,
		utilizationThreshold: utilizationThreshold,
		utilizationIndicator: utilizationIndicator,
		csvWriter:            w,
		currentIteration:     1,
		memBoChans:           make(map[string]*BoConfig),
		cpuBoChans:           make(map[string]*BoConfig),
	}

	return c, nil
}

func (ct *MemPSIThreshold) AddSLA(sla *sla.SLA) error {
	return nil
}

func (ct *MemPSIThreshold) GetName() string {
	return "MemPSIThreshold"
}

var (
	memMaxAlloc = 10000.0
	memMinInc   = 500.0
	cpuMaxAlloc = 10000.0
	cpuMinInc   = 500.0
)

type BoConfig struct {
	toBoChan   chan float64
	fromBoChan chan float64
	doneBoChan chan float64
}

func (ct *MemPSIThreshold) GetInitialConfiguration(workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, error) {
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
		ct.memBoChans[resource] = &BoConfig{
			make(chan float64, 1),
			make(chan float64),
			make(chan float64),
		}
		ct.cpuBoChans[resource] = &BoConfig{
			make(chan float64, 1),
			make(chan float64),
			make(chan float64),
		}
		log.Infof("%s.GetInitialConfiguration(): initial config for %s: %v", ct.GetName(), resource, config[resource])
		if err := ct.csvWriter.Write([]string{"0", resource, strconv.FormatInt(*config[resource].CPU, 10), strconv.FormatInt(*config[resource].Memory, 10), "", "", "", ""}); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		resource := resource
		go func() {
			ct.initBO(resource, ct.memBoChans[resource], memMaxAlloc, memMinInc)
		}()
		go func() {
			ct.initBO(resource, ct.cpuBoChans[resource], cpuMaxAlloc, cpuMinInc)
		}()

	}

	return config, nil
}

func (ct *MemPSIThreshold) initBO(resource string, chans *BoConfig, maxAlloc float64, minInc float64) {
	_ = <-chans.toBoChan
	rounds := 6
	upperBound := 15.0
	lowerBound := 5.0

	steps := math.Floor(maxAlloc / minInc)
	allocParam := bayesopt.UniformParam{
		Max:  steps,
		Min:  1,
		Name: "alloc",
	}

	objective := func(params map[bayesopt.Param]float64) float64 {
		chans.fromBoChan <- params[allocParam]
		y := <-chans.toBoChan
		log.Printf("%s iteration x=%f, y=%f", resource, params[allocParam], y)
		if y > upperBound {
			return math.Abs(y)
		} else if y < lowerBound {
			return math.Abs(y)
		} else {
			return 0
		}
	}

	o := bayesopt.New(
		[]bayesopt.Param{
			allocParam,
		},
		bayesopt.WithMinimize(true),
		bayesopt.WithRounds(rounds),
	)

	x, _, err := o.Optimize(objective)
	if err != nil {
		log.Fatal(err)
	}

	chans.doneBoChan <- x[allocParam]
}

func (ct *MemPSIThreshold) ConfigureNextStep(currentConfig map[string]*configuration.Configuration, workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, map[string]interface{}, bool, error) {

	isChanged := false
	//var err error

	newConfig := make(map[string]*configuration.Configuration)
	for _, resource := range ct.resources {
		newConfig[resource] = currentConfig[resource].DeepCopy()

		whatToCompareMem, _ := aggData.MemPsiUtilizations[resource].GetMean()
		whatToCompareCPU, _ := aggData.CPUPSIUtilizations[resource].GetMean()
		// send to BO
		ct.memBoChans[resource].toBoChan <- whatToCompareMem
		ct.cpuBoChans[resource].toBoChan <- whatToCompareCPU

		// get request from BO twice (CPU & Mem)
		// I know this way is bad.
		for i := 0; i < 2; i++ {
			log.Infof("New loop itter for %s", resource)
			select {
			case doneResult := <-ct.memBoChans[resource].doneBoChan:
				isChanged = false
				log.Infof("%s Final Mem config: %f which is %f b", resource, doneResult, math.Ceil(doneResult)*memMinInc)
			case newMemIncs := <-ct.memBoChans[resource].fromBoChan:
				newMem := int64(math.Ceil(newMemIncs) * memMinInc)
				newConfig[resource].Memory = int64Ptr(newMem)
				log.Printf("%s Changed mem to: %f which is %d b-- PSI: %f", resource, newMemIncs, newMem, whatToCompareMem)
				isChanged = true
			case doneResult := <-ct.cpuBoChans[resource].doneBoChan:
				isChanged = false
				log.Infof("%s Final CPU config: %f which is %f b", resource, doneResult, math.Ceil(doneResult)*cpuMinInc)
			case newCpuIncs := <-ct.cpuBoChans[resource].fromBoChan:
				newCpu := int64(math.Ceil(newCpuIncs) * cpuMinInc)
				newConfig[resource].CPU = int64Ptr(newCpu)

				log.Printf("%s Changed CPU to: %f which is %d b-- PSI: %f", resource, newCpuIncs, newCpu, whatToCompareCPU)
				isChanged = true
			}
		}

		loginMean, _ := aggData.ResponseTimes["login"].GetMean()
		getBookMean, _ := aggData.ResponseTimes["get-book"].GetMean()
		editBookMean, _ := aggData.ResponseTimes["edit-book"].GetMean()
		if err := ct.csvWriter.Write([]string{strconv.FormatInt(ct.currentIteration, 10), resource, strconv.FormatInt(*newConfig[resource].CPU, 10), strconv.FormatInt(*newConfig[resource].Memory, 10), fmt.Sprintf("%f", loginMean), fmt.Sprintf("%f", getBookMean), fmt.Sprintf("%f", editBookMean), fmt.Sprintf("%f", whatToCompareCPU), fmt.Sprintf("%f", whatToCompareMem)}); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		ct.csvWriter.Flush()

	}

	ct.currentIteration += 1

	return newConfig, make(map[string]interface{}), isChanged, nil

	//
	///*
	//	threshold = 10
	//	thresholdRange = 2
	//	if use > 12 scale up
	//	if use < 8 scale down
	//	else end
	//*/
	//thresholdRange := 8.0
	//
	//for _, resource := range ct.resources {
	//
	//	//isChanged = true
	//
	//	var whatToCompareMem float64
	//	var whatToCompareCpu float64
	//	var err error
	//	if ct.utilizationIndicator == "mean" {
	//		whatToCompareMem, err = aggData.MemPsiUtilizations[resource].GetMean()
	//		if err != nil {
	//			return nil, make(map[string]interface{}), false, errors.Wrapf(err, "error while computing mean of mem psi utilizations for %s.", resource)
	//		}
	//
	//		whatToCompareCpu, err = aggData.CPUPSIUtilizations[resource].GetMean()
	//		if err != nil {
	//			return nil, make(map[string]interface{}), false, errors.Wrapf(err, "error while computing mean of cpu psi utilizations for %s.", resource)
	//		}
	//	}
	//	if whatToCompareMem > (ct.utilizationThreshold + thresholdRange) {
	//		newMem := int64Ptr(*newConfig[resource].Memory + 500)
	//		log.Infof("%s.ConfigureNextStep() mem psi utilization for %s is %f is more than %f changing mem... from %d to %d", ct.GetName(), resource, whatToCompareMem, ct.utilizationThreshold, *newConfig[resource].Memory, *newMem)
	//		newConfig[resource].Memory = newMem
	//		isChanged = true
	//	} else if whatToCompareMem < (ct.utilizationThreshold - thresholdRange) {
	//		if *newConfig[resource].Memory > 500 {
	//			newMem := int64Ptr(*newConfig[resource].Memory - 500)
	//			log.Infof("%s.ConfigureNextStep() mem psi utilization for %s is %f is less than %f changing mem... from %d to %d", ct.GetName(), resource, whatToCompareMem, ct.utilizationThreshold, *newConfig[resource].Memory, *newMem)
	//			newConfig[resource].Memory = newMem
	//			isChanged = true
	//		}
	//	}
	//
	//	if whatToCompareCpu > (ct.utilizationThreshold + thresholdRange) {
	//		newCpu := int64Ptr(*newConfig[resource].CPU + 500)
	//		log.Infof("%s.ConfigureNextStep() cpu psi utilization for %s is %f is more than %f changing cpu... from %d to %d", ct.GetName(), resource, whatToCompareCpu, ct.utilizationThreshold, *newConfig[resource].CPU, *newCpu)
	//		newConfig[resource].CPU = newCpu
	//		isChanged = true
	//	} else if whatToCompareCpu < (ct.utilizationThreshold - thresholdRange) {
	//		if *newConfig[resource].CPU > 500 {
	//			newCpu := int64Ptr(*newConfig[resource].CPU - 500)
	//			log.Infof("%s.ConfigureNextStep() cpu psi utilization for %s is %f is less than %f changing cpu... from %d to %d", ct.GetName(), resource, whatToCompareCpu, ct.utilizationThreshold, *newConfig[resource].CPU, *newCpu)
	//			newConfig[resource].CPU = newCpu
	//			isChanged = true
	//		}
	//	}
	//
	//	loginMean, _ := aggData.ResponseTimes["login"].GetMean()
	//	getBookMean, _ := aggData.ResponseTimes["get-book"].GetMean()
	//	editBookMean, _ := aggData.ResponseTimes["edit-book"].GetMean()
	//	if err := ct.csvWriter.Write([]string{strconv.FormatInt(ct.currentIteration, 10), resource, strconv.FormatInt(*newConfig[resource].CPU, 10), strconv.FormatInt(*newConfig[resource].Memory, 10), fmt.Sprintf("%f", loginMean), fmt.Sprintf("%f", getBookMean), fmt.Sprintf("%f", editBookMean), strconv.FormatInt(*newConfig[resource].Memory, 10), fmt.Sprintf("%f", whatToCompareCpu), fmt.Sprintf("%f", whatToCompareMem)}); err != nil {
	//		log.Fatalln("error writing record to file", err)
	//	}
	//	ct.csvWriter.Flush()
	//	//if isChanged {
	//	//	//log.Infof("%s.ConfigureNextStep() Mem PSI utilization for %s is %f is less than %f not changing replica from %d", ct.GetName(), resource, whatToCompare, ct.utilizationThreshold, *newConfig[resource].ReplicaCount)
	//	//	//log.Infof("%s.ConfigureNextStep() CPU utilization for %s is %f is less than %f not changing replica from %d", ct.GetName(), resource, whatToCompare, ct.utilizationThreshold, *newConfig[resource].ReplicaCount)
	//	//}
	//}
	//ct.currentIteration += 1
	//
	//return newConfig, make(map[string]interface{}), isChanged, nil

}
