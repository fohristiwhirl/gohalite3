package core

import (
	"fmt"
)

func (self *Game) Pid() int { return self.pid }					// Note that simulated bots will be changing this
func (self *Game) Turn() int { return self.turn }
func (self *Game) Width() int { return self.width }
func (self *Game) Height() int { return self.height }
func (self *Game) Players() int { return self.players }

func (self *Game) HaliteAt(pos XYer) int {
	x := Mod(pos.GetX(), self.width)
	y := Mod(pos.GetY(), self.height)
	return self.halite[x][y]
}

func (self *Game) HaliteAtFast(x, y int) int {
	return self.halite[x][y]
}

func (self *Game) ShipAt(pos XYer) *Ship {			// Maybe nil
	x := Mod(pos.GetX(), self.width)
	y := Mod(pos.GetY(), self.height)
	return self.ship_xy_lookup[Point{x, y}]
}

func (self *Game) Sid(sid int) *Ship {				// Maybe nil
	return self.ship_id_lookup[sid]
}

func (self *Game) Dropoffs(pid int) []*Dropoff {	// Includes factory

	var ret []*Dropoff

	for _, dropoff := range self.dropoffs {
		if dropoff.Owner == pid {
			ret = append(ret, dropoff)
		}
	}

	return ret
}

func (self *Game) MyDropoffs() []*Dropoff {			// Includes factory
	return self.Dropoffs(self.pid)
}

func (self *Game) AllDropoffs() []*Dropoff {
	return self.dropoffs
}

func (self *Game) EnemyDropoffs() []*Dropoff {		// Includes factory

	var ret []*Dropoff

	for _, dropoff := range self.dropoffs {
		if dropoff.Owner != self.pid {
			ret = append(ret, dropoff)
		}
	}

	return ret
}

func (self *Game) Ships(pid int) []*Ship {

	var ret []*Ship

	for _, ship := range self.ships {
		if ship.Owner == pid {
			ret = append(ret, ship)
		}
	}

	return ret
}

func (self *Game) MyShips() []*Ship {
	return self.Ships(self.pid)
}

func (self *Game) AllShips() []*Ship {
	return self.ships
}

func (self *Game) EnemyShips() []*Ship {

	var ret []*Ship

	for _, ship := range self.ships {
		if ship.Owner != self.pid {
			ret = append(ret, ship)
		}
	}

	return ret
}

func (self *Game) Budget(pid int) int {
	return self.budgets[pid]
}

func (self *Game) MyBudget() int {
	return self.Budget(self.pid)
}

func (self *Game) Factory(pid int) *Dropoff {

	factory := self.dropoffs[pid]

	// Factories are stored in the dropoff list in player order... but best check...

	if factory.Owner != pid || factory.Factory == false {
		panic(fmt.Sprintf("self.dropoffs[%d] wasn't the right factory", pid))
	}

	return factory
}

func (self *Game) MyFactory() *Dropoff {
	return self.Factory(self.pid)
}

func (self *Game) EnemyFactories() []*Dropoff {

	var ret []*Dropoff

	// Factories are stored in the dropoff list in player order...

	for n := 0; n < self.players; n++ {
		if n != self.pid {
			ret = append(ret, self.dropoffs[n])
		}
	}

	return ret
}

func (self *Game) PlayerCanDropoffAt(pid int, pos XYer) bool {

	dropoffs := self.Dropoffs(pid)

	for _, dropoff := range dropoffs {
		if dropoff.X == pos.GetX() && dropoff.Y == pos.GetY() {
			return true
		}
	}
	return false
}

func (self *Game) ShipCanDropoffAt(ship *Ship, pos XYer) bool {
	return self.PlayerCanDropoffAt(ship.Owner, pos)
}

func (self *Game) Hash() string {
	return self.hash
}

func (self *Game) GroundHalite() int {
	var count int
	for x := 0; x < self.width; x++ {
		for y := 0; y < self.height; y++ {
			count += self.halite[x][y]
		}
	}
	return count
}

func (self *Game) TotalShips() int {
	return len(self.ships)
}

func (self *Game) Neighbours(x, y int) []Point {

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

type Change struct {
	X		int
	Y		int
	Delta	int
}

func (self *Game) Changes() []Change {

	var ret []Change

	for key, val := range self.box_deltas {				// Iterating over a map, order not deterministic
		ret = append(ret, Change{key.X, key.Y, val})
	}

	return ret
}
