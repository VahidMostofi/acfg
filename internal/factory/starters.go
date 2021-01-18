package factory

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/aggregators/endpointsagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/sysstructureagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/ussageagg"
	"github.com/vahidmostofi/acfg/internal/aggregators/workloadagg"
	"github.com/vahidmostofi/acfg/internal/autocfg"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/dataaccess"
)

func newConfigDatabase()(dataaccess.ConfigDatabase,error){
	var err error
	var cd dataaccess.ConfigDatabase
	switch viper.Get(constants.CONFIG_CONFIGURATION_DATABASE_TYPE) {
	case "s3":
		region := viper.GetString(constants.CONFIG_CONFIGURATION_DATABASE_S3_REGION)
		bucket := viper.GetString(constants.CONFIG_CONFIGURATION_DATABASE_S3_BUCKET)
		cd, err = dataaccess.NewAWSConfigurationDatabase(region, bucket)
		if err != nil{
			return nil, errors.Wrapf(err, "error while creating s3 config database %s %s",region, bucket)
		}
		break
	default:
		return nil, errors.New(fmt.Sprintf("unknown ConfigDatabase type: %s", viper.Get(constants.CONFIG_CONFIGURATION_DATABASE_TYPE) ))
	}
	return cd, nil
}

func getEndpointsFilters()(map[string]map[string]interface{},error){
	endpointsFilter, ok := viper.Get(constants.CONFIG_ENDPOINTS_FILTERS).(map[string]map[string]interface{})
	if !ok {
		return nil, errors.Errorf("cant find endpoints filters in configs using: %s with type map[string]map[string]interface{}", constants.CONFIG_ENDPOINTS_FILTERS)
	}
	return endpointsFilter, nil
}

func newEndpointsAggregator() (*endpointsagg.EndpointsAggregator, error){
	// endpoints filters
	endpointsFilter, err := getEndpointsFilters()
	if err != nil{ return nil, err}
	epagArgs := map[string]interface{}{
		"url": viper.Get(constants.CONFIG_INFLUXDB_URL), //IF YOU CHANGED THIS, CHANGE THE ERROR BELOW
		"token": viper.Get(constants.CONFIG_INFLUXDB_TOKEN),
		"organization": viper.Get(constants.CONFIG_INFLUXDB_ORG),
		"bucket": viper.Get(constants.CONFIG_INFLUXDB_BUCKET),
	}
	ep, err := endpointsagg.NewEndpointsAggregator(viper.GetString(constants.CONFIG_ENDPOINTS_AGGREGATOR_TYPE), epagArgs, endpointsFilter)
	if err != nil{
		return nil, errors.Wrapf(err, "error while creating EndpointsAggregator, these might be useful: \"%s, %s, %s, %s, %v",
			viper.Get(constants.CONFIG_INFLUXDB_URL),
			"some token value",
			viper.Get(constants.CONFIG_INFLUXDB_ORG),
			viper.Get(constants.CONFIG_INFLUXDB_BUCKET),
			endpointsFilter,
		)
	}
	return ep, err
}

func getSystemStructure() (*sysstructureagg.SystemStructure, error){
	// system structure
	tempConverted, ok := viper.Get(constants.CONFIG_ENDPOINTS_2_RESOURCES).(map[string][]string)
	if !ok {
		return nil, errors.Errorf("cant find endpoints to resources in configs using: %s with type map[string]map[string]interface{}", constants.CONFIG_ENDPOINTS_2_RESOURCES)
	}
	ss, err := sysstructureagg.NewSystemStructure(viper.GetString(constants.CONFIG_SYSTEM_STRUCTURE_AGGREGATOR_TYPE), tempConverted)
	if err != nil{
		return nil, errors.Wrapf(err, "error while creating system structure aggregator, these might be useful: %s, %v", viper.GetString(constants.CONFIG_SYSTEM_STRUCTURE_AGGREGATOR_TYPE), tempConverted)
	}
	return ss, err
}

