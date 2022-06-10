package progressbar

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProgressBarEmpty(t *testing.T) {
	pb := NewProgressBar("pb")
	res := pb.Draw()
	assert.Equal(t, "pb [                    ]  0/100", res)
}

func TestProgressBarUpdate(t *testing.T) {
	pb := NewProgressBar("pb")

	testCases := map[int]string{
		 1: "pb [                    ]  1/100",
		 5: "pb [*                   ]  5/100",
		10: "pb [**                  ] 10/100",
		12: "pb [**                  ] 12/100",
		50: "pb [**********          ] 50/100",
		99: "pb [******************* ] 99/100",
		100: "pb [********************] 100/100",
	}

	for val, res := range testCases {
		pb.Set(val)
		assert.Equal(t, res, pb.Draw())
	}
}