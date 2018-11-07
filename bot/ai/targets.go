package ai

import (
	"math/rand"
	hal "../core"
)

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




func StupidStep(frame *hal.Frame, pid int) {

	frame.SetPid(pid)

	happy_threshold := HappyThreshold(frame)

	my_ships := frame.MyShips()

	for _, ship := range my_ships {
		NewTurn(ship)
	}

	// TBC
*/
