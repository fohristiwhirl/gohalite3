package core

import (
	"fmt"
	"strings"
	"time"
)

type Frame struct {

	Constants
	ParseTime					time.Time

	players						int
	width						int
	height						int

	pid							int				// When simulating, if all sides are being played, this can be set by each bot
	__true_pid					int				// The PID of the player in the real game, regardless of sims; should almost never be read

	initial_ground_halite		int

	token_parser				*TokenParser

	turn						int
	highest_sid_seen			int							// Mostly for the simulator, which needs to generate unique new sids

	// All of the following are regenerated from scratch each turn...

	budgets						[]int
	halite						[][]int
	ships						[]*Ship						// Each ship contains a command field for the AI to set
	dropoffs					[]*Dropoff					// The first <player_count> items are always the factories
	ship_xy_lookup				map[Point]*Ship
	ship_id_lookup				map[int]*Ship
	generate					map[int]bool				// Whether the AI wants to send a "g" command

	// The following are cleared each parse / simgen / remake, then made when asked for, and cached...

	wealth_map					*WealthMap
	inspiration_map				map[int]*InspirationMap		// Likewise
	ground_halite				int							// Likewise
}

func NewGame() *Frame {

	frame := new(Frame)
	frame.turn = -1
	frame.highest_sid_seen = -1
	frame.token_parser = NewTokenParser()

	return frame
}

func (self *Frame) Zerofy() {

	// Put all the real content into a zeroed state.
	// The caller can then update it.

	self.budgets = make([]int, self.players)
	self.halite = Make2dIntArray(self.width, self.height)
	self.ships = nil
	self.dropoffs = nil
	self.ship_xy_lookup = make(map[Point]*Ship)
	self.ship_id_lookup = make(map[int]*Ship)
	self.generate = make(map[int]bool)

	// Cached stuff that gets remade when asked for...

	self.wealth_map = nil
	self.inspiration_map = nil
	self.ground_halite = 0
}

func (self *Frame) Remake() *Frame {			// This is a deep copy

	g := new(Frame)
	*g = *self			// Everything not explicitly changed will be the same

	g.ParseTime = time.Now()

	g.Zerofy()			// Clear all the data!

	for pid, val := range self.budgets {
		g.budgets[pid] = val
	}

	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			g.halite[x][y] = self.halite[x][y]
		}
	}

	for _, ship := range self.ships {
		remade := *ship
		remade.Frame = g
		g.ships = append(g.ships, &remade)
	}

	for _, dropoff := range self.dropoffs {
		remade := *dropoff
		remade.Frame = g
		g.dropoffs = append(g.dropoffs, &remade)
	}

	for _, ship := range g.ships {
		g.ship_xy_lookup[Point{ship.X, ship.Y}] = ship
		g.ship_id_lookup[ship.Sid] = ship
	}

	for key, val := range self.generate {
		g.generate[key] = val
	}

	return g
}

func (self *Frame) Hash() string {

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

	return HashFromString(strings.Join(s, "-"))
}

func (self *Frame) FixInspiration() {

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
