package task

import (
	"os"
	"sync"
)

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
