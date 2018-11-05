package ai

import (
	"fmt"
	hal "../core"
)

const (
	WEALTH_MAP_RADIUS = 4
)

type WealthMap struct {
	Values			[][]int
}

func NewWealthMap(game *hal.Game) *WealthMap {

	o := new(WealthMap)
	o.Values = hal.Make2dIntArray(game.Width(), game.Height())

	o.Init(game)		// Unlike some other maps, this one needs inited.
	return o
}

func (self *WealthMap) Init(game *hal.Game) {

	// Assumes the map is zeroed.
	// Can't be used as a way to update.

	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			self.Propagate(x, y, game.HaliteAtFast(x, y), WEALTH_MAP_RADIUS)
		}
	}
}

func (self *WealthMap) Update(game *hal.Game) {
	all_changed := game.Changes()
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

func (self *WealthMap) Flog(game *hal.Game) {
	for x := 0; x < game.Width(); x++ {
		for y := 0; y < game.Height(); y++ {
			s := fmt.Sprintf("Wealth: %v", self.Values[x][y])
			game.Flog(x, y, s, "")
		}
	}
}
