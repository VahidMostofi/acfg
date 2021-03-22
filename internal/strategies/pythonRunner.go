package strategies

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/sla"
	"github.com/vahidmostofi/acfg/internal/workload"
)

// PythonRunner ...
type PythonRunner struct {
	PythonPath    string
	ScriptPath    string
	index         int
	cmd           *exec.Cmd
	stdin         io.WriteCloser
	configCh      chan map[string]serviceConfig
	resources     []string
	initialCPU    int64
	initialMemory int64
}

type serviceConfig struct {
	CPUAmount float64 `json:"cpu_count"`
	// ContainerCount int     `json:"container_count"`
	// WorkerCount    int     `json:"worker_count"`
}

type dataToSend struct {
	Feedbacks []float64 `json:"feedbacks"`
}

func (pr *PythonRunner) Write(p []byte) (int, error) {
	config := make(map[string]serviceConfig)
	err := json.Unmarshal(p, &config)
	if err != nil {
		log.Println("PythonRunner: non json response:", string(p))
		if strings.Trim(string(p), "\n") == "done" {
			log.Println("PythonRunner: Python is done")
			pr.configCh <- config
			return len(p), nil
		} else {
			log.Println("PythonRunner: from python:", string(p))
			return len(p), nil
		}
	}
	pr.configCh <- config
	return len(p), nil
}

// NewPythonRunner ...
func NewPythonRunner(pythonPath, scriptPath string, resources []string, initialCPU, initialMemory int64) (*PythonRunner, error) {
	pr := &PythonRunner{PythonPath: pythonPath, ScriptPath: scriptPath, resources: resources, initialCPU: initialCPU, initialMemory: initialMemory}
	return pr, nil
}

// AddSLA ...
func (pr *PythonRunner) AddSLA(sla *sla.SLA) error {
	return nil
}

// GetName ...
func (pr *PythonRunner) GetName() string {
	return "PythonRunner"
}

// GetInitialConfiguration ...
func (pr *PythonRunner) GetInitialConfiguration(workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, error) {
	config := make(map[string]*configuration.Configuration)
	for _, resource := range pr.resources {
		config[resource] = &configuration.Configuration{}
		config[resource].ReplicaCount = int64Ptr(1)
		config[resource].CPU = int64Ptr(pr.initialCPU)
		config[resource].Memory = int64Ptr(pr.initialMemory)
		config[resource].ResourceType = "Deployment"
		log.Infof("%s.GetInitialConfiguration(): initial config for %s: %v", pr.GetName(), resource, config[resource])
	}
	return config, nil
}

// ConfigureNextStep ...
func (pr *PythonRunner) ConfigureNextStep(currentConfig map[string]*configuration.Configuration, workload *workload.Workload, aggData *aggregators.AggregatedData) (map[string]*configuration.Configuration, bool, error) {
	isChanged := false
	if pr.index == 0 {
		pr.configCh = make(chan map[string]serviceConfig)
		log.Println("PythonRunner: first iteration of configurer")
		ctx, _ := context.WithCancel(context.Background())
		pr.cmd = exec.CommandContext(ctx, pr.PythonPath, "-W", "ignore", pr.ScriptPath)
		stdin, err := pr.cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		// defer stdin.Close()
		pr.stdin = stdin
		pr.cmd.Stdout = pr
		pr.cmd.Stderr = os.Stderr
		err = pr.cmd.Start()
		log.Println("PythonRunner: started python program")
		if err != nil {
			panic(err)
		}
	} else {
		values := make([]float64, 0)
		keys := make([]string, 0)
		for key := range aggData.ResponseTimes {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			rst := aggData.ResponseTimes[key]
			// log.Println("response time for", serviceName, "is", *rst.ResponseTimes95Percentile)
			if os.Args[1] == "theory" {
				panic("are you sure?")
				// values = append(values, *rst.ResponseTimesMean)
			} else {
				v, e := rst.GetPercentile(95)
				if e != nil {
					return nil, false, errors.Wrap(e, "error while getting 95 percentile of response time")
				}
				values = append(values, v)
			}
		}
		// fmt.Println("values-response-times", keys, values)
		feedbacks := &dataToSend{Feedbacks: values}
		b, err := json.Marshal(feedbacks)
		if err != nil {
			panic(fmt.Errorf("error while converting feedbacks to json: %w", err))
		}
		log.Println("PythonRunner: sending feedback:", string(b)+"\n")
		io.WriteString(pr.stdin, string(b)+"\n")
		log.Println("PythonRunner: sent feedback:", string(b))
	}
	var suggestedValues map[string]serviceConfig
	log.Println("PythonRunner: waiting for config")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case suggestedValues = <-pr.configCh:
				if len(suggestedValues) > 0 {
					log.Println("PythonRunner: got the config")
					log.Printf("PythonRunner: %v", suggestedValues)
				} else {
					log.Warn("PythonRunner: the len of suggested Values is 0")
				}
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	if len(suggestedValues) == 0 {
		return nil, false, nil
	}
	newConfig := make(map[string]*configuration.Configuration)
	for key := range currentConfig {
		newConfig[key] = currentConfig[key].DeepCopy()
		newConfig[key].UpdateEqualWithNewCPUValue(int64(suggestedValues[key].CPUAmount), 500)
		isChanged = true
	}
	pr.index++
	return newConfig, isChanged, nil
}
