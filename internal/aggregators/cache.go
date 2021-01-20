package aggregators

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/acfg/internal/dataaccess"
	"gopkg.in/yaml.v2"
)

type ConfigDatabase interface{
	Store(code string, data *AggregatedData) error
	Retrieve(code string) (*AggregatedData, error) // if there is no config with this hash returns nil,false
}


// TODO make this more general, it should work with []byte and string
type AWSConfigurationDatabase struct{
	directoryName string
	s3db *dataaccess.S3Storage
}

func NewAWSConfigurationDatabase(s3Region, s3Bucket string) (*AWSConfigurationDatabase, error){
	a := &AWSConfigurationDatabase{directoryName: "cache"}
	s3, err := dataaccess.NewS3Storage(s3Region, s3Bucket)
	if err != nil {
		return nil, errors.Wrap(err, "error creating AWS s3 db.")
	}
	a.s3db = s3
	return a, nil
}

func(a *AWSConfigurationDatabase) Store(code string, data *AggregatedData) error{
	log.Debugf("ConfigCache: storing with %s", code)
	buffer, err := yaml.Marshal(data)
	if err != nil{
		return errors.Wrap(err, "error while marshaling aggregatedData")
	}
	key := a.directoryName + "/" + code + ".yaml"
	err = a.s3db.Store(key, buffer)
	return errors.Wrap(err, "error while saving file to AWS s3")
}

func(a *AWSConfigurationDatabase) Retrieve(code string) (*AggregatedData, error){
	log.Debugf("ConfigCache: retrieveing with %s", code)
	key := a.directoryName + "/" + code
	buffer, err := a.s3db.Retrieve(key)
	if buffer == nil{
		return nil, nil
	}
	ag := &AggregatedData{}
	err = yaml.Unmarshal(buffer, ag)
	return ag, errors.Wrapf(err, "error while deocing object with key %s", key)
}
