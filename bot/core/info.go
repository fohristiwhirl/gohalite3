package core

import (
	"fmt"
)

// Game info, NOT including trivial properties (e.g. width, players)

func (self *Game) Box(x, y int) *Box {
	x = mod(x, self.width)
	y = mod(y, self.height)
	return self.boxes[x][y]
}

func (self *Game) Ship(x, y int) *Ship {
	x = mod(x, self.width)
	y = mod(y, self.height)
	return self.ship_xy_lookup[Point{x, y}]
}

func (self *Game) Sid(sid int) *Ship {
	return self.ship_id_lookup[sid]
}

func (self *Game) Dropoffs(pid int) []*Dropoff {

	var ret []*Dropoff

	for _, dropoff := range self.dropoffs {
		if dropoff.Owner == pid {
			ret = append(ret, dropoff)
		}
	}

	return ret
}

func (self *Game) MyDropoffs() []*Dropoff {
	return self.Dropoffs(self.pid)
}

func (self *Game) PlayerShips(pid int) []*Ship {

	var ret []*Ship

	for _, ship := range self.ships {
		if ship.Owner == pid {
			ret = append(ret, ship)
		}
	}

	return ret
}

func (self *Game) MyShips() []*Ship {
	return self.PlayerShips(self.pid)
}

func (self *Game) PlayerBudget(pid int) int {
	return self.budgets[pid]
}

func (self *Game) MyBudget() int {
	return self.PlayerBudget(self.pid)
}

func (self *Game) MyFactory() *Dropoff {
	return self.Factory(self.pid)
}

func (self *Game) Factory(pid int) *Dropoff {

	factory := self.dropoffs[pid]

	if factory.Owner != pid || factory.Factory == false {
		panic(fmt.Sprintf("self.dropoffs[%d] wasn't the right factory", pid))
	}

	return factory
}

func (self *Game) ShipCanDropoffAt(ship *Ship, box *Box) bool {

	dropoffs := self.Dropoffs(ship.Owner)

	for _, dropoff := range dropoffs {
		if dropoff.X == box.X && dropoff.Y == box.Y {
			return true
		}
	}
	return false
}
