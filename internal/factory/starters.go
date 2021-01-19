package factory

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/aggregators/endpointsagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/sysstructureagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/ussageagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/workloadagg"
	"github.com/vahidmostofi/acfg/internal/autocfg"
	"github.com/vahidmostofi/acfg/internal/constants"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func newConfigDatabase()(aggregators.ConfigDatabase,error){
	var err error
	var cd aggregators.ConfigDatabase
	switch viper.Get(constants.AutoConfigureCacheDatabaseType) {
	case "s3":
		region := viper.GetString(constants.AutoConfigureCacheS3Region)
		bucket := viper.GetString(constants.AutoConfigureCacheS3Bucket)
		if len(region) == 0 || len(bucket) == 0{
			panic("len(region) or len(bucket) is 0")
		}
		cd, err = aggregators.NewAWSConfigurationDatabase(region, bucket)
		if err != nil{
			return nil, errors.Wrapf(err, "error while creating s3 config database %s %s",region, bucket)
		}
		break
	default:
		return nil, errors.New(fmt.Sprintf("unknown ConfigDatabase type: %s", viper.Get(constants.AutoConfigureCacheDatabaseType) ))
	}
	return cd, nil
}

func getEndpointsFilters()(map[string]map[string]interface{},error){
	endpointsFilter, err := parseMapMapInterface(viper.GetStringMap(constants.EndpointsFilters))
	if err != nil {
		log.Errorf("this is found: %v", viper.Get(constants.EndpointsFilters))
		return nil, errors.Errorf("cant find endpoints filters in configs using %s with type map[string]map[string]interface{}", constants.EndpointsFilters)
	}
	return endpointsFilter, nil
}

