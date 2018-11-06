package ai

import (
	"fmt"
	hal "../core"
)

type EnemyDistMap struct {
	Values			[][]int
}

func NewEnemyDistMap(frame *hal.Frame) *EnemyDistMap {
	o := new(EnemyDistMap)
	o.Values = hal.Make2dIntArray(frame.Width(), frame.Height())
	return o
}

func (self *EnemyDistMap) Update(frame *hal.Frame) {

	width := frame.Width()
	height := frame.Height()

	var hotpoints []hal.Point

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			self.Values[x][y] = 9999
		}
	}

	for _, ship := range frame.EnemyShips() {
		self.Values[ship.X][ship.Y] = 0
		hotpoints = append(hotpoints, hal.Point{ship.X, ship.Y})
	}

	for _, factory := range frame.EnemyFactories() {
		self.Values[factory.X][factory.Y] = 0
		hotpoints = append(hotpoints, hal.Point{factory.X, factory.Y})
	}

	for {

		var next_hotpoints []hal.Point

		for _, hotpoint := range hotpoints {

			neighbours := frame.Neighbours(hotpoint.X, hotpoint.Y)

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

func (self *EnemyDistMap) Flog(frame *hal.Frame) {
	for x := 0; x < frame.Width(); x++ {
		for y := 0; y < frame.Height(); y++ {
			s := fmt.Sprintf("Enemy dist: %v", self.Values[x][y])
			frame.Flog(x, y, s, "")
		}
	}
}
