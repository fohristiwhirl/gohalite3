package ai

import (
	hal "../core"
)

func NewTurn(ship *hal.Ship) {

	ship.Command = ""
	ship.Desires = nil

	ship.ClearTarget()

	if ship.Dist(ship.NearestDropoff()) > ship.Frame.Constants.MAX_TURNS - ship.Frame.Turn() - 3 {
		ship.FinalDash = true
	}

	if ship.FinalDash || ShouldReturn(ship) || (ship.Returning && ship.OnDropoff() == false) {
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

	ship.ClearTarget()

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

			dist := ship.Dist(hal.Point{x, y})
			score := HaliteDistScore(halite, dist)

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
		ship.Score = HaliteDistScore(ship.HaliteAt(), 0)
	}

	// Set the book for this ship. Note that for dropoff targets,
	// we already returned and so they aren't included here.

	target_book[ship.Target().X][ship.Target().Y] = true
}
