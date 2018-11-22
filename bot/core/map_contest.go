package core

import (
	"fmt"
	"../logging"
)

type ContestMap struct {
	Values			[][]int
}

// Strongly negative numbers are heavily in our area of influence.
// Strongly positive numbers are heavily in enemy area of influence.

func NewContestMap(a *FriendlyDistMap, b *EnemyDistMap) *ContestMap {

	self := new(ContestMap)
	self.Values = Make2dIntArray(len(a.Values), len(a.Values[0]))

	for x := 0; x < len(a.Values); x++ {
		for y := 0; y < len(a.Values[0]); y++ {
			self.Values[x][y] = a.Values[x][y] - b.Values[x][y]
		}
	}

	return self
}

func (self *ContestMap) Flog(turn int) {
	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			s := fmt.Sprintf("Contest: %v", self.Values[x][y])
			logging.Flog(turn, x, y, s, "")
		}
	}
}
