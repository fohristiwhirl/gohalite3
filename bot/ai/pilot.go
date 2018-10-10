package ai

import (
	"math/rand"
	"sort"

	hal "../core"
)

type State int

const (
	Normal = 				State(iota)
	Returning
)

type Pilot struct {
	Game					*hal.Game
	Overmind				*Overmind
	Ship					*hal.Ship
	Sid						int
	State					State
	Target					hal.XYer
	Desires					[]string
}

func (self *Pilot) GetGame() *hal.Game { return self.Game }
func (self *Pilot) GetX() int { return self.Ship.X }
func (self *Pilot) GetY() int { return self.Ship.Y }
func (self *Pilot) DxDy(other hal.XYer) (int, int) { return hal.DxDy(self, other) }
func (self *Pilot) Dist(other hal.XYer) int { return hal.Dist(self, other) }
func (self *Pilot) SamePlace(other hal.XYer) bool { return hal.SamePlace(self, other) }

func (self *Pilot) DesireNav(target hal.XYer) {

	self.Desires = nil

	dx, dy := self.DxDy(target)

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
		dislikes = append(dislikes, "s")
	} else {
		neutrals = append(neutrals, "s")
		neutrals = append(neutrals, "n")
	}

	rand.Shuffle(len(likes), func(i, j int) {
		likes[i], likes[j] = likes[j], likes[i]
	})

	rand.Shuffle(len(neutrals), func(i, j int) {
		neutrals[i], neutrals[j] = neutrals[j], neutrals[i]
	})

	rand.Shuffle(len(dislikes), func(i, j int) {
		dislikes[i], dislikes[j] = dislikes[j], dislikes[i]
	})

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

	if self.State == Normal { // && self.SamePlace(self.Target) {
		if ship.Halite < 800 {
			if ship.BoxUnder().Halite > 50 {
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
		if self.SamePlace(self.Target) {
			self.State = Returning			// FIXME: maybe choose new target
		} else {
			self.DesireNav(self.Target)
		}
	}

	if self.State == Returning {

		choice := game.MyDropoffs()[0]
		choice_dist := self.Dist(choice)

		for _, dropoff := range game.MyDropoffs()[1:] {
			if self.Dist(dropoff) < choice_dist {
				choice = dropoff
				choice_dist = self.Dist(dropoff)
			}
		}

		self.DesireNav(choice)
	}
}

func (self *Pilot) NewTarget() {

	type Foo struct {
		X		int
		Y		int
		Score	float64
	}

	self.Target = self.Game.Box(self.Ship.X, self.Ship.Y)

	game := self.Game
	width := game.Width()
	height := game.Height()

	pilots := self.Overmind.Pilots

	var all_foo []Foo

	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			dist := self.Dist(hal.Point{x, y})

			score := float64(game.Box(x, y).Halite) / float64((dist + 1))		// Avoid divide by zero

			all_foo = append(all_foo, Foo{
				X: x,
				Y: y,
				Score: score,
			})

		}
	}

	sort.Slice(all_foo, func(a, b int) bool {
		return all_foo[a].Score > all_foo[b].Score				// Highest first
	})

	FooLoop:
	for _, foo := range all_foo {
		for _, pilot := range pilots {
			if pilot != self && pilot.State == Normal && pilot.Target.GetX() == foo.X && pilot.Target.GetY() == foo.Y {
				continue FooLoop
			}
		}
		self.Target = game.Box(foo.X, foo.Y)
		break FooLoop
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
