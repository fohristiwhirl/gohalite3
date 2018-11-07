package ai

import (
	"math/rand"
	hal "../core"
)

func NewTurn(ship *hal.Ship, move_on_threshold int) {

	ship.Command = ""
	ship.Desires = nil

	if ship.Dist(ship.NearestDropoff()) > ship.Frame.Constants.MAX_TURNS - ship.Frame.Turn() - 3 {
		ship.FinalDash = true
	}

	if ship.FinalDash || ShouldReturn(ship) {
		ship.Returning = true
	} else {
		ship.Returning = false
	}

	// ------------------------------------------------------------

	if ship.Returning {
		ship.SetTarget(ship.NearestDropoff())
	}

	// If we're at our target and it has little halite, find a new one. Works if the target is dropoff too.

	if ship.TargetOK() && ship.Dist(ship.Target()) == 0 && ship.HaliteAt() < move_on_threshold {
		ship.ClearTarget()
	}
}

func ChooseTargets(frame *hal.Frame, my_ships []*hal.Ship, pid int) {

	for _, ship := range my_ships {

		if ship.TargetOK() {
			continue
		}

		x := rand.Intn(frame.Width())
		y := rand.Intn(frame.Height())

		ship.SetTarget(hal.Point{x, y})
	}
}




/*

		for n := 0; n < 10; n++ {

			x := rand.Intn(frame.Width())
			y := rand.Intn(frame.Height())

			ship.Target = hal.Point{x, y}

			sim_frame := frame.Remake()

			sim_frame.DeleteEnemies()

			// TBC




		}
	}
}
*/


/*
func StupidStep(frame *hal.Frame, pid int) {

	frame.SetPid(pid)

	happy_threshold := HappyThreshold(frame)

	my_ships := frame.MyShips()

	for _, ship := range my_ships {
		NewTurn(ship)
	}

	// TBC
*/
