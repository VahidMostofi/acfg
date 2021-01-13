package autoconfigurer

import (
	"github.com/vahidmostofi/acfg/internal/autocfg"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type AutoConfigurationAgent interface{
	GetName() string
	GetInitialConfiguration(workload *workload.Workload) (*autocfg.Configuration, error)
}
