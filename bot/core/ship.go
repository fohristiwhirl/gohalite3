package core

import (
	"fmt"
)

type Ship struct {
	Game						*Game

	Owner						int			// Player ID
	Sid							int			// Ship ID
	X							int
	Y							int
	Halite						int
	Inspired					bool

	Alive						bool		// Parser updates this so stale refs can be detected by the AI
	Command						string		// AI's chosen command this turn
}

func (self *Ship) HaliteAt() int {
	return self.Game.HaliteAt(self.X, self.Y)
}

func (self *Ship) NeighbourPoints() []Point {
	return self.Game.NeighbourPoints(self.X, self.Y)
}

func (self *Ship) OnDropoff() bool {

	// True iff ship is on one of its own dropoffs or its factory.

	dropoffs := self.Game.Dropoffs(self.Owner)

	for _, d := range dropoffs {
		if self.X == d.X && self.Y == d.Y {
			return true
		}
	}

	return false
}

func (self *Ship) LocationFromMove(s string) (int, int) {

	switch s {

	case "w":
		return mod(self.X - 1, self.Game.width), self.Y
	case "e":
		return mod(self.X + 1, self.Game.width), self.Y
	case "n":
		return self.X,                           mod(self.Y - 1, self.Game.height)
	case "s":
		return self.X,                           mod(self.Y + 1, self.Game.height)
	default:
		return self.X,                           self.Y
	}
}

func (self *Ship) Left() {
	self.Command = fmt.Sprintf("m %d w", self.Sid)
}

func (self *Ship) Right() {
	self.Command = fmt.Sprintf("m %d e", self.Sid)
}

func (self *Ship) Up() {
	self.Command = fmt.Sprintf("m %d n", self.Sid)
}

func (self *Ship) Down() {
	self.Command = fmt.Sprintf("m %d s", self.Sid)
}

func (self *Ship) Construct() {
	self.Command = fmt.Sprintf("c %d", self.Sid)
}

func (self *Ship) Clear() {
	self.Command = ""
}

func (self *Ship) Set(s string) {
	self.Command = s
}
