package main

import (
	ai "./ai"
	hal "./core"
)

func sim_check(game *hal.Game, config *ai.Config) string {

	// Returns the final hash that the real bot will see,
	// if the real bot is matched only against itself...

	game.Init()		// This should be harmless for the caller.

	var overminds []*ai.Overmind

	for pid := 0; pid < game.Players(); pid++ {
		overminds = append(overminds, ai.NewOvermind(game, config, pid))
	}

	for {
		if game.Turn() == game.Constants.MAX_TURNS - 1 {
			return game.Hash()
		}

		for _, o := range overminds {
			o.Step(game)
		}

		game = game.SimGen()
	}
}
