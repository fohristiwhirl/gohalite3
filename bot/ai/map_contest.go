package ai

import (
	"fmt"
	hal "../core"
)

type ContestMap struct {
	Overmind		*Overmind
	Values			[][]int
}

// Strongly negative numbers are heavily in our area of influence.
// Strongly positive numbers are heavily in enemy area of influence.

func NewContestMap(overmind *Overmind, frame *hal.Frame) *ContestMap {
	o := new(ContestMap)
	o.Overmind = overmind
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

func (self *ContestMap) Flog() {

	frame := self.Overmind.Frame

	for x := 0; x < frame.Width(); x++ {
		for y := 0; y < frame.Height(); y++ {
			s := fmt.Sprintf("Contest: %v", self.Values[x][y])
			frame.Flog(x, y, s, "")
		}
	}
}
