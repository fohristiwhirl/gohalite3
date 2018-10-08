package ai

import (
	hal "../core"
)

type Config struct {
	Timeseed				bool
}

type Overmind struct {
	Config					*Config
	Game					*hal.Game
	Pilots					[]*Pilot
}

type Pilot struct {
	Ship					*hal.Ship
	Sid						int
}

func NewOvermind(game *hal.Game, config *Config) *Overmind {
	ret := new(Overmind)
	ret.Game = game
	ret.Config = config
	return ret
}

func (self *Overmind) Step() {

	// budget := self.Game.MyBudget()

	known_ships := make(map[int]bool)		// Ships that have a pilot already

	for _, pilot := range self.Pilots {
		known_ships[pilot.Sid] = true
	}

	my_ships := self.Game.MyShips()

	for _, ship := range my_ships {
		if known_ships[ship.Sid] == false {
			pilot := new(Pilot)
			pilot.Ship = ship
			pilot.Sid = ship.Sid
			self.Pilots = append(self.Pilots, pilot)
		}
	}

	return
}
