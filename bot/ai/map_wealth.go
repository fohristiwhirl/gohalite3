package ai

import (
	"fmt"
	hal "../core"
)

const (
	WEALTH_MAP_RADIUS = 4
)

type WealthMap struct {
	Game			*hal.Game
	Values			[][]int
}

func NewWealthMap(game *hal.Game) *WealthMap {

	o := new(WealthMap)
	o.Game = game
	o.Values = hal.Make2dIntArray(game.Width(), game.Height())

	o.Init()		// Unlike some other maps, this one needs inited.
	return o
}

func (self *WealthMap) Init() {

	// Assumes the map is zeroed.
	// Can't be used as a way to update.

	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			self.Propagate(x, y, self.Game.HaliteAtFast(x, y), WEALTH_MAP_RADIUS)
		}
	}
}

func (self *WealthMap) Update() {
	all_changed := self.Game.Changes()
	for _, c := range all_changed {
		self.Propagate(c.X, c.Y, c.Delta, WEALTH_MAP_RADIUS)
	}
}

func (self *WealthMap) Propagate(ox, oy, value int, radius int) {

	width := self.Game.Width()
	height := self.Game.Height()

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
	for x := 0; x < self.Game.Width(); x++ {
		for y := 0; y < self.Game.Height(); y++ {
			s := fmt.Sprintf("Wealth: %v", self.Values[x][y])
			self.Game.Flog(x, y, s, "")
		}
	}
}
