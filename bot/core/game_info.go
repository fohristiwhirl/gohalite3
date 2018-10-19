package core

import (
	"fmt"
)

func (self *Game) Pid() int { return self.pid }
func (self *Game) Turn() int { return self.turn }
func (self *Game) Width() int { return self.width }
func (self *Game) Height() int { return self.height }
func (self *Game) Players() int { return self.players }

func (self *Game) BoxAt(pos XYer) *Box {
	x := Mod(pos.GetX(), self.width)
	y := Mod(pos.GetY(), self.height)
	return self.boxes[x][y]
}

func (self *Game) BoxAtFast(x, y int) *Box {		// For when caller is sure x and y are in bounds
	return self.boxes[x][y]
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
			count += self.boxes[x][y].Halite
		}
	}
	return count
}

func (self *Game) ChangedBoxes() []*Box {
	return self.changed_boxes					// No real need for a defensive copy.
}
