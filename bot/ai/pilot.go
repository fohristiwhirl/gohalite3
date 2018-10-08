package ai

import (
	"fmt"
	"math/rand"

	hal "../core"
)

type Pilot struct {
	Ship					*hal.Ship
	Sid						int
}

func (self *Pilot) Navigate(x, y int) {

	// FIXME: consider wraps

	ship := self.Ship

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
