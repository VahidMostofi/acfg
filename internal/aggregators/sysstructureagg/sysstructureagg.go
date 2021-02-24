package sysstructureagg

type SystemStructure struct{
	Endpoints2Resources map[string][]string
}

// NewSystemStructure ...
// kind: could be "predefined"; it will gets it from
func NewSystemStructure(kind string, e2r map[string][]string) (*SystemStructure,error){

	s := &SystemStructure{}

	if kind == "predefined"{
		s.Endpoints2Resources = e2r
	}

	return s, nil
}

// GetEndpoints2Resources ...
func (s *SystemStructure) GetEndpoints2Resources() map[string][]string {
	return s.Endpoints2Resources
}