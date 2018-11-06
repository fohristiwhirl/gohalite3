package ai

import (
	"fmt"
	"math/rand"
	"sort"

	hal "../core"
)

func (self *Overmind) NewTurn(ship *hal.Ship) {
	ship.Desires = nil
	ship.Target = ship.Point()
	ship.TargetOK = true
	ship.Score = 0

	if ship.OnDropoff() {
		ship.Returning = false
	}
}

func (self *Overmind) SetTarget(ship *hal.Ship) {

	// Note that the ship may still not move if it's happy where it is.

	if self.FinalDash(ship) {
		ship.Target = ship.NearestDropoff().Point()
		ship.Returning = true
		return
	}

	if ship.Halite > 500 || ship.Returning {
		ship.Target = ship.NearestDropoff().Point()
		ship.Returning = true
		return
	}

	self.TargetBestBox(ship)
	self.TargetBook[ship.Target.X][ship.Target.Y] = true		// Only for normal targets
}

func (self *Overmind) TargetBestBox(ship *hal.Ship) {

	ship.TargetOK = false
	ship.Score = -999999

	width := self.Frame.Width()
	height := self.Frame.Height()

	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			if self.TargetBook[x][y] {
				continue
			}

			halite := self.Frame.HaliteAtFast(x, y)

			if halite < self.IgnoreThreshold {
				continue
			}

			dist := ship.Dist(hal.Point{x, y})
			score := halite_dist_score(halite, dist)

			if score > ship.Score {
				ship.Target = hal.Point{x, y}
				ship.TargetOK = true
				ship.Score = score
			}
		}
	}

	// It's best not to set a default at top because it can confuse the logic.
	// i.e. we want to ignore boxes below a certain halite threshold, even if
	// they have a good score. Our default might be such a square, so comparing
	// against its score might lead us to reject a box we should pick.

	if ship.TargetOK == false {
		ship.Target = ship.Point()											// Default - my own square
		ship.TargetOK = true
		ship.Score = halite_dist_score(ship.HaliteAt(), 0)
	}
}

func (self *Overmind) SetDesires(ship *hal.Ship) {

	// Maybe we can't move...

	if ship.Halite < ship.MoveCost() {
		ship.Desires = []string{"o"}
		return
	}

	// Maybe we're on a mad dash to deliver stuff before end...

	if self.FinalDash(ship) {
		self.DesireNav(ship)
		return
	}

	// Maybe we're happy where we are...

	if self.ShouldMine(ship.Halite, ship, ship.Target) {
		ship.Desires = []string{"o"}
		return
	}

	// Normal case...

	self.DesireNav(ship)
}

func (self *Overmind) DesireNav(ship *hal.Ship) {

	ship.Desires = nil
	dx, dy := ship.DxDy(ship.Target)

	if dx == 0 && dy == 0 {
		ship.Desires = []string{"o"}
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

	if ship.Halite < 750 {

		sort.Slice(likes, func(a, b int) bool {

			halite_after_move := ship.Halite - ship.MoveCost()

			loc1 := ship.LocationAfterMove(likes[a])
			loc2 := ship.LocationAfterMove(likes[b])

			would_mine_1 := self.ShouldMine(halite_after_move, loc1, ship.Target)
			would_mine_2 := self.ShouldMine(halite_after_move, loc2, ship.Target)

			if would_mine_1 && would_mine_2 == false {				// Only mines at 1
				return true
			} else if would_mine_1 == false && would_mine_2 {		// Only mines at 2
				return false
			} else if would_mine_1 && would_mine_2 {				// Mines at both, choose higher
				return self.Frame.HaliteAtFast(loc1.X, loc1.Y) > self.Frame.HaliteAtFast(loc2.X, loc2.Y)
			} else {												// Mines at neither, choose lower
				return self.Frame.HaliteAtFast(loc1.X, loc1.Y) < self.Frame.HaliteAtFast(loc2.X, loc2.Y)
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

	ship.Desires = append(ship.Desires, likes...)
	ship.Desires = append(ship.Desires, neutrals...)
	ship.Desires = append(ship.Desires, dislikes...)
	ship.Desires = append(ship.Desires, "o")
}

func (self *Overmind) FinalDash(ship *hal.Ship) bool {
	return ship.Dist(ship.NearestDropoff()) > self.Frame.Constants.MAX_TURNS - self.Frame.Turn() - 3
}

func (self *Overmind) FlogTarget(ship *hal.Ship) {
	self.Frame.Flog(ship.X, ship.Y, fmt.Sprintf("Target: %d %d - Dist: %d", ship.Target.X, ship.Target.Y, ship.Dist(ship.Target)), "")
	self.Frame.Flog(ship.Target.X, ship.Target.Y, "", "LemonChiffon")
}
