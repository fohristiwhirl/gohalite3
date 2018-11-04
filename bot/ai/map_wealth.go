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
	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			self.Propagate(hal.Point{x, y}, self.Game.BoxAtFast(x, y).Halite, WEALTH_MAP_RADIUS)
		}
	}
}

func (self *WealthMap) Update() {
	all_changed := self.Game.ChangedBoxes()
	for _, box := range all_changed {
		self.Propagate(box, box.Delta, WEALTH_MAP_RADIUS)
	}
}

func (self *WealthMap) Propagate(origin hal.XYer, value int, radius int) {

	width := self.Game.Width()
	height := self.Game.Height()

	ox := origin.GetX()
	oy := origin.GetY()

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
