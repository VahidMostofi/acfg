package autoconfigurer

import (
	"github.com/vahidmostofi/acfg/internal/autocfg"
	"github.com/vahidmostofi/acfg/internal/workload"
)

type AutoConfigurationAgent interface{
	GetName() string
	GetInitialConfiguration(workload *workload.Workload) (map[string]*autocfg.Configuration, error)
	ConfigureNextStep(workload *workload.Workload) (map[string]*autocfg.Configuration, bool, error)
}
