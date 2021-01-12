package sysstructureagg

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/constants"
)

type SystemStructure struct{
	endpoints2Resources map[string][]string
}

// NewSystemStructure ...
// kind: could be "predefined"; it will gets it from CONFIG_ENDPOINTS_2_RESOURCES
func NewSystemStructure(kind string) (*SystemStructure,error){

	s := &SystemStructure{}

	if kind == "predefined"{
		temp := viper.Get(constants.CONFIG_ENDPOINTS_2_RESOURCES)

		tempConverted, ok := temp.(map[string][]string)
		if !ok {
			return nil, errors.Errorf("cant find endpoints to resources in configs using: %s with type map[string]map[string]interface{}", constants.CONFIG_ENDPOINTS_2_RESOURCES)
		}
		s.endpoints2Resources = tempConverted
	}

	return s, nil
}

// GetEndpoints2Resources ...
func (s *SystemStructure) GetEndpoints2Resources() map[string][]string {
	return s.endpoints2Resources
}