package ai

import (
	"fmt"
	"math/rand"
	"sort"

	hal "../core"
)

type Pilot struct {			// Best not to have a reference to the game since that changes in sims.
	Overmind				*Overmind
	Ship					*hal.Ship
	Sid						int
	Target					hal.Point	// This is NOT used to preserve target info between turns.
	Score					float32		// Score if our target is a mineable box.
	Desires					[]string
	Returning				bool
}

func (self *Pilot) NewTurn() {
	self.Desires = nil
	self.Target = hal.Point{self.Ship.X, self.Ship.Y}
	self.Score = 0

	if self.OnDropoff() {
		self.Returning = false
	}
}

func (self *Pilot) SetTarget() {

	// Note that the ship may still not move if it's happy where it is.

	if self.FinalDash() {
		self.Target = self.NearestDropoff().Point()
		self.Returning = true
		return
	}

	if self.Ship.Halite > 500 || self.Returning {
		self.Target = self.NearestDropoff().Point()
		self.Returning = true
		return
	}

	self.TargetBestBox()
	self.Overmind.TargetBook[self.Target.X][self.Target.Y] = true		// Only for normal targets
}

func (self *Pilot) TargetBestBox() {

	self.Target = hal.Point{-1, -1}		// Spooky
	self.Score = -999999

	frame := self.GetFrame()
	width := frame.Width()
	height := frame.Height()

	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			if self.Overmind.TargetBook[x][y] {
				continue
			}

			halite := frame.HaliteAtFast(x, y)

			if halite < self.Overmind.IgnoreThreshold {
				continue
			}

			dist := self.Dist(hal.Point{x, y})
			score := halite_dist_score(halite, dist)

			if score > self.Score {
				self.Target = hal.Point{x, y}
				self.Score = score
			}
		}
	}

	// It's best not to set a default at top because it can confuse the logic.
	// i.e. we want to ignore boxes below a certain halite threshold, even if
	// they have a good score. Our default might be such a square, so comparing
	// against its score might lead us to reject a box we should pick.

	if self.Target == (hal.Point{-1, -1}) {
		self.Target = self.Ship.Point()											// Default - my own square
		self.Score = halite_dist_score(self.Ship.HaliteAt(), 0)
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

	if self.Overmind.ShouldMine(self.Ship.Halite, self, self.Target) {
		self.Desires = []string{"o"}
		return
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

	// If lowish halite, prefer mining en route...

	frame := self.GetFrame()

	if self.Ship.Halite < 750 {

		sort.Slice(likes, func(a, b int) bool {

			halite_after_move := self.Ship.Halite - self.Ship.MoveCost()

			loc1 := self.LocationAfterMove(likes[a])
			loc2 := self.LocationAfterMove(likes[b])

			would_mine_1 := self.Overmind.ShouldMine(halite_after_move, loc1, self.Target)
			would_mine_2 := self.Overmind.ShouldMine(halite_after_move, loc2, self.Target)

			if would_mine_1 && would_mine_2 == false {				// Only mines at 1
				return true
			} else if would_mine_1 == false && would_mine_2 {		// Only mines at 2
				return false
			} else if would_mine_1 && would_mine_2 {				// Mines at both, choose higher
				return frame.HaliteAtFast(loc1.X, loc1.Y) > frame.HaliteAtFast(loc2.X, loc2.Y)
			} else {												// Mines at neither, choose lower
				return frame.HaliteAtFast(loc1.X, loc1.Y) < frame.HaliteAtFast(loc2.X, loc2.Y)
			}
		})

	} else {

		rand.Shuffle(len(likes), func(i, j int) {
			likes[i], likes[j] = likes[j], likes[i]
		})

	}

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
	return self.Dist(self.NearestDropoff()) > self.GetFrame().Constants.MAX_TURNS - self.GetFrame().Turn() - 3
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

func (self *Pilot) TargetHalite() int {
	return self.GetFrame().HaliteAt(self.Target)
}

func (self *Pilot) GetFrame() *hal.Frame { return self.Overmind.Frame }
func (self *Pilot) GetX() int { return self.Ship.X }
func (self *Pilot) GetY() int { return self.Ship.Y }
func (self *Pilot) DxDy(other hal.XYer) (int, int) { return hal.DxDy(self, other) }
func (self *Pilot) Dist(other hal.XYer) int { return hal.Dist(self, other) }
func (self *Pilot) SamePlace(other hal.XYer) bool { return hal.SamePlace(self, other) }

func (self *Pilot) FlogTarget() {
	self.GetFrame().Flog(self.Ship.X, self.Ship.Y, fmt.Sprintf("Target: %d %d - Dist: %d", self.Target.X, self.Target.Y, hal.Dist(self, self.Target)), "")
	self.GetFrame().Flog(self.Target.X, self.Target.Y, "", "LemonChiffon")
}
