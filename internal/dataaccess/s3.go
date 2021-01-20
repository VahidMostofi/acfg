package dataaccess

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type S3Storage struct{
	session *session.Session
	bucket string
	region string
}

func NewS3Storage (region, bucket string) (*S3Storage, error){
	s := &S3Storage{
		bucket: bucket,
		region: region,
	}
	session, err := session.NewSession(
		&aws.Config{
			Region: aws.String(s.region),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating aws session.")
	}
	s.session = session

	return s, nil
}

func (ss *S3Storage) Store(key string, data []byte) error{
	log.Debugf("S3.Store() for key: %s and size: %d", key, len(data))
	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err := s3.New(ss.session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(ss.bucket),
		Key:                  aws.String(key),
		Body:                 bytes.NewReader(data),
	})
	return errors.Wrap(err, "error while saving file to aws s3")
}

func (ss *S3Storage) Retrieve(key string) ([]byte, error){
	log.Debugf("S3.Retrieve() retrieveing with key: %s", key)

	oo, err := s3.New(ss.session).GetObject(&s3.GetObjectInput{
		Bucket:               aws.String(ss.bucket),
		Key:                  aws.String(key),
	})
	if err != nil{
		switch err.(awserr.Error).Code() {
		case s3.ErrCodeNoSuchKey:
			log.Debugf("not file found with key: %s", key)
			return nil,nil
		}
		return nil, errors.Wrapf(err, "error while getting object with key %s", key)
	}
	defer oo.Body.Close()

	res, err := ioutil.ReadAll(oo.Body)
	return res, errors.Wrapf(err, "error while reading buffer with key %s", key)
}