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
	Book					[][]*Pilot
}

func NewOvermind(game *hal.Game, config *Config) *Overmind {

	o := new(Overmind)
	o.Game = game
	o.Config = config

	return o
}

func (self *Overmind) Step() {

	self.ClearBook()
	self.UpdatePilots()

	for _, pilot := range self.Pilots {
		if pilot.Ship.Command == "" {
			pilot.SetDesires()
		}
	}

	for _, pilot := range self.Pilots {
		if len(pilot.Desires) == 0 {				// Should be impossible
			pilot.Desires = []string{"o"}
		}
	}

	for _, pilot := range self.Pilots {
		if pilot.Desires[0] == "o" {
			pilot.Ship.Move("o")
			self.SetBook(pilot, pilot.Ship.X, pilot.Ship.Y)
		}
	}

	for _, pilot := range self.Pilots {

		if pilot.Ship.Command != "" {
			continue
		}

		for _, desire := range pilot.Desires {

			new_x, new_y := pilot.LocationAfterMove(desire)

			if self.Booker(new_x, new_y) == nil {
				pilot.Ship.Move(desire)
				self.SetBook(pilot, new_x, new_y)
				break
			}
		}
	}

	factory := self.Game.MyFactory()

	if self.Game.MyBudget() >= 1000 && self.Booker(factory.X, factory.Y) == nil && self.Game.Turn() < self.Game.Constants.MAX_TURNS / 2 {
		self.Game.SetGenerate(true)
	}

	self.SanityCheck()

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
			pilot.Overmind = self
			pilot.Ship = ship
			pilot.Sid = ship.Sid
			pilot.State = Normal
			pilot.Target = self.Game.Box(ship.X, ship.Y)

			self.Pilots = append(self.Pilots, pilot)
		}
	}

	// Step 2: delete dead pilots...

	for n := len(self.Pilots) - 1; n >= 0 ; n-- {
		pilot := self.Pilots[n]
		if pilot.Ship.Alive == false {
			self.Pilots = append(self.Pilots[:n], self.Pilots[n+1:]...)
		}
	}

	// Step 3: other maintainence...

	for _, pilot := range self.Pilots {
		pilot.Desires = nil
	}
}

func (self *Overmind) ClearBook() {
	self.Book = make([][]*Pilot, self.Game.Width())
	for x := 0; x < self.Game.Width(); x++ {
		self.Book[x] = make([]*Pilot, self.Game.Height())
	}
}

func (self *Overmind) Booker(x, y int) *Pilot {

	x = mod(x, self.Game.Width())
	y = mod(y, self.Game.Height())

	return self.Book[x][y]
}

func (self *Overmind) SetBook(pilot *Pilot, x, y int) {

	x = mod(x, self.Game.Width())
	y = mod(y, self.Game.Height())

	self.Book[x][y] = pilot
}

func (self *Overmind) SanityCheck() {

	targets := make(map[hal.Point]bool)

	for _, pilot := range self.Pilots {
		if pilot.State == Normal {
			if targets[hal.Point{pilot.Target.GetX(), pilot.Target.GetY()}] {
				self.Game.Log("Multiple \"Normal\" ships looking at same target")
			} else {
				targets[hal.Point{pilot.Target.GetX(), pilot.Target.GetY()}] = true
			}
		}
	}
}
