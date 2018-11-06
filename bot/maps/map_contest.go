package maps

import (
	"fmt"
	hal "../core"
)

type ContestMap struct {
	Values			[][]int
}

// Strongly negative numbers are heavily in our area of influence.
// Strongly positive numbers are heavily in enemy area of influence.

func NewContestMap(frame *hal.Frame) *ContestMap {
	o := new(ContestMap)
	o.Values = hal.Make2dIntArray(frame.Width(), frame.Height())
	return o
}

func (self *ContestMap) Update(a *DistMap, b *EnemyDistMap) {
	for x := 0; x < len(a.Values); x++ {
		for y := 0; y < len(a.Values[0]); y++ {
			self.Values[x][y] = a.Values[x][y] - b.Values[x][y]
		}
	}
}

func (self *ContestMap) Flog(frame *hal.Frame) {
	for x := 0; x < frame.Width(); x++ {
		for y := 0; y < frame.Height(); y++ {
			s := fmt.Sprintf("Contest: %v", self.Values[x][y])
			frame.Flog(x, y, s, "")
		}
	}
}
