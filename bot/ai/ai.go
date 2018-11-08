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

	// Other initialisation..................................

	rand.Seed(int64(frame.MyBudget() + pid))

	// Ship cleanup and target choice........................

	my_ships := frame.MyShips()

	for _, ship := range my_ships {
		NewTurn(ship)	// May clear the ship's target.
	}

	target_book := hal.Make2dBoolArray(frame.Width(), frame.Height())

	for _, ship := range my_ships {
		SetTarget(ship, target_book)
	}

	TargetSwaps(my_ships, 4)

	// What each ship wants to do right now..................

	for _, ship := range my_ships {
		SetDesires(ship)
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

func ShouldMine(frame *hal.Frame, halite_carried int, pos, tar hal.XYer) bool {

	// Whether a ship -- if it were carrying n halite, at pos, with specified target -- would stop to mine.

	if halite_carried >= 800 {
		return false
	}

	pos_halite := frame.HaliteAt(pos)
	tar_halite := frame.HaliteAt(tar)

	// if frame.InspirationCheck(pos) { pos_halite *= 3 }
	// if frame.InspirationCheck(tar) { tar_halite *= 3 }

	if pos_halite > HappyThreshold(frame) {
		if pos_halite > tar_halite / 3 {			// This is a bit odd since the test even happens when target is dropoff.
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

func TargetSwaps(my_ships []*hal.Ship, cycles int) {

	for cycle := 0; cycle < cycles; cycle++ {

		swap_count := 0

		for i, ship_a := range my_ships {

			if ship_a.TargetIsDropoff() {
				continue
			}

			for _, ship_b := range my_ships[i + 1:] {

				if ship_b.TargetIsDropoff() {
					continue
				}

				a_dist_b := ship_a.Dist(ship_b.Target())
				b_dist_a := ship_b.Dist(ship_a.Target())

				alt_score_a := HaliteDistScore(ship_b.TargetHalite(), a_dist_b)
				alt_score_b := HaliteDistScore(ship_a.TargetHalite(), b_dist_a)

				if alt_score_a + alt_score_b > ship_a.Score + ship_b.Score {

					tmp := ship_a.Target()
					ship_a.SetTarget(ship_b.Target())
					ship_b.SetTarget(tmp)

					ship_a.Score = alt_score_a
					ship_b.Score = alt_score_b

					// ship_a.Frame.Log("Swapped targets for pilots %d, %d (cycle %d)", ship_a.Sid, ship_b.Sid, cycle)
					swap_count++
				}
			}
		}

		if swap_count == 0 {
			return
		}
	}
}
