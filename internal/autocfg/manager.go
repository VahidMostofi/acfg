package autocfg

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/aggregators/endpointsagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/sysstructureagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/ussageagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/workloadagg"
	"github.com/vahidmostofi/acfg/internal/clustermanager"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/loadgenerator"
	"github.com/vahidmostofi/acfg/internal/sla"
	"github.com/vahidmostofi/acfg/internal/strategies"
	"github.com/vahidmostofi/acfg/internal/workload"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type WaitTimes struct{ // if anything is added to this, you need to update starter which creates this object from config files
	WaitAfterConfigIsDeployed time.Duration
	LoadTestDuration time.Duration
	WaitAfterLoadGeneratorIsDone time.Duration
	WaitAfterLoadGeneratorStartes time.Duration
}

type ConfigurationValidation struct{
	TotalAvailableCPU int64
	TotalAvailableMemory int64
}

type AutoConfigManager struct{
	clusterManager          clustermanager.ClusterManager
	configurationValidation ConfigurationValidation
	usingHash               bool
	configDatabase          aggregators.ConfigDatabase
	waitTimes               WaitTimes
	endpointsAggregator     *endpointsagg.EndpointsAggregator
	systemStructure         *sysstructureagg.SystemStructure
	usageAggregator         *ussageagg.UsageAggregator
	workloadAggregator      workloadagg.WorkloadAggregator
	endpointsFilter         map[string]map[string]interface{}
	storePathPrefix         string
	cancelFunc              context.CancelFunc
	sla						*sla.SLA
	lg 						loadgenerator.LoadGenerator
}

type AutoConfigManagerArgs struct{
	Namespace           string
	DeploymentsToManage []string
	CfgValidation       ConfigurationValidation
	UsingHash           bool
	ConfigDatabase      aggregators.ConfigDatabase
	WaitTimes           WaitTimes
	EndpointsAggregator *endpointsagg.EndpointsAggregator
	SystemStructure     *sysstructureagg.SystemStructure
	UsageAggregator     *ussageagg.UsageAggregator
	WorkloadAggregator  workloadagg.WorkloadAggregator
	EndpointsFilter     map[string]map[string]interface{}
	StorePathPrefix     string
	SLA 				*sla.SLA
	LoadGenerator		loadgenerator.LoadGenerator
}

func NewAutoConfigManager(args *AutoConfigManagerArgs) (*AutoConfigManager,error){
	c, err := clustermanager.NewK8ClusterManager(args.Namespace, args.DeploymentsToManage)
	if err != nil{
		return nil, errors.Wrap(err, "error while creating kubernetes cluster clustermanager.")
	}

	a := &AutoConfigManager{
		clusterManager: c,
		configurationValidation: args.CfgValidation,
		usingHash: args.UsingHash,
		configDatabase: args.ConfigDatabase,
		waitTimes: args.WaitTimes,
		endpointsAggregator: args.EndpointsAggregator,
		systemStructure: args.SystemStructure,
		usageAggregator: args.UsageAggregator,
		workloadAggregator: args.WorkloadAggregator,
		endpointsFilter: args.EndpointsFilter,
		storePathPrefix: args.StorePathPrefix,
		sla: args.SLA,
		lg: args.LoadGenerator,
	}
	return a,nil
}

func (a *AutoConfigManager) aggregatedData(startTime, finishTime int64) (*aggregators.AggregatedData, error){
	log.Debugf("aggregating data")
	var err error

	// response times
	ag := &aggregators.AggregatedData{}

	// Response Times
	ag.ResponseTimes, err = a.endpointsAggregator.GetEndpointsResponseTimes(startTime, finishTime)
	if err != nil{
		return nil, errors.Wrap(err, "error while getting response times")
	}

	// Workload
	ag.HappenedWorkload, err = a.workloadAggregator.GetWorkload(startTime, finishTime, a.endpointsFilter)
	if err != nil{
		return nil, errors.Wrap(err, "error while getting the workload that happened")
	}

	// System Structure
	ag.SystemStructure = a.systemStructure

	// Resource Utilization
	// CPU
	ag.CPUUtilizations, err = a.usageAggregator.GetAggregatedCPUUtilizations(startTime, finishTime)
	if err != nil{
		return nil, errors.Wrap(err, "error while getting the CPU utilizations")
	}
	// Memory TODO

	return ag, nil
}

