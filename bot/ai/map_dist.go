package ai

import (
	"fmt"
	hal "../core"
)

type DistMap struct {
	Game			*hal.Game
	Values			[][]int
}

func NewDistMap(game *hal.Game) *DistMap {
	o := new(DistMap)
	o.Game = game
	o.Values = hal.Make2dIntArray(game.Width(), game.Height())
	return o
}

func (self *DistMap) Update() {

	width := self.Game.Width()
	height := self.Game.Height()

	var hotpoints []hal.Point

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			self.Values[x][y] = 9999
		}
	}

	for _, ship := range self.Game.MyShips() {
		self.Values[ship.X][ship.Y] = 0
		hotpoints = append(hotpoints, hal.Point{ship.X, ship.Y})
	}

	factory := self.Game.MyFactory()
	self.Values[factory.X][factory.Y] = 0
	hotpoints = append(hotpoints, hal.Point{factory.X, factory.Y})

	for {

		var next_hotpoints []hal.Point

		for _, hotpoint := range hotpoints {

			neighbours := self.Game.Neighbours(hotpoint.X, hotpoint.Y)

			for _, box := range neighbours {

				if self.Values[box.X][box.Y] == 9999 {

					self.Values[box.X][box.Y] = self.Values[hotpoint.X][hotpoint.Y] + 1
					next_hotpoints = append(next_hotpoints, hal.Point{box.X, box.Y})
				}
			}
		}

		hotpoints = next_hotpoints

		if len(hotpoints) == 0 {
			return
		}
	}
}

func (self *DistMap) Flog() {
	for x := 0; x < self.Game.Width(); x++ {
		for y := 0; y < self.Game.Height(); y++ {
			s := fmt.Sprintf("Friendly dist: %v", self.Values[x][y])
			self.Game.Flog(x, y, s, "")
		}
	}
}
