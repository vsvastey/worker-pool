package task

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

/* TODO: consider using https://github.com/spf13/afero
to abstract a FileSystem
and have an ability to unit test CopyFileTask
*/

// CopyFileTaskConfig is a configuration of CopyFileTask
type CopyFileTaskConfig struct {
	// Source is a name of a file to copy from
	Source string `yaml:"source"`
	// Destination is a name of a file to copy to
	Destination string `yaml:"destination"`
}

// CopyFileTask copies src file to dst
// src file must exist
// dst path must exist
type CopyFileTask struct {
	// src is a source filename
	src string
	// dst is a destination filename
	dst string
	// name is a description of the task
	name string
}

// NewCopyFileTask is a constructor of CopyFileTask
func NewCopyFileTask(config *CopyFileTaskConfig) (*CopyFileTask, error) {
	return &CopyFileTask{
		src:  config.Source,
		dst:  config.Destination,
		name: fmt.Sprintf("copy %s", filepath.Base(config.Source)),
	}, nil
}

func (cf CopyFileTask) Name() string {
	return cf.name
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
		if err != nil {
			return
		}
		dst, err := os.Create(cf.dst)
		if err != nil {
			return
		}
		dstWithStatus := NewWriterWithStatus(dst, info.Size(), res)

		_, err = io.Copy(dstWithStatus, src)
	}()
	return res
}
