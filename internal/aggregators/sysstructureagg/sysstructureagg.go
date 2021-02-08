package sysstructureagg

type SystemStructure struct {
	Endpoints2Resources map[string][]string
}

// NewSystemStructure ...
// kind: could be "predefined"; it will gets it from
func NewSystemStructure(kind string, e2r map[string][]string) (*SystemStructure, error) {

	s := &SystemStructure{}

	if kind == "predefined" {
		s.Endpoints2Resources = e2r
	}

	return s, nil
}

// GetEndpoints2Resources ...
func (s *SystemStructure) GetEndpoints2Resources() map[string][]string {
	return s.Endpoints2Resources
}

// GetResources2Endpoints ...
func (s *SystemStructure) GetResources2Endpoints() map[string][]string {
	res := make(map[string][]string)
	for _, resources := range s.Endpoints2Resources {
		for _, r := range resources {
			res[r] = make([]string, 0)
		}
	}

	for endpoint, resources := range s.Endpoints2Resources {
		for _, r := range resources {
			res[r] = append(res[r], endpoint)
		}
	}

	return res
}
