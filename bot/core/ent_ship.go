package core

import (
	"fmt"
)

type Ship struct {

	// A short lived data structure, valid only for 1 (possibly simulated) turn.

	Frame						*Frame
	Owner						int			// Player ID
	Sid							int			// Ship ID
	X							int
	Y							int
	Halite						int
	Inspired					bool
	Command						string		// Note that "o" means chosen-no-move while "" means no-choice-yet

	// Some stuff used by the AI...
	// The comms parser or simulated frame maker will have to save all of this from the previous turn's ship...

	Target						Point
	TargetOK					bool		// Has the target been validly set?
	Score						float32		// Score if our target is a mineable box.
	Desires						[]string
	Returning					bool
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
	player_dropoffs := self.Frame.Dropoffs(self.Owner)
	for _, dropoff := range player_dropoffs {
		if dropoff.X == self.X && dropoff.Y == self.Y {
			return true
		}
	}
	return false
}

func (self *Ship) MoveCost() int {
	if self.Inspired {
		return self.Frame.HaliteAt(self) / self.Frame.Constants.INSPIRED_MOVE_COST_RATIO
	} else {
		return self.Frame.HaliteAt(self) / self.Frame.Constants.MOVE_COST_RATIO
	}
}

func (self *Ship) NearestDropoff() *Dropoff {

	possibles := self.Frame.Dropoffs(self.Owner)

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
	return self.Frame.ShipCanDropoffAt(self, pos)
}

func (self *Ship) TargetIsDropoff() bool {
	return self.CanDropoffAt(self.Target)
}

func (self *Ship) TargetHalite() int {
	return self.Frame.HaliteAt(self.Target)
}

func (self *Ship) HaliteAt() int {
	return self.Frame.HaliteAtFast(self.X, self.Y)
}

func (self *Ship) LocationAfterMove(s string) Point {

	dx, dy := StringToDxDy(s)

	x := self.X + dx
	y := self.Y + dy

	x = Mod(x, self.Frame.Width())
	y = Mod(y, self.Frame.Height())

	return Point{x, y}
}

func (self *Ship) Point() Point {
	return Point{self.X, self.Y}
}

func (self *Ship) GetFrame() *Frame { return self.Frame }
func (self *Ship) GetX() int { return self.X }
func (self *Ship) GetY() int { return self.Y }
func (self *Ship) DxDy(other XYer) (int, int) { return DxDy(self, other) }
func (self *Ship) Dist(other XYer) int { return Dist(self, other) }
func (self *Ship) SamePlace(other XYer) bool { return SamePlace(self, other) }
