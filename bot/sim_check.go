package main

import (
	ai "./ai"
	hal "./core"
)

func sim_check(real_frame *hal.Frame) (string, int) {

	// Returns the final hash that the real bot will see,
	// if the real bot is matched only against itself...

	frame := real_frame.Remake(false)

	for {
		if frame.Turn() == frame.Constants.MAX_TURNS - 1 {
			return frame.Hash(), frame.GroundHalite()
		}

		for pid := 0; pid < frame.Players(); pid++ {
			ai.Step(frame, pid, true)
		}

		frame = frame.SimGen()
	}
}
