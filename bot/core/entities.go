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
	return fmt.Sprintf("Ship %v (%v,%v, owner %v)", self.Sid, self.X, self.Y, self.Owner)
}

func (self *Ship) ClearMove() {
	self.Command = ""
}

func (self *Ship) Move(s string) {		// Note that one cannot send "" - use ClearMove() instead

	if s == "e" || s == "w" || s == "s" || s == "n" || s == "c" || s == "o" {
		self.Command = s
	} else {
		panic(fmt.Sprintf("ship.Move() - Illegal move \"%s\" on %v", s, self))
	}
}

func (self *Ship) OnDropoff() bool {
	player_dropoffs := self.Game.Dropoffs(self.Owner)
	for _, dropoff := range player_dropoffs {
		if dropoff.X == self.X && dropoff.Y == self.Y {
			return true
		}
	}
	return false
}

func (self *Ship) Box() *Box {
	return self.Game.BoxAt(self)
}

func (self *Ship) MoveCost() int {
	if self.Inspired {
		return self.Box().Halite / self.Game.Constants.INSPIRED_MOVE_COST_RATIO
	} else {
		return self.Box().Halite / self.Game.Constants.MOVE_COST_RATIO
	}
}

func (self *Ship) NearestDropoff() *Dropoff {

	possibles := self.Game.Dropoffs(self.Owner)

	choice := possibles[0]
	choice_dist := self.Dist(choice)

	for _, dropoff := range possibles[1:] {
		dist := self.Dist(dropoff)
		if dist < choice_dist {
			choice = dropoff
			choice_dist = dist
		}
	}

	return choice
}

func (self *Ship) CanDropoffAt(pos XYer) bool {
	return self.Game.ShipCanDropoffAt(self, pos)
}

// ------------------------------------------------------------

type Dropoff struct {
	Game						*Game
	Factory						bool
	Owner						int			// Player ID
	X							int
	Y							int
}

func (self *Dropoff) Box() *Box {
	return self.Game.BoxAt(self)
}

type Box struct {
	Game						*Game
	X							int
	Y							int
	Halite						int
}

type Point struct {
	X							int
	Y							int
}

func (self Point) GetX() int { return self.X }
func (self Point) GetY() int { return self.Y }

// ------------------------------------------------------------

func (self *Box) GetGame() *Game { return self.Game }
func (self *Ship) GetGame() *Game { return self.Game }
func (self *Dropoff) GetGame() *Game { return self.Game }

func (self *Box) GetX() int { return self.X }
func (self *Ship) GetX() int { return self.X }
func (self *Dropoff) GetX() int { return self.X }

func (self *Box) GetY() int { return self.Y }
func (self *Ship) GetY() int { return self.Y }
func (self *Dropoff) GetY() int { return self.Y }

func (self *Box) DxDy(other XYer) (int, int) { return DxDy(self, other) }
func (self *Ship) DxDy(other XYer) (int, int) { return DxDy(self, other) }
func (self *Dropoff) DxDy(other XYer) (int, int) { return DxDy(self, other) }

func (self *Box) Dist(other XYer) int { return Dist(self, other) }
func (self *Ship) Dist(other XYer) int { return Dist(self, other) }
func (self *Dropoff) Dist(other XYer) int { return Dist(self, other) }

func (self *Box) SamePlace(other XYer) bool { return SamePlace(self, other) }
func (self *Ship) SamePlace(other XYer) bool { return SamePlace(self, other) }
func (self *Dropoff) SamePlace(other XYer) bool { return SamePlace(self, other) }
