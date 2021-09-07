# Homebase for MOAT and TRIM
Previous version of MOAT that works with Docker Swarm can be found [here](https://github.com/VahidMostofi/swarmmanager)
# ACFG Auto ConFiG
The new version of automatic configuration of microservice applications.

## Environment Variables
    
These environment variables work independent from the configuration file.
- `LOG_LEVEL` controls the log level. 

## Configurations
Any configuration item in the yaml file can be set using environment variables and in this case would overwrite the value of the config file.

The environment variables must start with `ACFG_` and each level of the config name would be separated with `_`. For example if you want to set the `URL` of `Args` of `EndpointAggregator` you should use `ACFG_ENDPOINTSAGGREGATOR_ARGS_URL`

### Configuring Aggregators
Aggregators are there to collect different resource metrics. Currently the only type is `influxdb`. See below on how to configure influxdb. **This implementation works with influxdb v2**.

These aggregators need to be configured:
  - **EndpointsAggregator**, information about response times and status of each request sent the each endpoint.
  - **ResourceUsageAggregator**, information about CPU and Memory usage. The memory is not implemented at this version.
  - **WorkloadAggregator**, information about the generated workload during one iteration. 

Beside these, the **SystemStructureAggregator** also need to be defined, the only type for this aggregator at this version is `predefined`. It specifies the path each endpoint (request-type) takes (the services involved in each request type).
For `predefined` type, it needs `Endpoints2Resources` which accepts a `map[string][]string`. See sample-configs.

### influxdb aggregator

It needs these values to work:

  - Args:
  
    - URL
    - Token
    - Bucket
    - Organization


### HowTos?

How to create the cluster?
```
k3d cluster create -p "9099:9099@loadbalancer"
```

=======
## For the latest implemenation of MOAT and TRIM, pleas refer to `develop` branch of the same repository.
