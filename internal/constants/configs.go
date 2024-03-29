package constants

const VersionCode = "VersionCode"

const EndpointsAggregatorType = "EndpointsAggregator.Type"
const EndpointsAggregatorArgsURL = "EndpointsAggregator.Args.URL"
const EndpointsAggregatorArgsToken = "EndpointsAggregator.Args.Token"
const EndpointsAggregatorArgsBucket = "EndpointsAggregator.Args.Bucket"
const EndpointsAggregatorArgsOrganization = "EndpointsAggregator.Args.Organization"

const ResourceUsageAggregatorType = "ResourceUsageAggregator.Type"
const ResourceUsageAggregatorArgsURL = "ResourceUsageAggregator.Args.URL"
const ResourceUsageAggregatorArgsToken = "ResourceUsageAggregator.Args.Token"
const ResourceUsageAggregatorArgsBucket = "ResourceUsageAggregator.Args.Bucket"
const ResourceUsageAggregatorArgsOrganization = "ResourceUsageAggregator.Args.Organization"

const WorkloadAggregatorType = "WorkloadAggregator.Type"
const WorkloadAggregatorArgsURL = "WorkloadAggregator.Args.URL"
const WorkloadAggregatorArgsToken = "WorkloadAggregator.Args.Token"
const WorkloadAggregatorArgsBucket = "WorkloadAggregator.Args.Bucket"
const WorkloadAggregatorArgsOrganization = "WorkloadAggregator.Args.Organization"

const SystemStructureAggregatorType = "SystemStructureAggregator.Type"
const SystemStructureAggregatorEndpoints2Resources = "SystemStructureAggregator.Endpoints2Resources"

const ResourceFilters = "ResourceFilters"   // needs to be map[string]map[string]interface{}
const EndpointsFilters = "EndpointsFilters" // needs to be map[string]map[string]interface{}

const TargetSystemNamespace = "TargetSystem.Namespace"
const TargetSystemDeploymentsToManage = "TargetSystem.DeploymentsToManage" // []string

const ConfigurationValidationTotalMemory = "ConfigurationValidation.TotalMemory"
const ConfigurationValidationTotalCpu = "ConfigurationValidation.TotalCPU"

const AutoConfigureUseCache = "AutoConfigure.UseCache" // bool
const AutoConfigureCacheDatabaseType = "AutoConfigure.CacheDatabaseType"
const AutoConfigureCacheS3Region = "AutoConfigure.CacheS3Region"
const AutoConfigureCacheS3Bucket = "AutoConfigure.CacheS3Bucket"
const AutoConfigureFSDirectory = "AutoConfigure.FSDirectory"

const ResultsDirectory = "Results.Directory" // TODO do we need this? shouldn't we move to s3?

const WaitTimesWaitAfterConfigIsDeployedSeconds = "WaitTimes.WaitAfterConfigIsDeployedSeconds"
const WaitTimesLoadTestDurationSeconds = "WaitTimes.LoadTestDurationSeconds"
const WaitTimesWaitAfterLoadGeneratorIsDoneSeconds = "WaitTimes.WaitAfterLoadGeneratorIsDoneSeconds"
const WaitTimesWaitAfterLoadGeneratorStartes = "WaitTimes.WaitAfterLoadGeneratorStartes"

const TestName = "testname"

const TargetSystemWorkload = "TargetSystem.Workload"
const TargetSystemWorkloadName = "TargetSystem.Workload.Name"
const TargetSystemWorkloadBody = "TargetSystem.Workload.Body"

const StrategyName = "StrategyName"
const SLAConditionsFile = "SLAConditionsFile"

const LoadGeneratorArgs = "LoadGenerator.Args"
const LoadGeneratorType = "LoadGenerator.Type"
const LoadGeneratorScriptPath = "LoadGenerator.ScriptPath"
const LoadGeneratorCommand = "LoadGenerator.Command"

const AutoScalingApproachName = "AutoScaling.ApproachName"

const KubeConfigPath = "KubeConfigPath"

const DumpStartTime = "DumpStartTime"
const DumpFinishTime = "DumpFinishTime"
const DumpOutputPath = "DumpOutputPath"
const DumpWithCPUInfo = "DumpWithCPUInfo"

const InfluxDBCustomCRTFile = "InfluxDBCustomCRTFile"
