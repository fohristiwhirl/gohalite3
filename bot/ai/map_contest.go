package ai

import (
	"fmt"
	hal "../core"
)

type ContestMap struct {
	Game			*hal.Game
	Values			[][]int
}

// Strongly negative numbers are heavily in our area of influence.
// Strongly positive numbers are heavily in enemy area of influence.

func NewContestMap(game *hal.Game) *ContestMap {
	o := new(ContestMap)
	o.Game = game
	o.Values = hal.Make2dIntArray(game.Width(), game.Height())
	return o
}

func (self *ContestMap) Update(a *DistMap, b *EnemyDistMap) {

	width := self.Game.Width()
	height := self.Game.Height()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			self.Values[x][y] = a.Values[x][y] - b.Values[x][y]
		}
	}
}

func (self *ContestMap) Flog() {
	for x := 0; x < self.Game.Width(); x++ {
		for y := 0; y < self.Game.Height(); y++ {
			s := fmt.Sprintf("Contest: %v", self.Values[x][y])
			self.Game.Flog(x, y, s, "")
		}
	}
}
