package factory

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	deploymentinfoagg "github.com/vahidmostofi/acfg/internal/aggregators/deploymentInfoAggregator"
	"github.com/vahidmostofi/acfg/internal/aggregators/endpointsagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/sysstructureagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/ussageagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/workloadagg"
	"github.com/vahidmostofi/acfg/internal/autocfg"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/loadgenerator"
	"github.com/vahidmostofi/acfg/internal/sla"
	"github.com/vahidmostofi/acfg/internal/workload"
	"gopkg.in/yaml.v2"
)

func newConfigDatabase() (aggregators.ConfigDatabase, error) {
	var err error
	var cd aggregators.ConfigDatabase
	switch viper.Get(constants.AutoConfigureCacheDatabaseType) {
	case "s3":
		region := viper.GetString(constants.AutoConfigureCacheS3Region)
		bucket := viper.GetString(constants.AutoConfigureCacheS3Bucket)
		if len(region) == 0 || len(bucket) == 0 {
			panic("len(region) or len(bucket) is 0")
		}
		cd, err = aggregators.NewAWSConfigurationDatabase(region, bucket)
		if err != nil {
			return nil, errors.Wrapf(err, "error while creating s3 config database %s %s", region, bucket)
		}
		break
	case "fs":
		directoryPath := viper.GetString(constants.AutoConfigureFSDirectory)
		if len(directoryPath) == 0 {
			panic("newConfigDatabase(): len(directoryPath) == 0")
		}
		cd = &aggregators.FSConfigurationDatabase{DirectoryName: directoryPath}
		break
	default:
		return nil, errors.New(fmt.Sprintf("unknown ConfigDatabase type: %s", viper.Get(constants.AutoConfigureCacheDatabaseType)))
	}
	return cd, nil
}

func GetEndpointsFilters() (map[string]map[string]interface{}, error) {
	endpointsFilter, err := parseMapMapInterface(viper.GetStringMap(constants.EndpointsFilters))
	if err != nil {
		log.Errorf("this is found: %v", viper.Get(constants.EndpointsFilters))
		return nil, errors.Errorf("cant find endpoints filters in configs using %s with type map[string]map[string]interface{}", constants.EndpointsFilters)
	}
	return endpointsFilter, nil
}

func newEndpointsAggregator() (*endpointsagg.EndpointsAggregator, error) {
	// endpoints filters
	endpointsFilter, err := GetEndpointsFilters()
	if err != nil {
		return nil, err
	}
	epagArgs := map[string]interface{}{
		"url":          viper.Get(constants.EndpointsAggregatorArgsURL), //IF YOU CHANGED THIS, CHANGE THE ERROR BELOW
		"token":        viper.Get(constants.EndpointsAggregatorArgsToken),
		"organization": viper.Get(constants.EndpointsAggregatorArgsOrganization),
		"bucket":       viper.Get(constants.EndpointsAggregatorArgsBucket),
	}
	ep, err := endpointsagg.NewEndpointsAggregator(viper.GetString(constants.EndpointsAggregatorType), epagArgs, endpointsFilter)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating EndpointsAggregator, these might be useful: \"%s, %s, %s, %s, %v",
			viper.Get(constants.EndpointsAggregatorArgsURL),
			"some token value",
			viper.Get(constants.EndpointsAggregatorArgsOrganization),
			viper.Get(constants.EndpointsAggregatorArgsBucket),
			endpointsFilter,
		)
	}
	return ep, err
}

func getSystemStructure() (*sysstructureagg.SystemStructure, error) {
	// system structure
	tempConverted := viper.GetStringMapStringSlice(constants.SystemStructureAggregatorEndpoints2Resources)
	//if !ok {
	// 	log.Errorf("this is found: %v", viper.Get(constants.SystemStructureAggregatorEndpoints2Resources))
	// 	return nil, errors.Errorf("cant find endpoints to resources in configs using %s with type map[string][]string", constants.SystemStructureAggregatorEndpoints2Resources)
	//}
	ss, err := sysstructureagg.NewSystemStructure(viper.GetString(constants.SystemStructureAggregatorType), tempConverted)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating system structure aggregator, these might be useful: %s, %v", viper.GetString(constants.SystemStructureAggregatorType), tempConverted)
	}
	return ss, err
}

func newDeploymentAggregator() (deploymentinfoagg.DeploymentInfoAggregator, error) {
	diagArgs := map[string]interface{}{
		"url":          viper.Get(constants.EndpointsAggregatorArgsURL), //IF YOU CHANGED THIS, CHANGE THE ERROR BELOW
		"token":        viper.Get(constants.EndpointsAggregatorArgsToken),
		"organization": viper.Get(constants.EndpointsAggregatorArgsOrganization),
		"bucket":       viper.Get(constants.EndpointsAggregatorArgsBucket),
	}

	resources, err := GetResources()
	if err != nil {
		return nil, err
	}

	di, err := deploymentinfoagg.NewDeploymentInfoAggregator(viper.GetString(constants.EndpointsAggregatorType), diagArgs, resources)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating DeploymentAggregator (EndpointsAggregator), these might be useful: \"%s, %s, %s, %s, %v",
			viper.Get(constants.EndpointsAggregatorArgsURL),
			"some token value",
			viper.Get(constants.EndpointsAggregatorArgsOrganization),
			viper.Get(constants.EndpointsAggregatorArgsBucket),
		)
	}
	return di, nil
}

