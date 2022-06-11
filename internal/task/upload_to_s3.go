package task

type UploadToS3Config struct {
	Filename string
	Bucket string
	NameInBucket string
}

type UploadToS3Task struct {
	config UploadToS3Config
}

func NewUploadToS3Task(config UploadToS3Config) *UploadToS3Task {
	return &UploadToS3Task{config: config}
}

func (ut UploadToS3Task) Name() string {
	return ""
}

func (ut *UploadToS3Task) Do() <-chan Status {
	return nil
}

