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
}

func NewOvermind(game *hal.Game, config *Config) *Overmind {
	ret := new(Overmind)
	ret.Game = game
	ret.Config = config

	return ret
}

func (self *Overmind) Step() {

	my_ships := self.Game.MyShips()
	budget := self.Game.MyBudget()

	if budget >= 1000 && len(my_ships) < 4 {
		self.Game.SetGenerate(true)
	}

	for n, ship := range my_ships {

		switch n % 4 {

		case 0:

			if ship.Halite == 0 {
				ship.Left()
			} else if ship.Halite > 200 {
				ship.Right()
			} else {
				if self.Game.Turn() % 2 == 0 {
					ship.Left()
				}
			}

		case 1:

			if ship.Halite == 0 {
				ship.Up()
			} else if ship.Halite > 200 {
				ship.Down()
			} else {
				if self.Game.Turn() % 2 == 0 {
					ship.Up()
				}
			}

		case 2:

			if ship.Halite == 0 {
				ship.Right()
			} else if ship.Halite > 200 {
				ship.Left()
			} else {
				if self.Game.Turn() % 2 == 0 {
					ship.Right()
				}
			}

		case 3:

			if ship.Halite == 0 {
				ship.Down()
			} else if ship.Halite > 200 {
				ship.Up()
			} else {
				if self.Game.Turn() % 2 == 0 {
					ship.Down()
				}
			}

		}
	}

	return
}
