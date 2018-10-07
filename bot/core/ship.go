package core

import (
	"fmt"
)

type Ship struct {
	game						*Game

	X							int
	Y							int
	Owner						int			// Player ID
	Id							int			// Ship ID
	Halite						int

	Command						string
}

func (self *Ship) HaliteAt() int {
	return self.game.HaliteAt(self.X, self.Y)
}

func (self *Ship) NeighbourPoints() []Point {
	return self.game.NeighbourPoints(self.X, self.Y)
}

func (self *Ship) LocationFromMove(s string) (int, int) {

	switch s {

	case "w":
		return mod(self.X - 1, self.game.width), self.Y
	case "e":
		return mod(self.X + 1, self.game.width), self.Y
	case "n":
		return self.X,                           mod(self.Y - 1, self.game.height)
	case "s":
		return self.X,                           mod(self.Y + 1, self.game.height)
	default:
		return self.X,                           self.Y
	}
}

func (self *Ship) Left() {
	self.Command = fmt.Sprintf("m %d w", self.Id)
}

func (self *Ship) Right() {
	self.Command = fmt.Sprintf("m %d e", self.Id)
}

func (self *Ship) Up() {
	self.Command = fmt.Sprintf("m %d n", self.Id)
}

func (self *Ship) Down() {
	self.Command = fmt.Sprintf("m %d s", self.Id)
}

func (self *Ship) Clear() {
	self.Command = ""
}

func (self *Ship) Construct() {
	self.Command = fmt.Sprintf("c %d", self.Id)
}
