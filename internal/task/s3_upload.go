package task

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

type S3UploadConfig struct {
	Filename     string `yaml:"filename"`
	Bucket       string `yaml:"bucket"`
	NameInBucket string `yaml:"name_in_bucket"`
	// TODO: move AWS connection parameters to an dedicated structure
	AccessKeyId     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Region          string `yaml:"region"`
	Endpoint        string `yaml:"endpoint"`
}

type S3UploadTask struct {
	filename        string
	bucket          string
	nameInBucket    string
	accessKeyId     string
	secretAccessKey string
	region          string
	endpoint        string
}

func NewS3Upload(config *S3UploadConfig) (*S3UploadTask, error) {
	return &S3UploadTask{
		filename:        config.Filename,
		bucket:          config.Bucket,
		nameInBucket:    config.NameInBucket,
		accessKeyId:     config.AccessKeyId,
		secretAccessKey: config.SecretAccessKey,
		region:          config.Region,
		endpoint:        config.Endpoint,
	}, nil
}

func (ut S3UploadTask) Name() string {
	return fmt.Sprintf("s3 %s", filepath.Base(ut.filename))
}

func (ut *S3UploadTask) Do() <-chan Status {
	res := make(chan Status)

	go func() {
		defer func() {
			res <- Status{Progress: 100}
			close(res)
		}()

		// Configure to use MinIO Server
		awsConfig := &aws.Config{
			Credentials:      credentials.NewStaticCredentials(ut.accessKeyId, ut.secretAccessKey, ""),
			Endpoint:         aws.String(ut.endpoint),
			Region:           aws.String(ut.region),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		}
		sess, err := session.NewSession(awsConfig)
		if err != nil {
			log.Errorf("Error creating AWS session: %v", err)
			return
		}

		src, err := os.Open(ut.filename)
		if err != nil {
			log.Errorf("Error opening file %s: %v", ut.filename, err)
			return
		}
		info, err := src.Stat()
		if err != nil {
			log.Errorf("Error getting file %s info: %v", ut.filename, err)
			return
		}
		reader := NewReaderWithStatus(src, info.Size(), res)

		uploader := s3manager.NewUploader(sess)
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: &ut.bucket,
			Key:    &ut.nameInBucket,
			Body:   reader,
		})
		if err != nil {
			log.Errorf("Error uploading file %s to S3: %v", ut.filename, err)
		}
	}()
	return res
}
