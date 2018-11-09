package core

import (
	"fmt"
	"../logging"
)

const (
	WEALTH_MAP_RADIUS = 4
)

type WealthMap struct {
	Values			[][]int
}

func NewWealthMap(frame *Frame) *WealthMap {

	self := new(WealthMap)

	self.Values = Make2dIntArray(frame.Width(), frame.Height())

	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			self.propagate(x, y, frame.HaliteAtFast(x, y), WEALTH_MAP_RADIUS)
		}
	}

	return self
}

func (self *WealthMap) propagate(ox, oy, value int, radius int) {

	width := len(self.Values)
	height := len(self.Values[0])

	for y := 0; y <= radius; y++ {

		startx := y - radius
		endx := radius - y

		for x := startx; x <= endx; x++ {

			loc_x := Mod(ox + x, width)
			loc_y := Mod(oy + y, height)

			self.Values[loc_x][loc_y] += value

			if y != 0 {

				loc_y = Mod(oy - y, height)

				self.Values[loc_x][loc_y] += value
			}
		}
	}
}

func (self *WealthMap) Flog(turn int) {
	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			s := fmt.Sprintf("Wealth: %v", self.Values[x][y])
			logging.Flog(turn, x, y, s, "")
		}
	}
}
