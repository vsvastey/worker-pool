package progressbar

import (
	"fmt"
	"strings"
	"sync"
)

type ProgressBar struct {
	cur int
	caption string
	mu sync.Mutex
}

func NewProgressBar(caption string) *ProgressBar {
	return &ProgressBar{
		caption: caption,
		mu: sync.Mutex{},
	}
}

func (pb *ProgressBar) Set(value int) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.cur = value
}

func (pb *ProgressBar) Draw() string {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	asterixCount := pb.cur / 5
	return fmt.Sprintf("%s [%-20s] %2d/100", pb.caption, strings.Repeat("*", asterixCount), pb.cur)
}