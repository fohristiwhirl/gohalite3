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

func main() {

	const (
		NAME = "Fohristiwhirl"
		VERSION = "16.b"				// hash is ??
	)

	config := new(ai.Config)

	flag.BoolVar(&config.Crash, "crash", false, "randomly crash")
	flag.BoolVar(&config.SimTest, "simtest", false, "test the simulator")
	flag.Parse()

	game := hal.NewGame()

	var longest_turn time.Duration
	var longest_turn_number int

	start_time := time.Now()

	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("%v", p)
			game.Log("Quitting: %v", p)
			game.Log("Last known hash: %s", game.Hash())
			game.Log("Longest turn (%d) took %v", longest_turn_number, longest_turn)
			game.Log("Real-world time elapsed: %v", time.Now().Sub(start_time))
			game.StopLog()
			game.StopFlog()
		}
	}()

	// -------------------------------------------------------------------------------

	game.PrePreParse()				// Reads very early data (including PID, needed for log names)
	true_pid := game.Pid()

	// Both of these fail harmlessly if the directory isn't there:
	game.StartLog(fmt.Sprintf("logs/log-%v.txt", true_pid))
	game.StartFlog(fmt.Sprintf("flogs/flog-%v-%v.json", game.Constants.GameSeed, true_pid))

	game.PreParse()					// Reads the map data.

	game.LogWithoutTurn("--------------------------------------------------------------------------------")
	game.LogWithoutTurn("%s %s starting up at %s", NAME, VERSION, time.Now().Format("2006-01-02 15:04:05"))

	overmind := ai.NewOvermind(game, config, true_pid)

	var player_strings []string
	for n := 0; n < game.Players(); n++ {
		player_strings = append(player_strings, "bot.exe")
	}

	game.LogWithoutTurn("./halite.exe --width %d --height %d -s %v %s", game.Width(), game.Height(), game.Constants.GameSeed, strings.Join(player_strings, " "))

	// -------------------------------------------------------------------------------

	if config.SimTest {
		prediction := sim_check(game, config)
		game.Log("Simulator predicts final hash %s", prediction)
	}

	fmt.Printf("%s %s\n", NAME, VERSION)

	for {
		game.Parse()

		if config.Crash {
			if rand.Intn(100) == 40 {
				fmt.Printf("g g\n")
			} else if rand.Intn(100) == 40 {
				time.Sleep(5 * time.Second)
			}
		}

		overmind.Step(game)
		game.Send()

		if time.Now().Sub(game.ParseTime) > longest_turn {
			longest_turn = time.Now().Sub(game.ParseTime)
			longest_turn_number = game.Turn()
		}
	}
}

