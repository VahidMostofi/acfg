package configuration

import (
	"strconv"
)

type Configuration struct{
	ResourceType string
	ReplicaCount *int64
	CPU *int64
	Memory *int64
	EnvironmentValues map[string]string
}

func (c *Configuration) DeepCopy() *Configuration{
	c2 := &Configuration{
		ResourceType: c.ResourceType,
		ReplicaCount: c.ReplicaCount,
		CPU: c.CPU,
		Memory: c.Memory,
		EnvironmentValues: make(map[string]string),
	}
	for key,value := range c.EnvironmentValues{
		c2.EnvironmentValues[key] = value
	}

	return c
}

func (c *Configuration) GetCPUStringForK8s() string{
	s := strconv.FormatInt(*c.CPU, 10) + "m"
	return s
}

func (c *Configuration) GetMemoryStringForK8s() string{
	s := strconv.FormatInt(*c.Memory, 10) + "Mi"
	return s
}