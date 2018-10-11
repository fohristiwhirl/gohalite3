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
			game.StopFlog()
		}
	}()

	// -------------------------------------------------------------------------------

	err := game.PrePreParse()			// Reads very early data and starts log file.

	// Both of these fail harmlessly if the directory isn't there:
	game.StartLog(fmt.Sprintf("logs/log%v.txt", game.Pid()))
	game.StartFlog(fmt.Sprintf("flogs/flog-%v-%v.json", game.Constants.GameSeed, game.Pid()))

	if err != nil {
		game.Log("%v", err)
	}

	game.PreParse()						// Reads the map data.

	game.LogWithoutTurn("--------------------------------------------------------------------------------")
	game.LogWithoutTurn("%s %s starting up at %s", NAME, VERSION, time.Now().Format("2006-01-02 15:04:05"))

	if config.Timeseed {
		seed := time.Now().UTC().UnixNano()
		rand.Seed(seed)
		game.LogWithoutTurn("Seeding own RNG: %v", seed)
	}

	overmind := ai.NewOvermind(game, config)
	fmt.Printf("%s %s\n", NAME, VERSION)

	game.LogWithoutTurn("./halite.exe --width %d --height %d -s %v    <%d players>", game.Width(), game.Height(), game.Constants.GameSeed, game.Players())

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
