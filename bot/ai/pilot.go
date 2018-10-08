package ai

import (
	"fmt"
	"math/rand"

	hal "../core"
)

type State int

const (
	Normal		= State(iota)
	Returning
)

type Pilot struct {
	Game					*hal.Game
	Ship					*hal.Ship
	Sid						int
	State					State
	TargetX					int
	TargetY					int
}

func (self *Pilot) Navigate(x, y int) {

	// FIXME: consider wraps

	ship := self.Ship

	dx := x - ship.X
	dy := y - ship.Y

	if dx == 0 && dy == 0 {
		ship.Clear()
		return
	}

	var options []string

	if dx > 0 {
		options = append(options, "e")
	}

	if dx < 0 {
		options = append(options, "w")
	}

	if dy > 0 {
		options = append(options, "s")
	}

	if dy < 0 {
		options = append(options, "n")
	}

	n := rand.Intn(len(options))

	direction := options[n]

	command := fmt.Sprintf("m %d %s", ship.Sid, direction)
	ship.Set(command)
}

func (self *Pilot) Fly() {

	game := self.Game
	ship := self.Ship

	if ship.Halite < game.HaliteAt(ship.X, ship.Y) / 10 {			// We can't move
		return
	}

	if ship.IsOnDropoff() {
		self.State = Normal
		self.NewTarget()
	}

	if self.State == Normal {
		if ship.X == self.TargetX && ship.Y == self.TargetY {
			if ship.Halite > 800 {
				self.State = Returning
			} else if game.HaliteAt(ship.X, ship.Y) < 50 {
				self.State = Returning								// FIXME: choose new target instead
			}
		} else {
			self.Navigate(self.TargetX, self.TargetY)
		}
	}

	if self.State == Returning {

		// FIXME: consider more than the first in the list...

		dropoff := game.MyDropoffs()[0]
		self.Navigate(dropoff.X, dropoff.Y)
	}
}

func (self *Pilot) NewTarget() {

	x := rand.Intn(self.Game.Width() / 2)
	y := rand.Intn(self.Game.Height() / 2)

	factory_x, factory_y := self.Game.MyFactoryXY()

	if factory_x >= self.Game.Width() / 2 {
		x += self.Game.Width() / 2
	}

	if factory_y >= self.Game.Height() / 2 {
		y += self.Game.Height() / 2
	}

	self.TargetX = x
	self.TargetY = y
}
