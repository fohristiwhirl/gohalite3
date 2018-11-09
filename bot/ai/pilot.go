package ai

import (
	"math/rand"
	"sort"

	hal "../core"
)

func SetDesires(ship *hal.Ship) {

	// Maybe we can't move...

	if ship.Halite < ship.MoveCost() {
		ship.Desires = []string{"o"}
		return
	}

	// Maybe we're on a mad dash to deliver stuff before end...

	if ship.FinalDash {
		DesireNav(ship)
		return
	}

	// Maybe we're happy where we are...

	if ShouldMine(ship.Frame, ship.Halite, ship, ship.Target()) {
		ship.Desires = []string{"o"}
		return
	}

	// Normal case...

	DesireNav(ship)
}

func DesireNav(ship *hal.Ship) {

	ship.Desires = nil
	frame := ship.Frame
	dx, dy := ship.DxDy(ship.Target())

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

	// We may prefer one square to the other if 2 are available...

	sort.Slice(likes, func(a, b int) bool {

		halite_after_move := ship.Halite - ship.MoveCost()

		loc1 := ship.LocationAfterMove(likes[a])
		loc2 := ship.LocationAfterMove(likes[b])

		would_mine_1 := ShouldMine(frame, halite_after_move, loc1, ship.Target())
		would_mine_2 := ShouldMine(frame, halite_after_move, loc2, ship.Target())

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
