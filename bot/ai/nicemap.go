package ai

import (
	"strconv"
	"strings"

	hal "../core"
)

const (
	NICE_RADIUS = 4
)

type NiceMap struct {
	Game			*hal.Game
	Values			[][]int
}

func NewNiceMap(game *hal.Game) *NiceMap {
	o := new(NiceMap)
	o.Game = game
	o.Values = make([][]int, game.Width())
	for x := 0; x < game.Width(); x++ {
		o.Values[x] = make([]int, game.Height())
	}
	o.Init()
	return o
}

func (self *NiceMap) Init() {
	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			self.Propagate(hal.Point{x, y}, self.Game.BoxAtFast(x, y).Halite, NICE_RADIUS)
		}
	}
}

func (self *NiceMap) Update() {
	all_changed := self.Game.ChangedBoxes()
	for _, box := range all_changed {
		self.Propagate(box, box.Delta, NICE_RADIUS)
	}
}

func (self *NiceMap) Propagate(origin hal.XYer, value int, radius int) {

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

func (self *NiceMap) Log() {

	for y := 0; y < len(self.Values[0]); y++ {

		var parts []string

		for x := 0; x < len(self.Values); x++ {
			parts = append(parts, strconv.Itoa(self.Values[x][y]))
		}

		line := strings.Join(parts, " ")

		self.Game.Log(line)
	}
}
