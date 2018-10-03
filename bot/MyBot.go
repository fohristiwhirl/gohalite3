package main

import (
	"fmt"
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

	// Do flag parsing here.

	game := hal.NewGame()

	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("%v", p)
			game.Log("Quitting: %v", p)
		}
	}()

	game.StartLog(fmt.Sprintf("log%d.txt", game.Pid()))
	game.LogWithoutTurn("--------------------------------------------------------------------------------")
	game.LogWithoutTurn("%s %s starting up at %s", NAME, VERSION, time.Now().Format("2006-01-02 15:04:05"))

	overmind := ai.NewOvermind(game, config)

	// Play the game...

	overmind.Step()
}
