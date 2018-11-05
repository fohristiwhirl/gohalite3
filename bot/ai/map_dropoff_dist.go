package ai

import (
	"fmt"
	hal "../core"
)

type DropoffDistMap struct {
	Values			[][]int
}

func NewDropoffDistMap(game *hal.Game) *DropoffDistMap {
	o := new(DropoffDistMap)
	o.Values = hal.Make2dIntArray(game.Width(), game.Height())
	return o
}

func (self *DropoffDistMap) Update(game *hal.Game) {

	width := game.Width()
	height := game.Height()

	var hotpoints []hal.Point

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			self.Values[x][y] = 9999
		}
	}

	for _, dropoff := range game.MyDropoffs() {
		self.Values[dropoff.X][dropoff.Y] = 0
		hotpoints = append(hotpoints, hal.Point{dropoff.X, dropoff.Y})
	}

	for {

		var next_hotpoints []hal.Point

		for _, hotpoint := range hotpoints {

			neighbours := game.Neighbours(hotpoint.X, hotpoint.Y)

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

func (self *DropoffDistMap) Flog(game *hal.Game) {
	for x := 0; x < game.Width(); x++ {
		for y := 0; y < game.Height(); y++ {
			s := fmt.Sprintf("Dropoff dist: %v", self.Values[x][y])
			game.Flog(x, y, s, "")
		}
	}
}
