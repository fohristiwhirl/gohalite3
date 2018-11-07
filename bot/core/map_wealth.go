package core

import (
	"fmt"
)

const (
	WEALTH_MAP_RADIUS = 4
)

type WealthMap struct {
	Values			[][]int
}

func NewWealthMap(frame *Frame) *WealthMap {
	o := new(WealthMap)
	o.Init(frame)		// Unlike some other maps, this one needs inited.
	return o
}

func (self *WealthMap) Init(frame *Frame) {

	self.Values = Make2dIntArray(frame.Width(), frame.Height())

	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			self.Propagate(x, y, frame.HaliteAtFast(x, y), WEALTH_MAP_RADIUS)
		}
	}
}

/*
func (self *WealthMap) Update(frame *Frame) {
	all_changed := frame.Changes()
	for _, c := range all_changed {
		self.Propagate(c.X, c.Y, c.Delta, WEALTH_MAP_RADIUS)
	}
}
*/

func (self *WealthMap) Propagate(ox, oy, value int, radius int) {

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

func (self *WealthMap) Flog(frame *Frame) {
	for x := 0; x < frame.Width(); x++ {
		for y := 0; y < frame.Height(); y++ {
			s := fmt.Sprintf("Wealth: %v", self.Values[x][y])
			frame.Flog(x, y, s, "")
		}
	}
}
