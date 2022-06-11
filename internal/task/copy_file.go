package task

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type CopyFileConfig struct {
	Source string
	Destination string
}

type CopyFileTask struct {
	config CopyFileConfig
}

func NewCopyFileTask(config CopyFileConfig) *CopyFileTask {
	return &CopyFileTask{config: config}
}

func (cf CopyFileTask) Name() string {
	return fmt.Sprintf("copy %s", filepath.Base(cf.config.Source))
}

func (cf *CopyFileTask) Do() <-chan Status {
	res := make(chan Status)

	go func() {
		defer func() {
			res <- Status{Progress: 100}
			close(res)
		}()

		src, err := os.Open(cf.config.Source)
		if err != nil {
			return
		}
		info, err := src.Stat()
		dst, err := os.Create(cf.config.Destination)
		dstWithStatus := NewWriterWithStatus(dst, info.Size(), res)
		if err != nil {
			return
		}
		_, err = io.Copy(dstWithStatus, src)
	}()
	return res
}

type WriterWithStatus struct {
	wp *os.File
	written int64
	statusChan chan<- Status
	mu sync.Mutex
	total int64
}

func NewWriterWithStatus(wp *os.File, size int64, ch chan<- Status) *WriterWithStatus {
	return &WriterWithStatus{wp: wp, total: size, statusChan: ch}
}

func (ws *WriterWithStatus) Write(p []byte) (n int, err error) {
	n, err = ws.wp.Write(p)
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.written += int64(n)
	ws.statusChan <- Status{Progress: int(100 * ws.written / ws.total)}
	return n, err
}


