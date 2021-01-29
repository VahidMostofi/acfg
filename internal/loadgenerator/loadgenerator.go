package loadgenerator

import (
	"github.com/vahidmostofi/acfg/internal/workload"
)

const containerName = "kkkk6localautoconfig"

type LoadGenerator interface{
	Start(workload *workload.Workload) error
	Stop() error
	GetFeedback() (map[string]interface{}, error)
}

func prepareLoadGenerator(workload *workload.Workload, info map[string]interface{}) ([]byte, error){

	return nil, nil
}
