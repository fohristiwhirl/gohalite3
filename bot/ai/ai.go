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

	// Various calls rely on this happening..................
	// Always have this first.

	frame.SetPid(pid)

	// Various other initialisation..........................

	rand.Seed(int64(frame.MyBudget() + pid))
	happy_threshold := HappyThreshold(frame)

	// Ship cleanup and target choice........................

	my_ships := frame.MyShips()

	for _, ship := range my_ships {
		NewTurn(ship)						// May clear the ship's target.
	}

	ChooseTargets(frame, my_ships, pid)		// Only sets targets for ships that need a new one.

	// What each ship wants to do right now..................

	for _, ship := range my_ships {
		SetDesires(ship, happy_threshold)
	}

	// Resolve the desired moves.............................

	move_book := Resolve(frame, my_ships)

	// Other.................................................

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

func ShouldMine(frame *hal.Frame, halite_carried int, pos, tar hal.XYer, happy_threshold int) bool {

	// Whether a ship -- if it were carrying n halite, at pos, with specified target -- would stop to mine.

	if halite_carried >= 800 {
		return false
	}

	pos_halite := frame.HaliteAt(pos)
	tar_halite := frame.HaliteAt(tar)

	if pos_halite > happy_threshold {
		if pos_halite > tar_halite / 3 {			// This is a bit odd since the test even happens when target is dropoff.
			return true
		}
	}

	return false
}

func HappyThreshold(frame *hal.Frame) int {			// Probably bad to call this a lot when simming, will be slow. So cache it.
	return frame.AverageGroundHalite() / 2
}
