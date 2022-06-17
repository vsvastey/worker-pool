package progressbar

import (
	"fmt"
	"strings"
	"sync"
)

type ProgressBar struct {
	cur     int
	caption string
	mu      sync.Mutex
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		mu: sync.Mutex{},
	}
}

func (pb *ProgressBar) Set(caption string, value int) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.caption = caption
	pb.cur = value
}

func (pb *ProgressBar) Draw() string {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	asterixCount := pb.cur / 5
	return fmt.Sprintf("%20s [%-20s] %3d/100", pb.caption, strings.Repeat("*", asterixCount), pb.cur)
}
