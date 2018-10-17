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
	Target					*hal.Box				// Currently this is not allowed to be nil
	Desires					[]string
	FinalDash				bool
}

func (self *Pilot) SetTarget() {

	if self.Dist(self.NearestDropoff()) > self.Game.Constants.MAX_TURNS - self.Game.Turn() - 3 {
		self.FinalDash = true
	}

	if self.FinalDash {
		self.Target = self.NearestDropoff().Box()
		return
	}

	if self.SamePlace(self.Target) {
		if self.Ship.Halite > 500 {
			self.Target = self.NearestDropoff().Box()
		} else {
			self.NewTarget()
		}
		return
	}
}

func (self *Pilot) SetDesires() {

	// Maybe we can't move...

	if self.Ship.Halite < self.MoveCost() {
		self.Desires = []string{"o"}
		return
	}

	// Maybe we're on a mad dash to deliver stuff before end...

	if self.FinalDash {
		self.DesireNav(self.Target)
		return
	}

	// Maybe we're happy where we are...

	if self.Ship.Halite < 800 {
		if self.Box().Halite > 50 {
			self.Desires = []string{"o"}
			return
		}
	}

	// Normal case...

	self.DesireNav(self.Target)
}

func (self *Pilot) NewTarget() {

	old_target := self.Target

	type Option struct {
		Box		*hal.Box
		Score	float32
	}

	self.Target = self.Box()

	game := self.Game
	width := game.Width()
	height := game.Height()

	pilots := self.Overmind.Pilots

	var all_options []Option

	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			box := game.BoxAt(hal.Point{x, y})
			dist := self.Dist(box)

			score := float32(box.Halite) / float32((dist + 1) * (dist + 1))		// Avoid divide by zero

			if box.Halite <= 20 {
				score -= 10000
			}

			all_options = append(all_options, Option{
				Box: box,
				Score: score,
			})

		}
	}

	sort.Slice(all_options, func(a, b int) bool {
		return all_options[a].Score > all_options[b].Score					// Highest first
	})

	Outer:
	for _, o := range all_options {
		for _, pilot := range pilots {
			if pilot.Target == o.Box {
				continue Outer
			}
		}
		self.Target = o.Box
		break Outer
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

func (self *Pilot) Flog() {

	if self.CanDropoffAt(self.Target) {
		self.Game.Flog(self.Game.Turn(), self.Ship.X, self.Ship.Y, "Returning")
		return
	}

	style := `color: #ffffff`
	if (self.Dist(self.Target) == 0) {
		style = `color: #d9b3ff`
	}
	msg := fmt.Sprintf(`Target: %v, %v &ndash; <span style="%v">dist: %v</span>`, self.Target.X, self.Target.Y, style, self.Dist(self.Target))
	self.Game.Flog(self.Game.Turn(), self.Ship.X, self.Ship.Y, msg)
}

func (self *Pilot) OnDropoff() bool {
	return self.Ship.OnDropoff()
}

func (self *Pilot) MoveCost() int {
	return self.Ship.MoveCost()
}

func (self *Pilot) NearestDropoff() *hal.Dropoff {
	return self.Ship.NearestDropoff()
}

func (self *Pilot) CanDropoffAt(pos hal.XYer) bool {
	return self.Ship.CanDropoffAt(pos)
}

func (self *Pilot) TargetIsDropoff() bool {
	return self.Ship.CanDropoffAt(self.Target)
}

func (self *Pilot) LocationAfterMove(s string) hal.Point {
	return self.Ship.LocationAfterMove(s)
}

func (self *Pilot) Box() *hal.Box {
	return self.Game.BoxAt(self)
}

func (self *Pilot) GetGame() *hal.Game { return self.Game }
func (self *Pilot) GetX() int { return self.Ship.X }
func (self *Pilot) GetY() int { return self.Ship.Y }
func (self *Pilot) DxDy(other hal.XYer) (int, int) { return hal.DxDy(self, other) }
func (self *Pilot) Dist(other hal.XYer) int { return hal.Dist(self, other) }
func (self *Pilot) SamePlace(other hal.XYer) bool { return hal.SamePlace(self, other) }