func (a *AutoConfigManager) isConfigurationValid(cs map[string]*configuration.Configuration) (string,bool){
	var totalCPU int64
	var totalMemory int64
	for _,config := range cs{

		totalCPU += (*config.CPU) * *(config.ReplicaCount)
		totalMemory += (*config.Memory) * *(config.ReplicaCount)
	}

	if totalMemory > a.configurationValidation.TotalAvailableMemory && a.configurationValidation.TotalAvailableMemory > 0 {
		return "not enough memory", false
	}
	// CPU in config should be times 1000
	maxAvailableCPU := a.configurationValidation.TotalAvailableCPU * 1000
	if totalCPU > maxAvailableCPU && a.configurationValidation.TotalAvailableCPU > 0 {
		return fmt.Sprintf("not enough CPU: %d vs %d", totalCPU, maxAvailableCPU), false
	}
	return "", true
}

func (a *AutoConfigManager) storeTestInformation(test *TestInformation) error{
	a.storePathPrefix = strings.ReplaceAll(a.storePathPrefix, "$STRATEGY.NAME", viper.GetString(constants.StrategyName))
	err := os.MkdirAll(filepath.Clean(a.storePathPrefix), os.ModePerm)
	if err != nil{
		return errors.Wrapf(err, "there was error creating necessary directories.")
	}
	filePath := filepath.Join(filepath.Clean(a.storePathPrefix), test.Name + ".yaml") // TODO is it always local? no s3?
	log.Infof("saving file at %s", filePath)
	fo, err := os.Create(filePath)
	if err != nil{
		return errors.Wrapf(err, "error while creating file for TestInformation")
	}
	// log where you saved it
	err = yaml.NewEncoder(fo).Encode(test)
	return errors.Wrap(err,"error while encoding testInforamation")
}


