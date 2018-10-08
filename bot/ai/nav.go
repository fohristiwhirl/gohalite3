package ai

import (
	"fmt"
	"math/rand"

	hal "../core"
)

func abs(a int) int {
	if a < 0 { return -a }
	return a
}

func (self *Overmind) Navigate(ship *hal.Ship, x, y int) {

	// FIXME: consider wraps

	dx := x - ship.X
	dy := y - ship.Y

	if dx == 0 && dy == 0 {
		ship.Clear()
		return
	}

	var options []string

	if dx > 0 {
		options = append(options, "e")
	}

	if dx < 0 {
		options = append(options, "w")
	}

	if dy > 0 {
		options = append(options, "s")
	}

	if dy < 0 {
		options = append(options, "n")
	}

	n := rand.Intn(len(options))

	direction := options[n]

	command := fmt.Sprintf("m %d s", ship.Sid, direction)
	ship.Set(command)
}
