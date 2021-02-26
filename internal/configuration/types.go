package configuration

import (
	"fmt"
	"math"
	"strconv"
)

const ResourceTypeDeployment = "Deployment"

type Configuration struct {
	ResourceType      string
	ReplicaCount      *int64
	CPU               *int64
	Memory            *int64
	EnvironmentValues map[string]string
}

func (c *Configuration) String() string {
	return fmt.Sprintf("%s: replica: %d, cpu: %d, memory: %dmb, envs: %v", c.ResourceType, *c.ReplicaCount, *c.CPU, *c.Memory, c.EnvironmentValues)
}

func (c *Configuration) UpdateEqualWithNewCPUValue(newCPU int64, maxCPUPerReplica int64) {
	replicaCount := int64(math.Ceil(float64(newCPU) / float64(maxCPUPerReplica)))
	cpuPerReplica := newCPU / replicaCount
	c.ReplicaCount = &replicaCount
	c.CPU = &cpuPerReplica
}

func (c *Configuration) DeepCopy() *Configuration {
	c2 := &Configuration{
		ResourceType:      c.ResourceType,
		ReplicaCount:      c.ReplicaCount,
		CPU:               c.CPU,
		Memory:            c.Memory,
		EnvironmentValues: make(map[string]string),
	}
	for key, value := range c.EnvironmentValues {
		c2.EnvironmentValues[key] = value
	}

	return c2
}

func (c *Configuration) GetCPUStringForK8s() string {
	s := strconv.FormatInt(*c.CPU, 10) + "m"
	return s
}

func (c *Configuration) GetMemoryStringForK8s() string {
	s := strconv.FormatInt(*c.Memory, 10) + "Mi"
	return s
}
