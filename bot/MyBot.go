package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
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
	flag.BoolVar(&config.RemakeTest, "remaketest", false, "test the frame remaker")
	flag.BoolVar(&config.SimTest, "simtest", false, "test the simulator")
	flag.Parse()

	frame := hal.NewGame()

	var longest_turn time.Duration
	var longest_turn_number int

	start_time := time.Now()

	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("%v", p)
			frame.Log("Quitting: %v", p)
			frame.Log("Last known hash: %s", frame.Hash())
			frame.Log("Longest turn (%d) took %v", longest_turn_number, longest_turn)
			frame.Log("Real-world time elapsed: %v", time.Now().Sub(start_time))
			frame.StopLog()
			frame.StopFlog()
		}
	}()

	// -------------------------------------------------------------------------------

	frame.PrePreParse()				// Reads very early data (including PID, needed for log names)
	true_pid := frame.Pid()

	// Both of these fail harmlessly if the directory isn't there:
	frame.StartLog(fmt.Sprintf("logs/log-%v.txt", true_pid))
	frame.StartFlog(fmt.Sprintf("flogs/flog-%v-%v.json", frame.Constants.GameSeed, true_pid))

	frame.PreParse()				// Reads the map data.

	frame.LogWithoutTurn("--------------------------------------------------------------------------------")
	frame.LogWithoutTurn("%s %s starting up at %s", NAME, VERSION, time.Now().Format("2006-01-02 15:04:05"))
	frame.LogWithoutTurn("Invoked as %s", strings.Join(os.Args, " "))

	overmind := ai.NewOvermind(frame, config, true_pid)

	var player_strings []string
	for n := 0; n < frame.Players(); n++ {
		player_strings = append(player_strings, "bot.exe")
	}

	frame.LogWithoutTurn("./halite.exe --width %d --height %d -s %v %s",
		frame.Width(), frame.Height(), frame.Constants.GameSeed, strings.Join(player_strings, " "))

	// -------------------------------------------------------------------------------

	if config.SimTest {
		prediction_hash, prediction_ground := sim_check(frame, config)
		frame.Log("Simulator predicts final hash %v", prediction_hash)
		frame.Log("Simulator predicts ground halite %v on turn N-1", prediction_ground)
	}

	fmt.Printf("%s %s\n", NAME, VERSION)

	for {
		frame.Parse()

		if config.RemakeTest {
			frame = frame.Remake(true)
		}

		if config.Crash {
			if rand.Intn(100) == 40 {
				fmt.Printf("g g\n")
			} else if rand.Intn(100) == 40 {
				time.Sleep(5 * time.Second)
			}
		}

		overmind.Step(frame)
		frame.Send()

		if time.Now().Sub(frame.ParseTime) > longest_turn {
			longest_turn = time.Now().Sub(frame.ParseTime)
			longest_turn_number = frame.Turn()
		}
	}
}

