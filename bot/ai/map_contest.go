package ai

import (
	"fmt"
	hal "../core"
)

type ContestMap struct {
	Values			[][]int
}

// Strongly negative numbers are heavily in our area of influence.
// Strongly positive numbers are heavily in enemy area of influence.

func NewContestMap(game *hal.Game) *ContestMap {
	o := new(ContestMap)
	o.Values = hal.Make2dIntArray(game.Width(), game.Height())
	return o
}

func (self *ContestMap) Update(game *hal.Game, a *DistMap, b *EnemyDistMap) {

	width := game.Width()
	height := game.Height()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			self.Values[x][y] = a.Values[x][y] - b.Values[x][y]
		}
	}
}

func (self *ContestMap) Flog(game *hal.Game) {
	for x := 0; x < game.Width(); x++ {
		for y := 0; y < game.Height(); y++ {
			s := fmt.Sprintf("Contest: %v", self.Values[x][y])
			game.Flog(x, y, s, "")
		}
	}
}
