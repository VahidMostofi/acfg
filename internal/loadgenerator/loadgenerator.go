package loadgenerator

import (
	"github.com/vahidmostofi/acfg/internal/workload"
	"io"
)

const containerName = "kkkk6localautoconfig"

type LoadGenerator interface{
	Start(workload *workload.Workload, reader io.Reader, extras map[string]string) error
	Stop() error
	GetFeedback() (map[string]interface{}, error)
}

func prepareLoadGenerator(workload *workload.Workload, info map[string]interface{}) ([]byte, error){

	return nil, nil
}
