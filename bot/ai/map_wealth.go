package ai

import (
	"fmt"
	hal "../core"
)

type WealthMap struct {
	Overmind		*Overmind
	Values			[][]int
}

func NewWealthMap(overmind *Overmind, frame *hal.Frame) *WealthMap {
	o := new(WealthMap)
	o.Overmind = overmind
	o.Values = hal.Make2dIntArray(frame.Width(), frame.Height())
	o.Init(frame)		// Unlike some other maps, this one needs inited.
	return o
}

func (self *WealthMap) Init(frame *hal.Frame) {

	// Assumes the map is zeroed.
	// Can't be used as a way to update.

	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			self.Propagate(x, y, frame.HaliteAtFast(x, y), WEALTH_MAP_RADIUS)
		}
	}
}

func (self *WealthMap) Update() {
	all_changed := self.Overmind.Frame.Changes()
	for _, c := range all_changed {
		self.Propagate(c.X, c.Y, c.Delta, WEALTH_MAP_RADIUS)
	}
}

func (self *WealthMap) Propagate(ox, oy, value int, radius int) {

	width := len(self.Values)
	height := len(self.Values[0])

	for y := 0; y <= radius; y++ {

		startx := y - radius
		endx := radius - y

		for x := startx; x <= endx; x++ {

			loc_x := hal.Mod(ox + x, width)
			loc_y := hal.Mod(oy + y, height)

			self.Values[loc_x][loc_y] += value

			if y != 0 {

				loc_y = hal.Mod(oy - y, height)

				self.Values[loc_x][loc_y] += value
			}
		}
	}
}

func (self *WealthMap) Flog() {

	frame := self.Overmind.Frame

	for x := 0; x < frame.Width(); x++ {
		for y := 0; y < frame.Height(); y++ {
			s := fmt.Sprintf("Wealth: %v", self.Values[x][y])
			frame.Flog(x, y, s, "")
		}
	}
}
