package autocfg

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/aggregators/endpointsagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/sysstructureagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/ussageagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/workloadagg"
	"github.com/vahidmostofi/acfg/internal/autocfg/autoconfigurer"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/dataaccess"
	"github.com/vahidmostofi/acfg/internal/manager"
	"github.com/vahidmostofi/acfg/internal/workload"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type WaitTimes struct{
	WaitAfterConfigIsDeployed time.Duration //TODO
	LoadTestDuration time.Duration //TODO
}

type ConfigurationValidation struct{
	TotalAvailableCPU int64
	TotalAvailableMemory int64
}

type AutoConfigManager struct{
	ClusterManager *manager.K8sManager
	ConfigurationValidation ConfigurationValidation
	UsingHash bool
	ConfigDatabase dataaccess.ConfigDatabase
	WaitTimes WaitTimes
	EndpointsAggregator *endpointsagg.EndpointsAggregator
	SystemStructure *sysstructureagg.SystemStructure
	UsageAggregator *ussageagg.UsageAggregator
	WorkloadAggregator workloadagg.WorkloadAggregator
	EndpointsFilter map[string]map[string]interface{}
	StorePathPrefix string
}

func NewAutoConfigManager() (*AutoConfigManager,error){
	c, err := manager.NewK8Manager("default", []string{"auth","books","gateway"}) //TODO
	if err != nil{
		return nil, errors.Wrap(err, "error while creating kubernetes cluster manager.")
	}

	a := &AutoConfigManager{
		ClusterManager: c,
	}
	return a,nil
}

func (a *AutoConfigManager) aggregatedData(startTime, finishTime int64) (*AggregatedData, error){
	var err error

	// response times
	ag := &AggregatedData{}

	// Response Times
	ag.ResponseTimes, err = a.EndpointsAggregator.GetEndpointsResponseTimes(startTime, finishTime)
	if err != nil{
		return nil, errors.Wrap(err, "error while getting response times")
	}

	// Workload
	ag.HappenedWorkload, err = a.WorkloadAggregator.GetWorkload(startTime, finishTime, a.EndpointsFilter)
	if err != nil{
		return nil, errors.Wrap(err, "error while getting the workload that happened")
	}

	// System Structure
	ag.SystemStructure = a.SystemStructure

	// Resource Utilization
	// CPU
	ag.CPUUtilizations, err = a.UsageAggregator.GetAggregatedCPUUtilizations(startTime, finishTime)
	if err != nil{
		return nil, errors.Wrap(err, "error while getting the CPU utilizations")
	}
	// Memory TODO

	return ag, nil
}

func (a *AutoConfigManager) isConfigurationValid(cs map[string]*Configuration) (string,bool){
	var totalCPU int64
	var totalMemory int64
	for _,config := range cs{

		totalCPU += (*config.CPU) * *(config.ReplicaCount)
		totalMemory += (*config.Memory) * *(config.ReplicaCount)
	}

	if totalMemory > a.ConfigurationValidation.TotalAvailableMemory && a.ConfigurationValidation.TotalAvailableMemory > 0 {
		return "not enough memory", false
	}

	if totalCPU > a.ConfigurationValidation.TotalAvailableCPU && a.ConfigurationValidation.TotalAvailableCPU > 0 {
		return "not enough CPU", false
	}

	return "", true
}

func (a *AutoConfigManager) storeTestInformation(test *TestInformation) error{
	filePath := a.StorePathPrefix + test.Name // TODO clean, concat
	log.Infof("saving file at %s", filePath)
	fo, err := os.Create(filePath)
	if err != nil{
		return errors.Wrapf(err, "error while creating file for TestInformation")
	}
	// log where you saved it
	err = yaml.NewEncoder(fo).Encode(test)
	return errors.Wrap(err,"error while encoding testInforamation")
}


