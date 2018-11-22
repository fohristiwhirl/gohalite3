package core

import (
	"fmt"
	"../logging"
)

type FriendlyDistMap struct {
	Values			[][]int
}

func NewFriendlyDistMap(frame *Frame) *FriendlyDistMap {

	width := frame.Width()
	height := frame.Height()

	self := new(FriendlyDistMap)
	self.Values = Make2dIntArray(width, height)

	var hotpoints []Point

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			self.Values[x][y] = 9999
		}
	}

	for _, ship := range frame.MyShips() {
		self.Values[ship.X][ship.Y] = 0
		hotpoints = append(hotpoints, Point{ship.X, ship.Y})
	}

	factory := frame.MyFactory()
	self.Values[factory.X][factory.Y] = 0
	hotpoints = append(hotpoints, Point{factory.X, factory.Y})

	for {

		var next_hotpoints []Point

		for _, hotpoint := range hotpoints {

			neighbours := frame.Neighbours(hotpoint.X, hotpoint.Y)

			for _, box := range neighbours {

				if self.Values[box.X][box.Y] == 9999 {

					self.Values[box.X][box.Y] = self.Values[hotpoint.X][hotpoint.Y] + 1
					next_hotpoints = append(next_hotpoints, Point{box.X, box.Y})
				}
			}
		}

		hotpoints = next_hotpoints

		if len(hotpoints) == 0 {
			return self
		}
	}
}

func (self *FriendlyDistMap) Flog(turn int) {
	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			s := fmt.Sprintf("Friendly dist: %v", self.Values[x][y])
			logging.Flog(turn, x, y, s, "")
		}
	}
}
