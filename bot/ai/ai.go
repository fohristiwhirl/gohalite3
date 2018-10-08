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

func NewOvermind(game *hal.Game, config *Config) *Overmind {
	ret := new(Overmind)
	ret.Game = game
	ret.Config = config
	return ret
}

func (self *Overmind) Step() {

	self.UpdatePilots()

	for _, pilot := range self.Pilots {
		pilot.Fly()
	}

	if self.Game.MyBudget() >= 1000 {
		self.Game.SetGenerate(true)
	}

	return
}

func (self *Overmind) UpdatePilots() {

	// Step 1: add new pilots...

	known_ships := make(map[int]bool)		// Ships that have a pilot already

	for _, pilot := range self.Pilots {
		known_ships[pilot.Sid] = true
	}

	my_ships := self.Game.MyShips()

	for _, ship := range my_ships {

		if known_ships[ship.Sid] == false {

			pilot := new(Pilot)
			pilot.Game = self.Game
			pilot.Ship = ship
			pilot.Sid = ship.Sid
			pilot.State = Normal
			pilot.TargetX = ship.X
			pilot.TargetY = ship.Y

			self.Pilots = append(self.Pilots, pilot)

			self.Game.Log("New pilot: %d", ship.Sid)
		}
	}

	// Step 2: delete dead pilots...

	for n := len(self.Pilots) - 1; n >= 0 ; n-- {
		pilot := self.Pilots[n]
		if pilot.Ship.Alive == false {
			self.Pilots = append(self.Pilots[:n], self.Pilots[n+1:]...)
			self.Game.Log("Dead pilot: %d", pilot.Ship.Sid)
		}
	}
}
