package ai


import (
	"math/rand"
	hal "../core"
)


/*
type ShipSim struct {
	realframe			*hal.Frame
	x					int
	y					int
	halite_carried		int
	turn				int
	target				*hal.Point
	cells				map[hal.Point]int
}
*/


func RunShipSim(real_ship *hal.Ship, proposed_target hal.Point) float32 {

	real_frame := real_ship.Frame
	x := real_ship.X
	y := real_ship.Y
	halite_carried := real_ship.Halite
	turn := real_frame.Turn()
	target := proposed_target
	cells := make(map[hal.Point]int)

	turn -= 1		// Make life convenient by incrementing the turn at start of loop...

	for {

		turn += 1

		halite_at_ship, ok := cells[hal.Point{x, y}]
		if ok == false {
			halite_at_ship = real_frame.HaliteAtFast(x, y)
		}

		halite_at_target, ok := cells[target]
		if ok == false {
			halite_at_target = real_frame.HaliteAt(target)
		}

		amount_to_mine := (halite_at_ship + real_frame.Constants.EXTRACT_RATIO - 1) / real_frame.Constants.EXTRACT_RATIO
		move_cost := halite_at_ship / real_frame.Constants.MOVE_COST_RATIO

		// FIXME: actually simulate the return as well
		if ShouldReturn(halite_carried) || halite_at_target == 0 {
			// Return a sort of efficiency stat
			return float32(halite_carried - real_ship.Halite) / float32(1 + turn - real_frame.Turn())
		}

		if move_cost > halite_carried || ShouldMine(real_frame, halite_carried, halite_at_ship, halite_at_target) {
			halite_carried += amount_to_mine
			cells[hal.Point{x, y}] = halite_at_ship - amount_to_mine
			continue
		}

		// DesireNav....................................................

		dx, dy := (&hal.Cell{real_frame, x, y}).DxDy(target)

		if dx == 0 && dy == 0 {
			halite_carried += amount_to_mine
			cells[hal.Point{x, y}] = halite_at_ship - amount_to_mine
			continue
		}

		var likes []string

		if dx > 0 {
			likes = append(likes, "e")
		} else if dx < 0 {
			likes = append(likes, "w")
		}

		if dy > 0 {
			likes = append(likes, "s")
		} else if dy < 0 {
			likes = append(likes, "n")
		}

		// FIXME: do proper sort

		rand.Shuffle(len(likes), func(i, j int) {
			likes[i], likes[j] = likes[j], likes[i]
		})

		move_x, move_y := hal.StringToDxDy(likes[0])

		x = hal.Mod(x + move_x, real_frame.Width())
		y = hal.Mod(y + move_y, real_frame.Height())

	}
}