func (a *AutoConfigManager) Run(testName string, autoConfigAgent autoconfigurer.AutoConfigurationAgent, inputWorkload *workload.Workload) error {
	ctx, cnF := context.WithCancel(context.Background())

	testInformation := &TestInformation{
		Name: testName,
		AutoconfiguringApproach:autoConfigAgent.GetName(),
		Iterations: make([]*IterationInformation,0),
		InputWorkload: inputWorkload,
		VersionCode: viper.GetString(constants.CONFIG_VERSION_CODE),
	}

	log.Debug("AutoConfigManager.Run() waiting for all deployments to be available ")
	// make sure everything is up and running, all pods and deployments
	a.ClusterManager.WaitAllDeploymentsAreStable(ctx)

	// get the currentConfiguration from aut configuration . initialConfiguration()
	log.Debug("AutoConfigManager.Run() getting configuration with GetInitialConfiguration()")
	currentConfig, err := autoConfigAgent.GetInitialConfiguration(inputWorkload)
	if err != nil{
		return errors.Wrap(err, "error getting InitialConfiguration")
	}

	for{
		// one iteration starts here
		iterInfo := &IterationInformation{
			Configuration: currentConfig,
			AggregatedData: nil,
		}
		// get the hash of the configuration; if the hash and its value already exists in the database && we are using hash
		hashCode,err := GetHash(iterInfo.Configuration, testInformation.VersionCode)
		if err != nil{
			return errors.Wrap(err,"error while getting hash code from configuration")
		}

		if a.UsingHash {
			ag, err := a.ConfigDatabase.Retrieve(hashCode)
			if err != nil {
				return errors.Wrapf(err, "error while retrieving configuration with hash %s", hashCode)
			}
			if ag == nil {
				log.Debugf("AutoConfigManager.Run() no aggregatedData is found with hash code %s", hashCode)
			} else {
				log.Debugf("AutoConfigManager.Run() aggregatedData is found with hash code %s", hashCode)
				iterInfo.AggregatedData = ag
			}
		}
		// else: (no aggregated data is found in the ConfigDatabase)
		if iterInfo.AggregatedData != nil{

			// deploy the new configuration and wait for it to be deployed
			log.Infof("AutoConfigManager.Run() deploying the configuration")
			a.ClusterManager.UpdateConfigurationsAndWait(ctx, iterInfo.Configuration)
			log.Infof("AutoConfigManager.Run() configurations deployed and ready")

			log.Debugf("AutoConfigManager.Run() waiting %d seconds", a.WaitTimes.WaitAfterConfigIsDeployed)
			time.Sleep(a.WaitTimes.WaitAfterConfigIsDeployed)

			iterInfo.StartTime = time.Now().Unix()

			// start the load generator and wait a few seconds for it
			log.Debugf("AutoConfigManager.Run() load generator is starting")
			// TODO starting the load generator

			// wait for the specific duration and then stop the load generator
			log.Infof("AutoConfigManager.Run() load generator is started, waiting %d seconds", a.WaitTimes.LoadTestDuration)
			time.Sleep(a.WaitTimes.LoadTestDuration)
			// TODO stopping the load generator

			iterInfo.AggregatedData , err = a.aggregatedData(iterInfo.StartTime, iterInfo.FinishTime)
			if err != nil{
				return errors.Wrapf(err, "error while aggregating data from %d to %d", iterInfo.StartTime, iterInfo.FinishTime)
			}

			iterInfo.FinishTime = time.Now().Unix()

			// store the aggregated data
			err = a.ConfigDatabase.Store(hashCode, iterInfo.AggregatedData)
			if err != nil{
				return errors.Wrapf(err, "error while storing aggregated data for hash code %s", hashCode)
			}
		}

		testInformation.Iterations = append(testInformation.Iterations, iterInfo)
		a.storeTestInformation(testInformation)

		// pass all these information(data+) to the auto configuring agent and get the new configuration from it
		currentConfig, isDone, err := autoConfigAgent.ConfigureNextStep(inputWorkload)
		if err != nil{
			return errors.Wrap(err, "error while getting next configuration")
		}

		if reason, isValid := a.isConfigurationValid(currentConfig); !isValid{
			log.Infof("AutoConfigManager.Run() the new configuration is not valid because: %s; Breaking out of the loop", reason)
			break
		}

		if isDone{
			log.Infof("AutoConfigManager.Run() the autoconfiguring agent thinks we are done")
			break
		}
	}
}
