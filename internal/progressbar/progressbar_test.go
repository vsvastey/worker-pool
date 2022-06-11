package progressbar

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProgressBarEmpty(t *testing.T) {
	pb := NewProgressBar()
	res := pb.Draw()
	assert.Equal(t, " [                    ]  0/100", res)
}

func TestProgressBarUpdate(t *testing.T) {
	pb := NewProgressBar()

	type params struct {
		caption string
		progress int
	}
	testCases := map[params]string{
		 {"pb", 1}: "pb [                    ]  1/100",
		{"idle", 5}: "idle [*                   ]  5/100",
		{"pb", 10}: "pb [**                  ] 10/100",
		{"randomstr", 12}: "randomstr [**                  ] 12/100",
		{"pb", 50}: "pb [**********          ] 50/100",
		{"pb", 99}: "pb [******************* ] 99/100",
		{"pb", 100}: "pb [********************] 100/100",
	}

	for val, res := range testCases {
		pb.Set(val.caption, val.progress)
		assert.Equal(t, res, pb.Draw())
	}
}