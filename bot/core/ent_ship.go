package core

import (
	"fmt"
)

type Ship struct {

	// A short lived data structure, valid only for 1 turn.

	Game						*Game
	Owner						int			// Player ID
	Sid							int			// Ship ID
	X							int
	Y							int
	Halite						int
	Inspired					bool
	Command						string		// Note that "o" means chosen-no-move while "" means no-choice-yet
}

func (self *Ship) String() string {
	return fmt.Sprintf("Ship %v (%v,%v, owner %v, command \"%v\")", self.Sid, self.X, self.Y, self.Owner, self.Command)
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

func (self *Ship) MoveCost() int {
	if self.Inspired {
		return self.Game.HaliteAt(self) / self.Game.Constants.INSPIRED_MOVE_COST_RATIO
	} else {
		return self.Game.HaliteAt(self) / self.Game.Constants.MOVE_COST_RATIO
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

func (self *Ship) HaliteAt() int {
	return self.Game.HaliteAtFast(self.X, self.Y)
}

func (self *Ship) LocationAfterMove(s string) Point {

	dx, dy := StringToDxDy(s)

	x := self.X + dx
	y := self.Y + dy

	x = Mod(x, self.Game.Width())
	y = Mod(y, self.Game.Height())

	return Point{x, y}
}

func (self *Ship) Point() Point {
	return Point{self.X, self.Y}
}

func (self *Ship) GetGame() *Game { return self.Game }
func (self *Ship) GetX() int { return self.X }
func (self *Ship) GetY() int { return self.Y }
func (self *Ship) DxDy(other XYer) (int, int) { return DxDy(self, other) }
func (self *Ship) Dist(other XYer) int { return Dist(self, other) }
func (self *Ship) SamePlace(other XYer) bool { return SamePlace(self, other) }
