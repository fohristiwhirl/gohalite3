package core

import (
	"fmt"
)

type Dropoff struct {

	// A short lived data structure, valid only for 1 turn. (Well, it's not like they move, but...)

	Game						*Game
	Factory						bool
	Owner						int
	X							int
	Y							int
}

func (self *Dropoff) String() string {
	return fmt.Sprintf("Dropoff (%v,%v, owner %v)", self.X, self.Y, self.Owner)
}

func (self *Dropoff) Point() Point {
	return Point{self.X, self.Y}
}

func (self *Dropoff) GetGame() *Game { return self.Game }
func (self *Dropoff) GetX() int { return self.X }
func (self *Dropoff) GetY() int { return self.Y }
func (self *Dropoff) DxDy(other XYer) (int, int) { return DxDy(self, other) }
func (self *Dropoff) Dist(other XYer) int { return Dist(self, other) }
func (self *Dropoff) SamePlace(other XYer) bool { return SamePlace(self, other) }

// ------------------------------------------------------------

type Point struct {
	X							int
	Y							int
}

func (self Point) GetX() int { return self.X }
func (self Point) GetY() int { return self.Y }
