package ai

import (
	"sort"

	hal "../core"
)

type State int

const (
	Normal = 				State(iota)
	Returning
)

type Box struct {
	X						int
	Y						int
	Score					float64
}

type Pilot struct {
	Game					*hal.Game
	Overmind				*Overmind
	Ship					*hal.Ship
	Sid						int
	State					State
	TargetX					int
	TargetY					int
	Desires					[]string
}

func (self *Pilot) DesireNav(x, y int) {

	self.Desires = nil

	dx, dy := self.Ship.DxDy(x, y)

	if dx == 0 && dy == 0 {
		self.Desires = []string{"o"}
		return
	}

	var likes []string
	var neutrals []string		// Perhaps badly named
	var dislikes []string

	if dx > 0 {
		likes = append(likes, "e")
		dislikes = append(dislikes, "w")
	} else if dx < 0 {
		likes = append(likes, "w")
		dislikes = append(dislikes, "e")
	} else {
		neutrals = append(neutrals, "e")
		neutrals = append(neutrals, "w")
	}

	if dy > 0 {
		likes = append(likes, "s")
		dislikes = append(dislikes, "n")
	} else if dy < 0 {
		likes = append(likes, "n")
		dislikes = append(dislikes, "n")
	} else {
		neutrals = append(neutrals, "s")
		neutrals = append(neutrals, "n")
	}

	self.Desires = append(self.Desires, likes...)
	self.Desires = append(self.Desires, neutrals...)
	self.Desires = append(self.Desires, dislikes...)
	self.Desires = append(self.Desires, "o")
}

func (self *Pilot) SetDesires() {

	game := self.Game
	ship := self.Ship

	// Maybe we want to stay still...

	if ship.Halite < self.Ship.MoveCost() {			// We can't move
		self.Desires = []string{"o"}
		return
	}

	if self.State == Normal {
		if ship.Halite < 800 {
			if ship.HaliteUnder() > 50 {			// We're happy here
				self.Desires = []string{"o"}
				return
			}
		}
	}

	// So we're not holding still...

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
			self.DesireNav(self.TargetX, self.TargetY)
		}
	}

	if self.State == Returning {

		choice := game.MyDropoffs()[0]
		choice_dist := self.Ship.Dist(choice.X, choice.Y)

		for _, dropoff := range game.MyDropoffs()[1:] {
			if self.Ship.Dist(dropoff.X, dropoff.Y) < choice_dist {
				choice = dropoff
				choice_dist = self.Ship.Dist(dropoff.X, dropoff.Y)
			}
		}

		self.DesireNav(choice.X, choice.Y)
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

			score := float64(game.HaliteAt(x, y)) / float64((dist + 1))		// Avoid divide by zero

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
			if pilot != self && pilot.State == Normal && pilot.TargetX == box.X && pilot.TargetY == box.Y {
				continue BoxLoop
			}
		}
		self.TargetX = box.X
		self.TargetY = box.Y
		break BoxLoop
	}
}

func (self *Pilot) LocationAfterMove(s string) (int, int) {

	dx, dy := string_to_dxdy(s)

	x := self.Ship.X + dx
	y := self.Ship.Y + dy

	x = mod(x, self.Game.Width())
	y = mod(y, self.Game.Height())

	return x, y
}
