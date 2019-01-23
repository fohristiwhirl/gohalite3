package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"./ai"
	"./config"
	"./logging"
	hal "./core"
)

func main() {

	const (
		NAME = "Fohristiwhirl"
		VERSION = "27b"
	)

	config.ParseCommandLine()

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
			logging.StopLog()
			logging.StopFlog()
		}
	}()

	// -------------------------------------------------------------------------------

	frame.PrePreParse()				// Reads very early data (including PID, needed for log names)
	true_pid := frame.Pid()

	// Both of these fail harmlessly if the directory isn't there:
	logging.StartLog(fmt.Sprintf("logs/log-%v.txt", true_pid))
	logging.StartFlog(fmt.Sprintf("flogs/flog-%v-%v.json", frame.Constants.GameSeed, true_pid))

	frame.PreParse()				// Reads the map data.

	config.SetGenMin(frame.Width(), frame.Players())

	logging.Log("--------------------------------------------------------------------------------")
	logging.Log("%s %s starting up at %s", NAME, VERSION, time.Now().Format("2006-01-02 15:04:05"))
	logging.Log("Invoked as %s", strings.Join(os.Args, " "))

	var player_strings []string
	for n := 0; n < frame.Players(); n++ {
		player_strings = append(player_strings, "bot.exe")
	}

	logging.Log("./halite.py --width %d --height %d -s %v %s",
		frame.Width(), frame.Height(), frame.Constants.GameSeed, strings.Join(player_strings, " "))

	logging.Log("GenMin is %v", config.GenMin)

	// -------------------------------------------------------------------------------

	if config.SimTest {
		logging.Suppress()
		prediction_hash, prediction_ground := sim_check(frame)
		logging.Allow()
		logging.Log("Simulator predicts final hash %v", prediction_hash)
		logging.Log("Simulator predicts ground halite %v on turn N-1", prediction_ground)
	}

	fmt.Printf("%s %s\n", NAME, VERSION)

	for {
		frame.Parse()

		if config.RemakeTest {
			frame = frame.Remake()
		}

		if config.Crash {
			if rand.Intn(100) == 40 {
				fmt.Printf("g g\n")
			} else if rand.Intn(100) == 40 {
				time.Sleep(5 * time.Second)
			}
		}

		ai.Step(frame, true_pid, true)
		frame.Send()

		if time.Now().Sub(frame.ParseTime) > longest_turn {
			longest_turn = time.Now().Sub(frame.ParseTime)
			longest_turn_number = frame.Turn()
		}
	}
}

func sim_check(real_frame *hal.Frame) (string, int) {

	// Returns the final hash that the real bot will see,
	// if the real bot is matched only against itself...

	frame := real_frame.Remake()

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
