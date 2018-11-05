package core

import (
	"fmt"
	"strings"
	"time"
)

type Game struct {

	Constants
	ParseTime					time.Time

	players						int
	width						int
	height						int

	pid							int				// When simulating, if all sides are being played, this can be set by each bot
	__true_pid					int				// The PID of the player in the real game, regardless of sims; should almost never be read

	logfile						*Logfile
	flogfile					*Flogfile
	token_parser				*TokenParser

	turn						int
	hash						string

	// All of the following are regenerated from scratch each turn...

	budgets						[]int
	halite						[][]int
	ships						[]*Ship			// Each ship contains a command field for the AI to set
	dropoffs					[]*Dropoff		// The first <player_count> items are always the factories
	ship_xy_lookup				map[Point]*Ship
	ship_id_lookup				map[int]*Ship
	box_deltas					map[Point]int
	generate					map[int]bool	// Whether the AI wants to send a "g" command
}

func NewGame() *Game {

	game := new(Game)
	game.turn = -1
	game.token_parser = NewTokenParser()

	return game
}

func (self *Game) set_hash() {

	var s []string

	for _, budget := range self.budgets {
		z := fmt.Sprintf("%d", budget)
		s = append(s, z)
	}

	for x := 0; x < self.width; x++ {
		for y := 0; y < self.height; y++ {
			z := fmt.Sprintf("%d", self.halite[x][y])
			s = append(s, z)
		}
	}

	// To hash the ships and dropoffs consistently we need to do it 1 player at a time...

	for pid := 0; pid < self.players; pid++ {

		ships := self.Ships(pid)

		for _, ship := range ships {
			z := fmt.Sprintf("%d %d %d %d", ship.Owner, ship.X, ship.Y, ship.Halite)		// Don't use ship.Sid, not consistent across engines
			s = append(s, z)
		}

		// FIXME: there's some chance of dropoffs coming in alternate orders.
		// Need to sort dropoffs somehow.

		dropoffs := self.Dropoffs(pid)

		for _, dropoff := range dropoffs {
			z := fmt.Sprintf("%d %d %d", dropoff.Owner, dropoff.X, dropoff.Y)
			s = append(s, z)
		}
	}

	self.hash = HashFromString(strings.Join(s, "-"))
}

func (self *Game) fix_inspiration() {

	for _, ship := range self.ships {

		hits := 0

		for y := 0; y <= self.Constants.INSPIRATION_RADIUS; y++ {

			startx := y - self.Constants.INSPIRATION_RADIUS
			endx := self.Constants.INSPIRATION_RADIUS - y

			for x := startx; x <= endx; x++ {

				other := self.ShipAt(Point{ship.X + x, ship.Y + y})			// Handles bounds automagically
				if other != nil {
					if other.Owner != ship.Owner {
						hits++
					}
				}

				if y != 0 {
					other := self.ShipAt(Point{ship.X + x, ship.Y - y})		// Handles bounds automagically
					if other != nil {
						if other.Owner != ship.Owner {
							hits++
						}
					}
				}
			}
		}

		if hits >= self.Constants.INSPIRATION_SHIP_COUNT {
			ship.Inspired = true
		} else {
			ship.Inspired = false
		}
	}
}
