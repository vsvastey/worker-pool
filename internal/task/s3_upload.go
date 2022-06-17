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

// S3UploadTaskConfig is configuration of S3UploadTask
type S3UploadTaskConfig struct {
	// Caption of file to upload
	Filename string `yaml:"filename"`
	// Bucket to which the file should be uploaded
	Bucket string `yaml:"bucket"`
	// Caption that the file should have in the bucket
	NameInBucket string `yaml:"name_in_bucket"`
	// TODO: move AWS connection parameters to an dedicated structure
	// S3 connection parameter
	AccessKeyId string `yaml:"access_key_id"`
	// S3 connection parameter
	SecretAccessKey string `yaml:"secret_access_key"`
	// S3 connection parameter
	Region string `yaml:"region"`
	// S3 connection parameter
	Endpoint string `yaml:"endpoint"`
}

// S3UploadTask is a Task that uploads provided file to S3
type S3UploadTask struct {
	filename        string
	bucket          string
	nameInBucket    string
	accessKeyId     string
	secretAccessKey string
	region          string
	endpoint        string
}

// NewS3UploadTask is a constructor for S3UploadTask
func NewS3UploadTask(config *S3UploadTaskConfig) (*S3UploadTask, error) {
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

func (ut S3UploadTask) Caption() string {
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
