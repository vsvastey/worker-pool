package task

import (
	"io"
	"sync"
)

type ReaderWithStatus struct {
	reader     io.Reader
	read       int64
	statusChan chan<- Status
	mu         sync.Mutex
	total      int64
}

func NewReaderWithStatus(reader io.Reader, size int64, ch chan<- Status) *ReaderWithStatus {
	return &ReaderWithStatus{reader: reader, total: size, statusChan: ch}
}

func (rs *ReaderWithStatus) Read(p []byte) (n int, err error) {
	n, err = rs.reader.Read(p)
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.read += int64(n)
	rs.statusChan <- Status{Progress: int(100 * rs.read / rs.total)}
	return n, err
}
