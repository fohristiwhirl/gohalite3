package core

import (
	"../logging"
)

type InspirationMap struct {
	Threshold		int
	Values			[][]int
}

func NewInspirationMap(frame *Frame) *InspirationMap {

	self := new(InspirationMap)

	self.Threshold = frame.Constants.INSPIRATION_SHIP_COUNT
	self.Values = Make2dIntArray(frame.Width(), frame.Height())

	width := frame.Width()
	height := frame.Height()

	for _, ship := range frame.EnemyShips() {

		for y := 0; y <= frame.Constants.INSPIRATION_RADIUS; y++ {

			startx := y - frame.Constants.INSPIRATION_RADIUS
			endx := frame.Constants.INSPIRATION_RADIUS - y

			for x := startx; x <= endx; x++ {

				ox := Mod(ship.X + x, width)
				oy := Mod(ship.Y + y, height)
				self.Values[ox][oy] += 1

				if y != 0 {
					oy = Mod(ship.Y - y, height)
					self.Values[ox][oy] += 1
				}
			}
		}
	}

	return self
}

func (self *InspirationMap) Check(pos XYer) bool {
	return self.Values[pos.GetX()][pos.GetY()] >= self.Threshold
}

func (self *InspirationMap) Flog(turn int) {
	for x := 0; x < len(self.Values); x++ {
		for y := 0; y < len(self.Values[0]); y++ {
			if self.Values[x][y] >= 2 {
				logging.Flog(turn, x, y, "", "darkslategray")
			}
		}
	}
}
