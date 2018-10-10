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
	Command						string		// Note that "o" means chosen-no-move while "" means no-choice-yet
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

func (self *Ship) DxDy(x, y int) (int, int) {
	return self.Game.DxDy(self.X, self.Y, x, y)
}

func (self *Ship) Dist(x, y int) int {
	return self.Game.Dist(self.X, self.Y, x, y)
}

func (self *Ship) Move(s string) {

	// Note that one cannot send "" - use ClearMove() instead

	if s == "e" || s == "w" || s == "s" || s == "n" || s == "c" || s == "o" {
		self.Command = s
	} else {
		panic(fmt.Sprintf("ship.Move() - Illegal move \"%s\" on %v", s, self))
	}
}

func (self *Ship) ClearMove() {
	self.Command = ""
}