func getResourceFilters() (map[string]map[string]interface{},error){
	// resource filters
	resourceFilters, ok := viper.Get(constants.CONFIG_RESOURCE_FILTERS).(map[string]map[string]interface{})
	if !ok {
		return nil, errors.Errorf("cant find resource filters in configs using: %s with type map[string]map[string]interface{}", constants.CONFIG_RESOURCE_FILTERS)
	}
	return resourceFilters, nil
}

func newUsageAggregator()(*ussageagg.UsageAggregator,error){
	rfs, err := getResourceFilters()
	if err != nil{ return nil, err}

	// usage Aggregator
	uagArgs := map[string]interface{}{
		"url": viper.Get(constants.CONFIG_INFLUXDB_URL), //IF YOU CHANGED THIS, CHANGE THE ERROR BELOW
		"token": viper.Get(constants.CONFIG_INFLUXDB_TOKEN),
		"organization": viper.Get(constants.CONFIG_INFLUXDB_ORG),
		"bucket": viper.Get(constants.CONFIG_INFLUXDB_BUCKET),
	}
	ug, err := ussageagg.NewUsageAggregator(viper.GetString(constants.CONFIG_USAGE_AGGREGATOR_TYPE), uagArgs, rfs)
	if err != nil{
		return nil, errors.Wrapf(err, "error while creating EndpointsAggregator, these might be useful: \"%s, %s, %s, %s, %v",
			viper.Get(constants.CONFIG_INFLUXDB_URL),
			"some token value",
			viper.Get(constants.CONFIG_INFLUXDB_ORG),
			viper.Get(constants.CONFIG_INFLUXDB_BUCKET),
			rfs,
		)
	}
	return ug, err
}

func newWorkloadAggregator()(workloadagg.WorkloadAggregator, error){
	url := viper.GetString(constants.CONFIG_INFLUXDB_URL)
	token := viper.GetString(constants.CONFIG_INFLUXDB_TOKEN)
	organization := viper.GetString(constants.CONFIG_INFLUXDB_ORG)
	bucket := viper.GetString(constants.CONFIG_INFLUXDB_BUCKET)
	wg, err := workloadagg.NewInfluxDBWA(url, token, organization, bucket)
	if err != nil{
		return nil, errors.Wrapf(err, "error while creating workload aggregator with %s %s %s %s", url, "some token", organization, bucket)
	}
	return wg, err
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
	ug, err := newUsageAggregator()
	if err != nil{ return nil, err}

	// workload aggregator
	wg, err := newWorkloadAggregator()
	if err != nil{ return nil, err}

	// endpointsFilters
	epf, err := getEndpointsFilters()
	if err != nil {return nil, err}

	args := &autocfg.AutoConfigManagerArgs {
		Namespace: viper.GetString(constants.CONFIG_SYSTEM_NAMESPACE),
		DeploymentsToManage: viper.GetStringSlice(constants.CONFIG_SYSTEM_DEPLOYMENTS_TO_MANAGE),
		CfgValidation: autocfg.ConfigurationValidation{
			TotalAvailableMemory: viper.GetInt64(constants.CONFIG_VALIDATION_TOTAL_MEMORY),
			TotalAvailableCPU: viper.GetInt64(constants.CONFIG_VALIDATION_TOTAL_CPU),
		},
		UsingHash: viper.GetBool(constants.CONFIG_VALIDATION_USE_CACHE),
		ConfigDatabase: cd,
		WaitTimes: viper.Get(constants.CONFIG_WAIT_TIMES).(autocfg.WaitTimes),
		EndpointsAggregator: ep,
		SystemStructure: ss,
		UsageAggregator: ug,
		WorkloadAggregator: wg,
		EndpointsFilter: epf,
		StorePathPrefix: viper.GetString(constants.CONFIG_RESULTS_PREFIX_PATH),
	}

	acfgManager, err := autocfg.NewAutoConfigManager(args)
	if err != nil{
		return nil, errors.Wrap(err,"error while creating AutoConfigManager object")
	}

	return acfgManager, errors.Wrap(err, "error while creating AutoConfigManager")
}