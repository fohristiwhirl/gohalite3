package ai

import (
	"fmt"
	hal "../core"
)

func NewTurn(ship *hal.Ship) {

	ship.Command = ""
	ship.Desires = nil

	if ship.OnDropoff() {
		ship.ClearTarget()
	}

	if ship.TargetOK() && ship.TargetHalite() < IgnoreThreshold(ship.Frame) {
		ship.ClearTarget()
	}

	if ship.Dist(ship.NearestDropoff()) > ship.Frame.Constants.MAX_TURNS - ship.Frame.Turn() - 3 {
		ship.FinalDash = true
	}

	if ship.FinalDash || ShouldReturn(ship.Halite) || (ship.Returning && ship.OnDropoff() == false) {
		ship.Returning = true
	} else {
		ship.Returning = false
	}
}

func SetTarget(ship *hal.Ship, target_book [][]bool) {

	if ship.Returning {
		ship.SetTarget(ship.NearestDropoff())
		return
	}

	if ship.TargetOK() {
		return
	}

	frame := ship.Frame
	width := frame.Width()
	height := frame.Height()

	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			if target_book[x][y] {
				continue
			}

			halite := frame.HaliteAtFast(x, y)
/*
			if frame.InspirationCheck(hal.Point{x, y}) {
				halite *= 3
			}
*/
			if halite < IgnoreThreshold(frame) {
				continue
			}

			score := RunShipSim(ship, hal.Point{x, y})

			if ship.TargetOK() == false || score > ship.Score {
				ship.SetTarget(hal.Point{x, y})
				ship.Score = score
			}
		}
	}

	// It's best not to set a default at top because it can confuse the logic.
	// i.e. we want to ignore boxes below a certain halite threshold, even if
	// they have a good score. Our default might be such a square, so comparing
	// against its score might lead us to reject a box we should pick.

	if ship.TargetOK() == false {
		ship.SetTarget(ship)										// Default - my own square
		ship.Score = RunShipSim(ship, hal.Point{ship.X, ship.Y})
	}

	// Set the book for this ship. Note that for dropoff targets,
	// we already returned and so they aren't included here.

	target_book[ship.Target().X][ship.Target().Y] = true
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

					// ship_a.Frame.Log("Swapped targets for ships %d, %d (cycle %d)", ship_a.Sid, ship_b.Sid, cycle)
					swap_count++
				}
			}
		}

		if swap_count == 0 {
			return
		}
	}
}

func FlogTarget(ship *hal.Ship) {
	ship.Frame.Flog(ship.X, ship.Y, fmt.Sprintf("Target: %d %d - Dist: %d", ship.Target().X, ship.Target().Y, ship.Dist(ship.Target())), "")
	ship.Frame.Flog(ship.Target().X, ship.Target().Y, "", "LemonChiffon")
}
