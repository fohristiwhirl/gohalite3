package core

import (
	"fmt"
)

func (self *Frame) Pid() int { return self.pid }					// Note that simulated bots will be changing this
func (self *Frame) Turn() int { return self.turn }
func (self *Frame) Width() int { return self.width }
func (self *Frame) Height() int { return self.height }
func (self *Frame) Players() int { return self.players }

func (self *Frame) HaliteAt(pos XYer) int {
	x := Mod(pos.GetX(), self.width)
	y := Mod(pos.GetY(), self.height)
	return self.halite[x][y]
}

func (self *Frame) HaliteAtFast(x, y int) int {
	return self.halite[x][y]
}

func (self *Frame) ShipAt(pos XYer) *Ship {			// Maybe nil
	x := Mod(pos.GetX(), self.width)
	y := Mod(pos.GetY(), self.height)
	return self.ship_xy_lookup[Point{x, y}]
}

func (self *Frame) Sid(sid int) *Ship {				// Maybe nil
	return self.ship_id_lookup[sid]
}

func (self *Frame) Dropoffs(pid int) []*Dropoff {	// Includes factory

	var ret []*Dropoff

	for _, dropoff := range self.dropoffs {
		if dropoff.Owner == pid {
			ret = append(ret, dropoff)
		}
	}

	return ret
}

func (self *Frame) MyDropoffs() []*Dropoff {		// Includes factory
	return self.Dropoffs(self.pid)
}

func (self *Frame) AllDropoffs() []*Dropoff {
	return self.dropoffs
}

func (self *Frame) EnemyDropoffs() []*Dropoff {		// Includes factory

	var ret []*Dropoff

	for _, dropoff := range self.dropoffs {
		if dropoff.Owner != self.pid {
			ret = append(ret, dropoff)
		}
	}

	return ret
}

func (self *Frame) Ships(pid int) []*Ship {

	var ret []*Ship

	for _, ship := range self.ships {
		if ship.Owner == pid {
			ret = append(ret, ship)
		}
	}

	return ret
}

func (self *Frame) MyShips() []*Ship {
	return self.Ships(self.pid)
}

func (self *Frame) AllShips() []*Ship {
	return self.ships
}

func (self *Frame) EnemyShips() []*Ship {

	var ret []*Ship

	for _, ship := range self.ships {
		if ship.Owner != self.pid {
			ret = append(ret, ship)
		}
	}

	return ret
}

func (self *Frame) Budget(pid int) int {
	return self.budgets[pid]
}

func (self *Frame) MyBudget() int {
	return self.Budget(self.pid)
}

func (self *Frame) Factory(pid int) *Dropoff {

	factory := self.dropoffs[pid]

	// Factories are stored in the dropoff list in player order... but best check...

	if factory.Owner != pid || factory.Factory == false {
		panic(fmt.Sprintf("self.dropoffs[%d] wasn't the right factory", pid))
	}

	return factory
}

func (self *Frame) MyFactory() *Dropoff {
	return self.Factory(self.pid)
}

func (self *Frame) EnemyFactories() []*Dropoff {

	var ret []*Dropoff

	// Factories are stored in the dropoff list in player order...

	for n := 0; n < self.players; n++ {
		if n != self.pid {
			ret = append(ret, self.dropoffs[n])
		}
	}

	return ret
}

func (self *Frame) PlayerCanDropoffAt(pid int, pos XYer) bool {

	dropoffs := self.Dropoffs(pid)

	for _, dropoff := range dropoffs {
		if dropoff.X == pos.GetX() && dropoff.Y == pos.GetY() {
			return true
		}
	}
	return false
}

func (self *Frame) ShipCanDropoffAt(ship *Ship, pos XYer) bool {
	return self.PlayerCanDropoffAt(ship.Owner, pos)
}

func (self *Frame) WealthMap() *WealthMap {		// Return cached value if available.
	if self.wealth_map == nil {
		self.wealth_map = NewWealthMap(self)
	}
	return self.wealth_map
}

func (self *Frame) InspirationMap() *InspirationMap {
	if self.inspiration_map == nil {
		self.inspiration_map = make(map[int]*InspirationMap)
	}
	if self.inspiration_map[self.pid] == nil {
		self.inspiration_map[self.pid] = NewInspirationMap(self)
	}
	return self.inspiration_map[self.pid]
}

func (self *Frame) DropoffDistMap() *DropoffDistMap {
	if self.dropoff_dist_map == nil {
		self.dropoff_dist_map = make(map[int]*DropoffDistMap)
	}
	if self.dropoff_dist_map[self.pid] == nil {
		self.dropoff_dist_map[self.pid] = NewDropoffDistMap(self)
	}
	return self.dropoff_dist_map[self.pid]
}

func (self *Frame) InspirationCheck(pos XYer) bool {
	return self.InspirationMap().Check(pos)
}

func (self *Frame) InitialGroundHalite() int {
	return self.initial_ground_halite
}

func (self *Frame) GroundHalite() int {		// Return cached value if available. PreParse() relies on it also working without a valid value.

	if self.ground_halite > 0 {
		return self.ground_halite
	}

	for x := 0; x < self.width; x++ {
		for y := 0; y < self.height; y++ {
			self.ground_halite += self.halite[x][y]
		}
	}
	return self.ground_halite
}

func (self *Frame) AverageGroundHalite() int {
	return self.GroundHalite() / (self.width * self.height)
}

func (self *Frame) TotalShips() int {
	return len(self.ships)
}

func (self *Frame) Neighbours(x, y int) []Point {

	ret := make([]Point, 0, 4)

	x1, y1 := x + 1, y
	x2, y2 := x - 1, y
	x3, y3 := x, y + 1
	x4, y4 := x, y - 1

	x1 = Mod(x1, self.width)
	x2 = Mod(x2, self.width)
	x3 = Mod(x3, self.width)
	x4 = Mod(x4, self.width)

	y1 = Mod(y1, self.height)
	y2 = Mod(y2, self.height)
	y3 = Mod(y3, self.height)
	y4 = Mod(y4, self.height)

	ret = append(ret, Point{x1, y1})
	ret = append(ret, Point{x2, y2})
	ret = append(ret, Point{x3, y3})
	ret = append(ret, Point{x4, y4})

	return ret
}
