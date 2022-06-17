package task

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type CopyFileTaskConfig struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

type CopyFileTask struct {
	src string
	dst string
}

func NewCopyFileTask(config *CopyFileTaskConfig) (*CopyFileTask, error) {
	return &CopyFileTask{
		src: config.Source,
		dst: config.Destination,
	}, nil
}

func (cf CopyFileTask) Caption() string {
	return fmt.Sprintf("copy %s", filepath.Base(cf.src))
}

func (cf *CopyFileTask) Do() <-chan Status {
	res := make(chan Status)

	go func() {
		defer func() {
			res <- Status{Progress: 100}
			close(res)
		}()

		src, err := os.Open(cf.src)
		if err != nil {
			log.Errorf("Error opening file %s: %v", cf.src, err)
			return
		}
		info, err := src.Stat()
		if err != nil {
			log.Errorf("Error on getting file %s information: %v", cf.src, err)
			return
		}

		dst, err := os.Create(cf.dst)
		if err != nil {
			log.Errorf("Error creating file %s: %v", cf.dst, err)
			return
		}

		dstWithStatus := NewWriterWithStatus(dst, info.Size(), res)
		_, err = io.Copy(dstWithStatus, src)
		if err != nil {
			log.Errorf("Error copying file %s to file %s: %v", cf.src, cf.dst, err)
			return
		}
	}()
	return res
}
