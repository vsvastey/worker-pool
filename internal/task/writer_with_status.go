package task

import (
	"io"
	"sync"
)

// WriterWithStatus is a Writer interface implementation that
// wraps an already existing Writer and reports to a statusChan channel
// how much data (in percent) has already been written
type WriterWithStatus struct {
	wp         io.Writer
	written    int64
	statusChan chan<- Status
	mu         sync.Mutex
	total      int64
}

func NewWriterWithStatus(wp io.Writer, size int64, ch chan<- Status) *WriterWithStatus {
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