func newEndpointsAggregator() (*endpointsagg.EndpointsAggregator, error){
	// endpoints filters
	endpointsFilter, err := getEndpointsFilters()
	if err != nil{ return nil, err}
	epagArgs := map[string]interface{}{
		"url": viper.Get(constants.EndpointsAggregatorArgsURL), //IF YOU CHANGED THIS, CHANGE THE ERROR BELOW
		"token": viper.Get(constants.EndpointsAggregatorArgsToken),
		"organization": viper.Get(constants.EndpointsAggregatorArgsOrganization),
		"bucket": viper.Get(constants.EndpointsAggregatorArgsBucket),
	}
	ep, err := endpointsagg.NewEndpointsAggregator(viper.GetString(constants.EndpointsAggregatorType), epagArgs, endpointsFilter)
	if err != nil{
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

func getSystemStructure() (*sysstructureagg.SystemStructure, error){
	// system structure
	tempConverted := viper.GetStringMapStringSlice(constants.SystemStructureAggregatorEndpoints2Resources)
	//if !ok {
	//	log.Errorf("this is found: %v", viper.Get(constants.SystemStructureAggregatorEndpoints2Resources))
	//	return nil, errors.Errorf("cant find endpoints to resources in configs using %s with type map[string][]string", constants.SystemStructureAggregatorEndpoints2Resources)
	//}
	ss, err := sysstructureagg.NewSystemStructure(viper.GetString(constants.SystemStructureAggregatorType), tempConverted)
	if err != nil{
		return nil, errors.Wrapf(err, "error while creating system structure aggregator, these might be useful: %s, %v", viper.GetString(constants.SystemStructureAggregatorType), tempConverted)
	}
	return ss, err
}

func getResourceFilters() (map[string]map[string]interface{},error){
	// resource filters
	resourceFilters, err := parseMapMapInterface(viper.GetStringMap(constants.ResourceFilters))
	if err != nil{
		log.Errorf("this is found: %v", viper.Get(constants.ResourceFilters))
		return nil, errors.Errorf("cant find resource filters in configs using: %s with type map[string]map[string]interface{}", constants.ResourceFilters)
	}
	return resourceFilters, nil
}

func newResourceUsageAggregator()(*ussageagg.UsageAggregator,error){
	rfs, err := getResourceFilters()
	if err != nil{ return nil, err}

	// usage Aggregator
	uagArgs := map[string]interface{}{
		"url": viper.Get(constants.ResourceUsageAggregatorArgsURL), //IF YOU CHANGED THIS, CHANGE THE ERROR BELOW
		"token": viper.Get(constants.ResourceUsageAggregatorArgsToken),
		"organization": viper.Get(constants.ResourceUsageAggregatorArgsOrganization),
		"bucket": viper.Get(constants.ResourceUsageAggregatorArgsBucket),
	}
	ug, err := ussageagg.NewUsageAggregator(viper.GetString(constants.ResourceUsageAggregatorType), uagArgs, rfs)
	if err != nil{
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

func newWorkloadAggregator()(workloadagg.WorkloadAggregator, error){
	if viper.GetString(constants.WorkloadAggregatorType) == "influxdb"{
		url := viper.GetString(constants.WorkloadAggregatorArgsURL)
		token := viper.GetString(constants.WorkloadAggregatorArgsToken)
		organization := viper.GetString(constants.WorkloadAggregatorArgsOrganization)
		bucket := viper.GetString(constants.WorkloadAggregatorArgsBucket)
		wg, err := workloadagg.NewInfluxDBWA(url, token, organization, bucket)
		if err != nil{
			return nil, errors.Wrapf(err, "error while creating workload aggregator with %s %s %s %s", url, "some token", organization, bucket)
		}
		return wg, err
	}
	return nil, errors.New("unknown worker aggregator type: " + viper.GetString(constants.WorkloadAggregatorType))
}

func getStoreDirectory() string{
	path := viper.GetString(constants.ResultsDirectory)
	parts := strings.Split(path, "/")
	for i := range parts{
		if len(parts[i]) < 1{
			continue
		}
		if parts[i][0] == '$'{
			parts[i] = viper.GetString(parts[i][1:])
		}
	}
	p := filepath.Join(parts...)
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil{
		panic(err)
	}
	return p
}

func getWaitTimes() autocfg.WaitTimes{
	w := autocfg.WaitTimes{
		WaitAfterConfigIsDeployed: time.Duration(viper.GetInt(constants.WaitTimesWaitAfterConfigIsDeployedSeconds)) * time.Second,
		LoadTestDuration: time.Duration(viper.GetInt(constants.WaitTimesLoadTestDurationSeconds)) * time.Second,
	}
	return w
}

// NewAutoConfigureManager returns new *AutoConfigManager
//
func NewAutoConfigureManager() (*autocfg.AutoConfigManager,error){

	// create a new config database
	cd, err := newConfigDatabase()
	if err != nil{ return nil, err}

	// create new endpoints aggregator
	ep, err := newEndpointsAggregator()
	if err != nil{ return nil, err}

	// system structure
	ss, err := getSystemStructure()
	if err != nil{ return nil, err}

	// usage aggregator
	ug, err := newResourceUsageAggregator()
	if err != nil{ return nil, err}

	// workload aggregator
	wg, err := newWorkloadAggregator()
	if err != nil{ return nil, err}

	// endpointsFilters
	epf, err := getEndpointsFilters()
	if err != nil {return nil, err}

	args := &autocfg.AutoConfigManagerArgs {
		Namespace: viper.GetString(constants.TargetSystemNamespace),
		DeploymentsToManage: viper.GetStringSlice(constants.TargetSystemDeploymentsToManage),
		CfgValidation: autocfg.ConfigurationValidation{
			TotalAvailableMemory: viper.GetInt64(constants.ConfigurationValidationTotalMemory),
			TotalAvailableCPU: viper.GetInt64(constants.ConfigurationValidationTotalCpu),
		},
		UsingHash: viper.GetBool(constants.AutoConfigureUseCache),
		ConfigDatabase: cd,
		WaitTimes: getWaitTimes(),
		EndpointsAggregator: ep,
		SystemStructure: ss,
		UsageAggregator: ug,
		WorkloadAggregator: wg,
		EndpointsFilter: epf,
		StorePathPrefix: getStoreDirectory(),
	}

	acfgManager, err := autocfg.NewAutoConfigManager(args)
	if err != nil{
		return nil, errors.Wrap(err,"error while creating AutoConfigManager object")
	}

	return acfgManager, errors.Wrap(err, "error while creating AutoConfigManager")
}

func parseMapMapInterface(in map[string]interface{}) (map[string]map[string]interface{},error){
	res := make(map[string]map[string]interface{})
	for key,value := range in{
		v, ok := value.(map[string]interface{})
		if !ok{
			return nil, errors.New(fmt.Sprintf("cant convert %s to %s", reflect.TypeOf(v), "map[string]interface{}"))
		}
		res[key] = v
	}
	return res, nil
}