package ai

import (
	"math/rand"
	"sort"

	hal "../core"
)

const (
	DROPOFF_SPACING = 12
	NICE_THRESHOLD = 8000
)

func Step(frame *hal.Frame, pid int, allow_build bool) {

	frame.SetPid(pid)		// Always have this first.

	rand.Seed(int64(frame.MyBudget() + pid))

	my_ships := frame.MyShips()

	for _, ship := range my_ships {
		NewTurn(ship)
	}

	target_book := hal.Make2dBoolArray(frame.Width(), frame.Height())

	for _, ship := range my_ships {
		SetTarget(ship, target_book)
	}

	TargetSwaps(my_ships, 4)

	for _, ship := range my_ships {
		SetDesires(ship)
	}

	move_book := Resolve(frame, my_ships)

	if allow_build {
		MaybeBuild(frame, my_ships, move_book)
	}

	for _, ship := range my_ships {
		FlogTarget(ship)
	}

	return
}

func MaybeBuild(frame *hal.Frame, my_ships []*hal.Ship, move_book *MoveBook) {

	budget := frame.MyBudget()

	factory := frame.MyFactory()
	willing := true

	if frame.InitialGroundHalite() / (frame.GroundHalite() + 1) >= 2 {		// remember int division, also div-by-zero
		willing = false
	}

	if frame.Turn() >= frame.Constants.MAX_TURNS / 2 {
		willing = false
	}

	if budget >= frame.Constants.NEW_ENTITY_ENERGY_COST && move_book.Booker(factory) == nil && willing {
		frame.SetGenerate(true)
		budget -= frame.Constants.NEW_ENTITY_ENERGY_COST
	}

	// -------------------------------------------

	var possible_constructs []*hal.Ship

	for _, ship := range my_ships {

		if ship.Dist(ship.NearestDropoff()) < DROPOFF_SPACING {
			continue
		}

		if frame.WealthMap().Values[ship.X][ship.Y] < NICE_THRESHOLD {
			continue
		}

		if frame.HaliteAtFast(ship.X, ship.Y) == 0 {		// Cheap way to avoid building on enemy dropoff / factory
			continue
		}

		possible_constructs = append(possible_constructs, ship)
	}

	sort.Slice(possible_constructs, func (a, b int) bool {

		return	frame.WealthMap().Values[possible_constructs[a].X][possible_constructs[a].Y] <
				frame.WealthMap().Values[possible_constructs[b].X][possible_constructs[b].Y]
	})

	for _, ship := range possible_constructs {
		if ship.Halite + frame.HaliteAtFast(ship.X, ship.Y) + budget >= frame.Constants.DROPOFF_COST {
			ship.Command = "c"
			frame.Log("Ship %d building dropoff (wmap: %d)", ship.Sid, frame.WealthMap().Values[ship.X][ship.Y])
			break
		}
	}
}

func ShouldMine(frame *hal.Frame, halite_carried, halite_at_ship, halite_at_target int) bool {

	// Whether a ship...
	//		- if it were carrying <halite_carried>
	//		- with <halite_at_ship> underneath it
	//		- with <halite_at_target> at its target
	// ...would stop to mine.

	if halite_carried >= 800 {
		return false
	}

	if halite_at_ship > HappyThreshold(frame) {
		if halite_at_ship > halite_at_target / 3 {			// This is a bit odd since the test even happens when target is dropoff.
			return true
		}
	}

	return false
}

func ShouldReturn(ship *hal.Ship) bool {			// Could consider dist to dropoff, etc
	return ship.Halite > 500
}

func HappyThreshold(frame *hal.Frame) int {
	return frame.AverageGroundHalite() / 2
}

func IgnoreThreshold(frame *hal.Frame) int {
	return frame.AverageGroundHalite() * 2 / 3
}

func HaliteDistScore(halite, dist int) float32 {
	return float32(halite) / float32((dist + 1) * (dist + 1))	// Avoid div-by-zero
}