func GetResourceFilters() (map[string]map[string]interface{}, error) {
	// resource filters
	resourceFilters, err := parseMapMapInterface(viper.GetStringMap(constants.ResourceFilters))
	if err != nil {
		log.Errorf("this is found: %v", viper.Get(constants.ResourceFilters))
		return nil, errors.Errorf("cant find resource filters in configs using: %s with type map[string]map[string]interface{}", constants.ResourceFilters)
	}
	return resourceFilters, nil
}

func GetResources() ([]string, error) {
	t, err := GetResourceFilters()
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)
	for v := range t {
		res = append(res, v)
	}
	return res, nil
}

func newResourceUsageAggregator() (*ussageagg.UsageAggregator, error) {
	rfs, err := GetResourceFilters()
	if err != nil {
		return nil, err
	}

	// usage Aggregator
	uagArgs := map[string]interface{}{
		"url":          viper.Get(constants.ResourceUsageAggregatorArgsURL), //IF YOU CHANGED THIS, CHANGE THE ERROR BELOW
		"token":        viper.Get(constants.ResourceUsageAggregatorArgsToken),
		"organization": viper.Get(constants.ResourceUsageAggregatorArgsOrganization),
		"bucket":       viper.Get(constants.ResourceUsageAggregatorArgsBucket),
	}
	ug, err := ussageagg.NewUsageAggregator(viper.GetString(constants.ResourceUsageAggregatorType), uagArgs, rfs)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating ResourceUsageAggregator, these might be useful: \"%s, %s, %s, %s, %v",
			viper.Get(constants.ResourceUsageAggregatorArgsURL),
			"some token value",
			viper.Get(constants.ResourceUsageAggregatorArgsOrganization),
			viper.Get(constants.ResourceUsageAggregatorArgsBucket),
			rfs,
		)
	}
	return ug, err
}

func newWorkloadAggregator() (workloadagg.WorkloadAggregator, error) {
	if viper.GetString(constants.WorkloadAggregatorType) == "influxdb" {
		epf, err := GetEndpointsFilters()
		if err != nil {
			return nil, err
		}

		url := viper.GetString(constants.WorkloadAggregatorArgsURL)
		token := viper.GetString(constants.WorkloadAggregatorArgsToken)
		organization := viper.GetString(constants.WorkloadAggregatorArgsOrganization)
		bucket := viper.GetString(constants.WorkloadAggregatorArgsBucket)
		wg, err := workloadagg.NewInfluxDBWA(url, token, organization, bucket, epf)
		if err != nil {
			return nil, errors.Wrapf(err, "error while creating workload aggregator with %s %s %s %s", url, "some token", organization, bucket)
		}
		return wg, err
	}
	return nil, errors.New("unknown worker aggregator type: " + viper.GetString(constants.WorkloadAggregatorType))
}

func getStoreDirectory() string {
	path := viper.GetString(constants.ResultsDirectory)
	parts := strings.Split(path, "/")
	for i := range parts {
		if len(parts[i]) < 1 {
			continue
		}
		if parts[i][0] == '$' {
			value := viper.GetString(parts[i][1:])
			if len(value) > 0 {
				parts[i] = value
			}
		}
	}
	p := filepath.Join(parts...)
	if path[0] == '/' {
		p = "/" + p
	}
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return p
}

func getWaitTimes() autocfg.WaitTimes {
	w := autocfg.WaitTimes{
		WaitAfterConfigIsDeployed:     time.Duration(viper.GetInt(constants.WaitTimesWaitAfterConfigIsDeployedSeconds)) * time.Second,
		LoadTestDuration:              time.Duration(viper.GetInt(constants.WaitTimesLoadTestDurationSeconds)) * time.Second,
		WaitAfterLoadGeneratorIsDone:  time.Duration(viper.GetInt(constants.WaitTimesWaitAfterLoadGeneratorIsDoneSeconds)) * time.Second,
		WaitAfterLoadGeneratorStartes: time.Duration(viper.GetInt(constants.WaitTimesWaitAfterLoadGeneratorStartes)) * time.Second,
	}
	return w
}

func NewAutoScalerManager() (*autocfg.AutoScalerManager, error) {
	a, err := NewAutoConfigureManager()
	return &autocfg.AutoScalerManager{
		Replicas:          make(map[string]int),
		AutoConfigManager: a,
	}, err
}

