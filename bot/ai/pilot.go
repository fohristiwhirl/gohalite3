package ai

import (
	"fmt"
	"math/rand"
	"sort"

	hal "../core"
)

type State int

const (
	Normal		= State(iota)
	Returning
)

type Box struct {
	X					int
	Y					int
	Score				int
}

type Pilot struct {
	Game					*hal.Game
	Overmind				*Overmind
	Ship					*hal.Ship
	Sid						int
	State					State
	TargetX					int
	TargetY					int
	Holding					bool		// Intends to not move this turn
}

func (self *Pilot) Navigate(x, y int) {

	dx, dy := self.Ship.DxDy(x, y)

	if dx == 0 && dy == 0 {
		self.Ship.Clear()
		return
	}

	var options []string

	if dx > 0 {
		if self.Overmind.Booker(self.Ship.X + 1, self.Ship.Y) == nil {
			options = append(options, "e")
		}
	}

	if dx < 0 {
		if self.Overmind.Booker(self.Ship.X - 1, self.Ship.Y) == nil {
			options = append(options, "w")
		}
	}

	if dy > 0 {
		if self.Overmind.Booker(self.Ship.X, self.Ship.Y + 1) == nil {
			options = append(options, "s")
		}
	}

	if dy < 0 {
		if self.Overmind.Booker(self.Ship.X, self.Ship.Y - 1) == nil {
			options = append(options, "n")
		}
	}

	if len(options) == 0 {
		return
	}

	n := rand.Intn(len(options))
	direction := options[n]

	self.Prepare(direction)
}

func (self *Pilot) MaybeHold() {

	ship := self.Ship

	if ship.Halite < self.Ship.MoveCost() {								// We can't move
		self.Hold()
		return
	}

	if self.State == Normal {
		if ship.Halite < 800 {
			if ship.HaliteUnder() > 50 {
				self.Hold()
				return
			}
		}
	}
}

func (self *Pilot) Hold() {
	self.Holding = true
	self.Overmind.SetBook(self, self.Ship.X, self.Ship.Y)
}

func (self *Pilot) Fly() {

	if self.Holding {
		return
	}

	game := self.Game
	ship := self.Ship

	if ship.OnDropoff() {
		self.State = Normal
		self.NewTarget()
	}

	// We're not holding, so if we're on our target square, we're either
	// about to return or about to change target.

	if self.State == Normal {
		if ship.X == self.TargetX && ship.Y == self.TargetY {
			self.State = Returning			// FIXME: maybe choose new target
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

	self.TargetX = self.Ship.X
	self.TargetY = self.Ship.Y

	game := self.Game
	width := game.Width()
	height := game.Height()

	pilots := self.Overmind.Pilots

	var boxes []Box

	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			dist := self.Ship.Dist(x, y)

			score := game.HaliteAt(x, y) / (dist + 1)		// Avoid divide by zero

			boxes = append(boxes, Box{
				X: x,
				Y: y,
				Score: score,
			})

		}
	}

	sort.Slice(boxes, func(a, b int) bool {
		return boxes[a].Score > boxes[b].Score				// Highest first
	})

	BoxLoop:
	for _, box := range boxes {
		for _, pilot := range pilots {
			if pilot.TargetX == box.X && pilot.TargetY == box.Y {
				continue BoxLoop
			}
			self.TargetX = box.X
			self.TargetY = box.Y
			break BoxLoop
		}
	}
}

func (self *Pilot) Prepare(d string) {

	command := fmt.Sprintf("m %d %s", self.Ship.Sid, d)
	self.Ship.Set(command)		// But undo this if it's not a valid direction. See default.

	switch d {

	case "e":
		self.Overmind.SetBook(self, self.Ship.X + 1, self.Ship.Y)
	case "w":
		self.Overmind.SetBook(self, self.Ship.X - 1, self.Ship.Y)
	case "s":
		self.Overmind.SetBook(self, self.Ship.X, self.Ship.Y + 1)
	case "n":
		self.Overmind.SetBook(self, self.Ship.X, self.Ship.Y - 1)
	default:
		self.Ship.Clear()

	}
}
