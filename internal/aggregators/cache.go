package aggregators

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ConfigDatabase interface{
	Store(code string, data *AggregatedData) error
	Retrieve(code string) (*AggregatedData, error) // if there is no config with this hash returns nil,false
}


// TODO make this more general, it should work with []byte and string
type AWSConfigurationDatabase struct{
	session *session.Session
	bucket string
	region string
	directoryName string
}

func NewAWSConfigurationDatabase(s3Region, s3Bucket string) (*AWSConfigurationDatabase, error){
	a := &AWSConfigurationDatabase{
		directoryName: "cache",
		bucket: s3Bucket,
		region: s3Region,
	}
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(s3Region),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating aws session")
	}

	a.session = sess

	return a, nil
}

func(a *AWSConfigurationDatabase) Store(code string, data *AggregatedData) error{
	buffer, err := yaml.Marshal(data)
	if err != nil{
		return errors.Wrap(err, "error while marshaling aggregatedData")
	}
	//size := int64(len(buffer))

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(a.session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(a.bucket),
		Key:                  aws.String(a.directoryName + "/" + code),
		Body:                 bytes.NewReader(buffer),
		//ACL:                  aws.String("private"),
		//ContentLength:        aws.Int64(size),
		//ContentType:          aws.String(http.DetectContentType(buffer)),
		//ContentDisposition:   aws.String("attachment"),
		//ServerSideEncryption: aws.String("AES256"),
	})
	return errors.Wrap(err, "error while saving file to aws s3")
}

func(a *AWSConfigurationDatabase) Retrieve(code string) (*AggregatedData, error){
	log.Debugf("retrieveing with %s", code)
	key := a.directoryName + "/" + code
	oo, err := s3.New(a.session).GetObject(&s3.GetObjectInput{
		Bucket:               aws.String(a.bucket),
		Key:                  aws.String(key),
	})
	if err != nil{
		switch err.(awserr.Error).Code() {
		case s3.ErrCodeNoSuchKey:
			log.Debugf("no file found with key: %s", code)
			return nil,nil
		}
		return nil, errors.Wrapf(err, "error while getting object with key %s", key)
	}
	defer oo.Body.Close()

	ag := &AggregatedData{}
	err = yaml.NewDecoder(oo.Body).Decode(ag)
	if err != nil{
		return nil, errors.Wrapf(err, "error while deocing object with key %s", key)
	}

	return ag, nil
}
