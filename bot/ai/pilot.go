package ai

import (
	"fmt"
	"math/rand"

	hal "../core"
)

type Pilot struct {
	Game					*hal.Game
	Overmind				*Overmind
	Ship					*hal.Ship
	Sid						int
	Target					*hal.Box	// Currently this is not allowed to be nil. It is also NOT used to preserve target info between turns.
	Score					float32		// Score if our target is a mineable box.
	Desires					[]string
	Returning				bool
}

func (self *Pilot) NewTurn() {
	self.Desires = nil
	self.Target = self.Box()
	self.Score = 0

	if self.OnDropoff() {
		self.Returning = false
	}
}

func (self *Pilot) SetTarget() {

	// Note that the ship may still not move if it's happy where it is.

	if self.FinalDash() {
		self.Target = self.NearestDropoff().Box()
		self.Returning = true
		return
	}

	if self.Ship.Halite > 500 {					// || self.Returning {		-- deleted for test
		self.Target = self.NearestDropoff().Box()
		self.Returning = true
		return
	}

	self.TargetBestBox()
	self.Overmind.TargetBook[self.Target.X][self.Target.Y] = true		// Only for normal targets
}

func (self *Pilot) TargetBestBox() {

	self.Target = nil
	self.Score = -999999

	game := self.Game
	width := game.Width()
	height := game.Height()

	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			if self.Overmind.TargetBook[x][y] {
				continue
			}

			box := game.BoxAtFast(x, y)

			if box.Halite < self.Overmind.IgnoreThreshold {
				continue
			}

			dist := self.Dist(box)
			score := halite_dist_score(box.Halite, dist)

			if score > self.Score {
				self.Target = box
				self.Score = score
			}
		}
	}

	// It's best not to set a default at top because it can confuse the logic.
	// i.e. we want to ignore boxes below a certain halite threshold, even if
	// they have a good score. Our default might be such a square, so comparing
	// against its core might lead us to reject a box we should pick.

	if self.Target == nil {
		self.Target = self.Box()											// Default - my own square
		self.Score = halite_dist_score(self.Box().Halite, 0)
	}
}

func (self *Pilot) SetDesires() {

	// Maybe we can't move...

	if self.Ship.Halite < self.MoveCost() {
		self.Desires = []string{"o"}
		return
	}

	// Maybe we're on a mad dash to deliver stuff before end...

	if self.FinalDash() {
		self.DesireNav(self.Target)
		return
	}

	// Maybe we're happy where we are...

	if self.Ship.Halite < 800 {
		if self.Box().Halite > self.Overmind.HappyThreshold {
			if self.Box().Halite > self.Target.Halite / 3 {
				self.Desires = []string{"o"}
				return
			}
		}
	}

	// Normal case...

	self.DesireNav(self.Target)
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

func (self *Pilot) FinalDash() bool {
	return self.Dist(self.NearestDropoff()) > self.Game.Constants.MAX_TURNS - self.Game.Turn() - 3
}

func (self *Pilot) Flog() {

	if self.CanDropoffAt(self.Target) {
		self.Game.Flog(self.Game.Turn(), self.Ship.X, self.Ship.Y, "Returning")
		return
	}

	msg := fmt.Sprintf("Target: %v, %v &ndash; dist: %v", self.Target.X, self.Target.Y, self.Dist(self.Target))
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
