package task

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func (cf CopyFileTask) Name() string {
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
			return
		}
		info, err := src.Stat()
		dst, err := os.Create(cf.dst)
		dstWithStatus := NewWriterWithStatus(dst, info.Size(), res)
		if err != nil {
			return
		}
		_, err = io.Copy(dstWithStatus, src)
	}()
	return res
}
