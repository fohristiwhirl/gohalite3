package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	ai "./ai"
	hal "./core"
)

const (
	NAME = "Fohristiwhirl"
	VERSION = "0"
)

func main() {

	config := new(ai.Config)

	flag.BoolVar(&config.Timeseed, "timeseed", false, "seed RNG with time")
	flag.Parse()

	game := hal.NewGame()

	var longest_turn time.Duration
	var longest_turn_number int

	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("%v", p)
			game.Log("Quitting: %v", p)
			game.Log("Longest turn (%d) took %v", longest_turn_number, longest_turn)
		}
	}()

	// -------------------------------------------------------------------------------

	game.PrePreParse(NAME, VERSION)		// Reads very early data and starts log file.

	if config.Timeseed {
		seed := time.Now().UTC().UnixNano()
		rand.Seed(seed)
		game.LogWithoutTurn("Seeding own RNG: %v", seed)
	}

	game.PreParse()						// Reads the map data.

	overmind := ai.NewOvermind(game, config)
	fmt.Printf("%s %s\n", NAME, VERSION)

	game.LogWithoutTurn("./halite.exe --width %d --height %d -s %v", game.Width(), game.Height(), game.Constants.GameSeed)

	// -------------------------------------------------------------------------------

	for {
		game.Parse()

		start_time := time.Now()

		if config.Timeseed == false {
			rand.Seed(int64(game.Turn() + game.Width() + game.Pid()))
		}

		overmind.Step()

		game.Send()

		if time.Now().Sub(start_time) > longest_turn {
			longest_turn = time.Now().Sub(start_time)
			longest_turn_number = game.Turn()
		}
	}
}