func (a *AutoConfigManager) Run(testName string, autoConfigStrategyAgent strategies.Strategy, inputWorkload *workload.Workload) error {
	ctx, cnF := context.WithCancel(context.Background())
	a.cancelFunc = cnF

	// adding SLA to strategy
	autoConfigStrategyAgent.AddSLA(a.sla)

	testInformation := &TestInformation{
		Name: testName,
		AutoconfiguringApproach:autoConfigStrategyAgent.GetName(),
		Iterations: make([]*IterationInformation,0),
		InputWorkload: inputWorkload,
		VersionCode: viper.GetString(constants.VersionCode),
		AllSettings: viper.AllSettings(),
	}

	log.Debug("AutoConfigManager.Run() waiting for all deployments to be available ")
	// make sure everything is up and running, all pods and deployments
	a.clusterManager.WaitAllDeploymentsAreStable(ctx)

	// get the currentConfiguration from aut configuration . initialConfiguration()
	log.Debug("AutoConfigManager.Run() getting configuration with GetInitialConfiguration()")
	currentConfig, err := autoConfigStrategyAgent.GetInitialConfiguration(inputWorkload, nil) // TODO aggData is nil
	if err != nil{
		return errors.Wrap(err, "error getting InitialConfiguration")
	}
	iterationId := 0
	for{
		iterationId++
		log.Infof("starting iteration %d", iterationId)
		// one iteration starts here
		iterInfo := &IterationInformation{
			Configuration: currentConfig,
			AggregatedData: nil,
			LoadGeneratorFeedback: nil,
		}

		if a.usingHash {
			// get the hash of the configuration; if the hash and its value already exists in the database && we are using hash
			hashCode,err := GetHash(iterInfo.Configuration, testInformation.VersionCode, inputWorkload)
			if err != nil{
				return errors.Wrap(err,"error while getting hash code from configuration")
			}

			ag, err := a.configDatabase.Retrieve(hashCode)
			if err != nil {
				return errors.Wrapf(err, "error while retrieving configuration with hash %s", hashCode)
			}
			if ag == nil {
				log.Infof("AutoConfigManager.Run() no aggregatedData is found with hash code %s", hashCode)
			} else {
				log.Infof("AutoConfigManager.Run() aggregatedData is found with hash code %s", hashCode)
				iterInfo.AggregatedData = ag
			}
		}
		// else: (no aggregated data is found in the ConfigDatabase)
		if iterInfo.AggregatedData == nil{

			// deploy the new configuration and wait for it to be deployed
			log.Infof("AutoConfigManager.Run() checking configuration before deploying")
			reason, ok := a.isConfigurationValid(iterInfo.Configuration)
			if !ok{
				log.Infof("breaking the loop because configuration is not valid: %s", reason)
				break
			}
			log.Infof("AutoConfigManager.Run() deploying the configuration")
			a.clusterManager.UpdateConfigurationsAndWait(ctx, iterInfo.Configuration)
			log.Infof("AutoConfigManager.Run() configurations deployed and ready")

			log.Infof("AutoConfigManager.Run() waiting %s.", a.waitTimes.WaitAfterConfigIsDeployed.String())
			time.Sleep(a.waitTimes.WaitAfterConfigIsDeployed)			

			// start the load generator and wait a few seconds for it
			log.Debugf("AutoConfigManager.Run() load generator is starting")
			a.lg.Start(inputWorkload)
			time.Sleep(a.waitTimes.WaitAfterLoadGeneratorStartes)

			iterInfo.StartTime = time.Now().Unix()

			// wait for the specific duration and then stop the load generator
			log.Infof("AutoConfigManager.Run() load generator is started, waiting %s while load generator is running.", a.waitTimes.LoadTestDuration.String())
			time.Sleep(a.waitTimes.LoadTestDuration)
			a.lg.Stop()
			iterInfo.LoadGeneratorFeedback, err = a.lg.GetFeedback()

			iterInfo.FinishTime = time.Now().Unix()
			log.Infof("AutoConfigManager.Run() load generator is done, waiting %s.", a.waitTimes.WaitAfterLoadGeneratorIsDone.String())
			time.Sleep(a.waitTimes.WaitAfterLoadGeneratorIsDone)
			iterInfo.AggregatedData , err = a.aggregatedData(iterInfo.StartTime, iterInfo.FinishTime)
			if err != nil{
				return errors.Wrapf(err, "error while aggregating data from %d to %d", iterInfo.StartTime, iterInfo.FinishTime)
			}

			if a.usingHash{
				hashCode,err := GetHash(iterInfo.Configuration, testInformation.VersionCode, inputWorkload)
				if err != nil{
					return errors.Wrap(err,"error while getting hash code from configuration")
				}

				// store the aggregated data
				err = a.configDatabase.Store(hashCode, iterInfo.AggregatedData)
				if err != nil{
					return errors.Wrapf(err, "error while storing aggregated data for hash code %s", hashCode)
				}
			}
		}

		// at this point we have the aggregated data, we either found it with cache or by running the load generator
		// so we print some info to the log about this iteration and the aggregated data in it.
		for endpointName, responseTimes := range iterInfo.AggregatedData.ResponseTimes{
			log.Infof("response times for %s: %s", endpointName, responseTimes.String())
		}
		log.Infof("workload that happend during this iteration: %s", iterInfo.AggregatedData.HappenedWorkload.String())

		testInformation.Iterations = append(testInformation.Iterations, iterInfo)
		err = a.storeTestInformation(testInformation)
		if err != nil{
			return errors.Wrapf(err, "error while saving aggregated results.")
		}

		// pass all these information(data+) to the auto configuring agent and get the new configuration from it
		currentConfig, isChanged, err := autoConfigStrategyAgent.ConfigureNextStep(iterInfo.Configuration, inputWorkload, iterInfo.AggregatedData)
		isDone := !isChanged
		doneReason := ""
		iterInfo.Configuration = currentConfig
		if err != nil{
			return errors.Wrap(err, "error while getting next configuration")
		}

		if reason, isValid := a.isConfigurationValid(iterInfo.Configuration); !isValid{
			log.Infof("AutoConfigManager.Run() the new configuration is not valid because: %s; Breaking out of the loop", reason)
			break
		}

		if isDone{
			doneReason += "the strategy agent is done."
			log.Infof("AutoConfigManager.Run() the autoconfiguring agent thinks we are done")
			break
		}
	}
	// TODO store the reason we are going out of the loop, cant do better? the config is not valid? we are done?
	// TODO we should listen to signals and undeploy everything. Graceful shutdown.
	return nil
}

// TODO use this!
func CheckCondition(data *aggregators.AggregatedData, condition sla.Condition) (bool, error){
	if condition.Type == "ResponseTime"{
		value := condition.GetComputeFunction()(*data.ResponseTimes[condition.EndpointName])
		if value <= condition.Threshold{
			return true, nil
		}
	}
	return false, nil
}