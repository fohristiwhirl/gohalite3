package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	ai "./ai"
	hal "./core"
)

const (
	NAME = "Fohristiwhirl"
	VERSION = "10"				// hash is c3e44812d9ff34ceda884bc045dcc0f5e605605e
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
			game.Log("Last known hash: %s", game.Hash())
			game.Log("Longest turn (%d) took %v", longest_turn_number, longest_turn)
			game.StopLog()
			game.StopFlog()
		}
	}()

	// -------------------------------------------------------------------------------

	game.PrePreParse()				// Reads very early data (including PID, needed for log names)

	// Both of these fail harmlessly if the directory isn't there:
	game.StartLog(fmt.Sprintf("logs/log-%v.txt", game.Pid()))
	game.StartFlog(fmt.Sprintf("flogs/flog-%v-%v.json", game.Constants.GameSeed, game.Pid()))

	game.PreParse()					// Reads the map data.

	game.LogWithoutTurn("--------------------------------------------------------------------------------")
	game.LogWithoutTurn("%s %s starting up at %s", NAME, VERSION, time.Now().Format("2006-01-02 15:04:05"))

	if config.Timeseed {
		seed := time.Now().UTC().UnixNano()
		rand.Seed(seed)
		game.LogWithoutTurn("Seeding own RNG: %v", seed)
	}

	overmind := ai.NewOvermind(game, config)
	fmt.Printf("%s %s\n", NAME, VERSION)

	var player_strings []string
	for n := 0; n < game.Players(); n++ {
		player_strings = append(player_strings, "bot.exe")
	}

	game.LogWithoutTurn("./halite.exe --width %d --height %d -s %v %s", game.Width(), game.Height(), game.Constants.GameSeed, strings.Join(player_strings, " "))

	// -------------------------------------------------------------------------------

	for {
		game.Parse()

		if config.Timeseed == false {
			rand.Seed(int64(game.Turn() + game.Width() + game.Pid()))
		}

		overmind.Step()
		game.Send()

		if time.Now().Sub(game.ParseTime) > longest_turn {
			longest_turn = time.Now().Sub(game.ParseTime)
			longest_turn_number = game.Turn()
		}
	}
}
