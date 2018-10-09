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

func (self *Ship) String() string {
	return fmt.Sprintf("Ship %v (%v,%v - owner %v)", self.Sid, self.X, self.Y, self.Owner)
}

func (self *Ship) HaliteUnder() int {
	return self.Game.HaliteAt(self.X, self.Y)
}

func (self *Ship) NeighbourPoints() []Point {
	return self.Game.NeighbourPoints(self.X, self.Y)
}

func (self *Ship) MoveCost() int {
	if self.Inspired {
		return self.HaliteUnder() / self.Game.Constants.INSPIRED_MOVE_COST_RATIO
	} else {
		return self.HaliteUnder() / self.Game.Constants.MOVE_COST_RATIO
	}
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

func (self *Ship) DxDy(x, y int) (int, int) {
	return self.Game.DxDy(self.X, self.Y, x, y)
}

func (self *Ship) Dist(x, y int) int {
	return self.Game.Dist(self.X, self.Y, x, y)
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
