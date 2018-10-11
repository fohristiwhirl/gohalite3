package ai

import (
	"fmt"
	"math/rand"
	"sort"

	hal "../core"
)

type Pilot struct {
	Game					*hal.Game
	Overmind				*Overmind
	Ship					*hal.Ship
	Sid						int
	Target					*hal.Box
	Desires					[]string
}

func (self *Pilot) GetGame() *hal.Game { return self.Game }
func (self *Pilot) GetX() int { return self.Ship.X }
func (self *Pilot) GetY() int { return self.Ship.Y }
func (self *Pilot) DxDy(other hal.XYer) (int, int) { return hal.DxDy(self, other) }
func (self *Pilot) Dist(other hal.XYer) int { return hal.Dist(self, other) }
func (self *Pilot) SamePlace(other hal.XYer) bool { return hal.SamePlace(self, other) }

func (self *Pilot) Flog() {
	style := `color: #ffffff`
	if (self.Dist(self.Target) == 0) {
		style = `color: #d9b3ff`
	}
	msg := fmt.Sprintf(`Target: <span style="%s">%d %d</span>`, style, self.Target.X, self.Target.Y)
	self.Game.Flog(self.Game.Turn(), self.Ship.X, self.Ship.Y, msg)
}

func (self *Pilot) SetDesires() {

	ship := self.Ship

	// Maybe we want to stay still...

	if ship.Halite < ship.MoveCost() {				// We can't move
		self.Desires = []string{"o"}
		return
	}

	if ship.Halite < 800 {
		if ship.Box().Halite > 20 {					// We are happy where we are
			self.Desires = []string{"o"}
			return
		}
	}

	// We're not holding, so if we're on our target square, we're either
	// about to return or about to change target.

	if self.SamePlace(self.Target) {
		if ship.Halite > 500 {
			self.Target = ship.NearestDropoff().Box()
		} else {
			self.NewTarget()
		}
	}

	self.DesireNav(self.Target)
}

func (self *Pilot) NewTarget() {

	old_target := self.Target

	type Foo struct {
		Box		*hal.Box
		Score	float32
	}

	self.Target = self.Ship.Box()

	game := self.Game
	width := game.Width()
	height := game.Height()

	pilots := self.Overmind.Pilots

	var all_foo []Foo

	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			box := game.Box(x, y)

			dist := self.Dist(box)

			score := float32(box.Halite) / float32((dist + 1) * (dist + 1))		// Avoid divide by zero

			all_foo = append(all_foo, Foo{
				Box: box,
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
			if pilot.Target == foo.Box {
				continue FooLoop
			}
		}
		self.Target = foo.Box
		break FooLoop
	}

	if self.Target.Halite <= old_target.Halite {
		self.Target = old_target
	}
}

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

func (self *Pilot) LocationAfterMove(s string) (int, int) {

	dx, dy := string_to_dxdy(s)

	x := self.Ship.X + dx
	y := self.Ship.Y + dy

	x = mod(x, self.Game.Width())
	y = mod(y, self.Game.Height())

	return x, y
}
