package main

import (
	ai "./ai"
	hal "./core"
)

func sim_check(frame *hal.Frame, config *ai.Config) (string, int) {

	// Returns the final hash that the real bot will see,
	// if the real bot is matched only against itself...

	frame.Init()		// This should be harmless for the caller.

	var overminds []*ai.Overmind

	for pid := 0; pid < frame.Players(); pid++ {
		overminds = append(overminds, ai.NewOvermind(frame, config, pid))
	}

	for {
		if frame.Turn() == frame.Constants.MAX_TURNS - 1 {
			return frame.Hash(), frame.GroundHalite()
		}

		for _, o := range overminds {
			o.Step(frame)
		}

		frame = frame.SimGen()
	}
}
