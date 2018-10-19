package core

type Box struct {
	Game						*Game
	X							int
	Y							int
	Halite						int
	Delta						int		// Change in Halite since last turn
}

func (self *Box) GetGame() *Game { return self.Game }
func (self *Box) GetX() int { return self.X }
func (self *Box) GetY() int { return self.Y }
func (self *Box) DxDy(other XYer) (int, int) { return DxDy(self, other) }
func (self *Box) Dist(other XYer) int { return Dist(self, other) }
func (self *Box) SamePlace(other XYer) bool { return SamePlace(self, other) }

// ------------------------------------------------------------

type Dropoff struct {
	Game						*Game
	Factory						bool
	Owner						int
	X							int
	Y							int
}

func (self *Dropoff) Box() *Box {
	return self.Game.BoxAt(self)
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
