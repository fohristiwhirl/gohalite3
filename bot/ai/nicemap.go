package ai

import (
	"strconv"
	"strings"

	hal "../core"
)

type NiceMap struct {
	Values			[][]int
}

func NewNiceMap(width, height int) *NiceMap {
	o := new(NiceMap)
	o.Values = make([][]int, width)
	for x := 0; x < width; x++ {
		o.Values[x] = make([]int, height)
	}
	return o
}

func (self *NiceMap) Init(game *hal.Game) {
	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			self.Propagate(game, hal.Point{x, y}, game.BoxAtFast(x, y).Halite, 4)		// FIXME: 0 is a test
		}
	}
}

func (self *NiceMap) Propagate(game *hal.Game, origin hal.XYer, value int, radius int) {

	width := len(self.Values)
	height := len(self.Values[0])

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

func (self *NiceMap) Log(game *hal.Game) {

	for y := 0; y < len(self.Values[0]); y++ {

		var parts []string

		for x := 0; x < len(self.Values); x++ {
			parts = append(parts, strconv.Itoa(self.Values[x][y]))
		}

		line := strings.Join(parts, " ")

		game.Log(line)
	}
}
