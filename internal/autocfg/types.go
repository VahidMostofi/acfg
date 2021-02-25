package autocfg

import (
	"crypto/md5"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/acfg/internal/aggregators"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"github.com/vahidmostofi/acfg/internal/constants"
	"github.com/vahidmostofi/acfg/internal/workload"
	"k8s.io/apimachinery/pkg/util/json"
)

//TODO add more args to this hash function. probably make it work with ...
func GetHash(c map[string]*configuration.Configuration, version string, inputWorkload *workload.Workload) (string, error) {
	// panic("you need to add name and system and other stuff to this!")
	b, err := json.Marshal(c)
	if err != nil {
		return "", errors.Wrap(err, "cant convert configuration to json")
	}
	b = append(b, []byte(version)...)

	temp, err := json.Marshal(inputWorkload)
	if err != nil {
		panic(err)
	}
	b = append(b, temp...)

	b = append(b, []byte(viper.GetString(constants.TestName))...)

	log.Debugf("hashing with %s", string(b))
	s := md5.Sum(b)
	h := fmt.Sprintf("%x", s)
	return h, err
}

type IterationInformation struct {
	Configuration         map[string]*configuration.Configuration `yaml:"configurations"`
	StartTime             int64                                   `yaml:"startTime"`
	FinishTime            int64                                   `yaml:"finishTime"`
	AggregatedData        *aggregators.AggregatedData             `yaml:"aggregatedData"`
	LoadGeneratorFeedback map[string]interface{}                  `yaml:"loadGeneratorFeedback"`
}

type TestInformation struct {
	Name                    string                  `yaml:"name"`
	VersionCode             string                  `yaml:"version"`
	AutoconfiguringApproach string                  `yaml:"autoConfigApproach"`
	Iterations              []*IterationInformation `yaml:"iterations"`
	InputWorkload           *workload.Workload      `yaml:"workload"`
	AllSettings             map[string]interface{}  `yaml:"allSettings"`
}