// NewAutoConfigureManager returns new *AutoConfigManager
//
func NewAutoConfigureManager() (*autocfg.AutoConfigManager, error) {

	// create a new config database
	cd, err := newConfigDatabase()
	if err != nil {
		return nil, err
	}

	// create new endpoints aggregator
	ep, err := newEndpointsAggregator()
	if err != nil {
		return nil, err
	}

	// system structure
	ss, err := getSystemStructure()
	if err != nil {
		return nil, err
	}

	// usage aggregator
	ug, err := newResourceUsageAggregator()
	if err != nil {
		return nil, err
	}

	// workload aggregator
	wg, err := newWorkloadAggregator()
	if err != nil {
		return nil, err
	}

	// endpointsFilters
	epf, err := GetEndpointsFilters()
	if err != nil {
		return nil, err
	}

	// get SLA
	sla, err := getSLA()
	if err != nil {
		return nil, err
	}

	wl := workload.GetTargetWorkload()
	// load generator
	lg, err := getLoadGenerator(wl)
	if err != nil {
		return nil, err
	}

	di, err := newDeploymentAggregator()
	if err != nil {
		return nil, err
	}

	args := &autocfg.AutoConfigManagerArgs{
		Namespace:           viper.GetString(constants.TargetSystemNamespace),
		DeploymentsToManage: viper.GetStringSlice(constants.TargetSystemDeploymentsToManage),
		CfgValidation: autocfg.ConfigurationValidation{
			TotalAvailableMemory: viper.GetInt64(constants.ConfigurationValidationTotalMemory),
			TotalAvailableCPU:    viper.GetInt64(constants.ConfigurationValidationTotalCpu),
		},
		UsingHash:                viper.GetBool(constants.AutoConfigureUseCache),
		ConfigDatabase:           cd,
		WaitTimes:                getWaitTimes(),
		EndpointsAggregator:      ep,
		SystemStructure:          ss,
		UsageAggregator:          ug,
		WorkloadAggregator:       wg,
		EndpointsFilter:          epf,
		StorePathPrefix:          getStoreDirectory(),
		SLA:                      sla,
		LoadGenerator:            lg,
		DeploymentInfoAggregator: di,
	}

	acfgManager, err := autocfg.NewAutoConfigManager(args)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating AutoConfigManager object")
	}

	return acfgManager, errors.Wrap(err, "error while creating AutoConfigManager")
}

func parseMapMapInterface(in map[string]interface{}) (map[string]map[string]interface{}, error) {
	res := make(map[string]map[string]interface{})
	for key, value := range in {
		v, ok := value.(map[string]interface{})
		if !ok {
			return nil, errors.New(fmt.Sprintf("cant convert %s to %s", reflect.TypeOf(v), "map[string]interface{}"))
		}
		res[key] = v
	}
	return res, nil
}

func getSLA() (*sla.SLA, error) {
	sla := &sla.SLA{
		Conditions: make([]sla.Condition, 0),
	}

	path := viper.GetString(constants.SLAConditionsFile)
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrapf(err, "error while reading sla conditions file %s", path)
	}

	err = yaml.NewDecoder(f).Decode(&sla.Conditions)
	if err != nil {
		return nil, errors.Wrapf(err, "error while parsing sla conditions file %s", path)
	}

	return sla, nil
}

func getLoadGenerator(w workload.Workload) (loadgenerator.LoadGenerator, error) {
	if strings.ToLower(viper.GetString(constants.LoadGeneratorType)) == "k6" {
		lg := &loadgenerator.K6LocalLoadGenerator{Args: w}
		r, err := os.Open(viper.GetString(constants.LoadGeneratorScriptPath))
		if err != nil {
			return nil, errors.Wrap(err, "error while getting load generator")
		}
		d, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, errors.Wrap(err, "error while reading load generator script")
		}
		lg.Data = d
		return lg, nil
	} else if strings.ToLower(viper.GetString(constants.LoadGeneratorType)) == "jmeter" {
		lg := &loadgenerator.JMeterLocalDocker{Command: viper.GetString(constants.LoadGeneratorCommand)}
		r, err := os.Open(viper.GetString(constants.LoadGeneratorScriptPath))
		if err != nil {
			return nil, errors.Wrap(err, "error while getting load generator (jmeter)")
		}
		d, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, errors.Wrap(err, "error while reading load generator script (jmeter)")
		}
		lg.Data = d
		return lg, nil
	}

	return nil, errors.New("unknown load generator type " + viper.GetString(constants.LoadGeneratorType))
}

type EnsembleAggregatorArgs struct {
	WithEndpointsAggregator  bool
	WithWorkloadAggregator   bool
	WithUsageAggregator      bool
	WithDeploymentAggregator bool
}

func NewEnsembleAggregator(eaa EnsembleAggregatorArgs) (*aggregators.Ensemble, error) {
	var err error
	e := &aggregators.Ensemble{}
	if eaa.WithEndpointsAggregator {
		e.Endpoints, err = newEndpointsAggregator()
		if err != nil {
			return nil, err
		}
	}

	if eaa.WithWorkloadAggregator {
		e.Workload, err = newWorkloadAggregator()
		if err != nil {
			return nil, err
		}
	}

	if eaa.WithUsageAggregator || true {
		e.Usage, err = newResourceUsageAggregator()
		if err != nil {
			return nil, err
		}
	}

	if eaa.WithDeploymentAggregator {
		e.DeploymentInfo, err = newDeploymentAggregator()
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}
