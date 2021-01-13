package autocfg

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/autocfg/autoconfigurer"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/manager"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type AutoConfigManager struct{
	ClusterManager *manager.K8sManager
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

func (a *AutoConfigManager) Run(autoConfigAgent autoconfigurer.AutoConfigurationAgent, inputWorkload *workload.Workload) error {
	ctx, cnF := context.WithCancel(context.Background())

	testInformation := &TestInformation{
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

	// one iteration starts here
	hashCode := currentConfig.GetHash(testInformation.VersionCode)
	// get the hash of the configuration; if the hash and its value already exists in the database && we are using hash
	// get the data+ from the database
	// else:
		// deploy the new configuration and wait for it to be deployed

		// start the load generator and wait a few seconds for it
		// record the start and finish time
		// wait for the specific duration and then stop the load generator

		// aggregate the results from the load generator
		// aggregate the results about the response times
		// aggregate the results about the workload that happened
		// aggregate the results about the resource utilization CPU-memory
		// aggregate information about system structure
		// all these data create the data+
		// pack data+ and append it as a step alongside with

	// pass all these information(data+) to the auto configuring agent and get the new configuration from it

}
