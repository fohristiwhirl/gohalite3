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
	InitialGroundHalite		int
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

	if self.InitialGroundHalite == 0 {				// This will be at turn 0
		self.InitialGroundHalite = self.Game.GroundHalite()
	}

	self.ClearBook()
	self.UpdatePilots()

	self.FindSwaps()

	for _, pilot := range self.Pilots {
		pilot.SetDesires()
	}

	for _, pilot := range self.Pilots {
		if len(pilot.Desires) == 0 {				// Should be impossible
			self.Game.Log("Pilot %d had no desired!", pilot.Sid)
			pilot.Desires = []string{"o"}
		}
	}

	for _, pilot := range self.Pilots {
		if pilot.Desires[0] == "o" {
			pilot.Ship.Move("o")
			self.SetBook(pilot, pilot)
		}
	}

	for _, pilot := range self.Pilots {

		if pilot.Ship.Command != "" {
			continue
		}

		for _, desire := range pilot.Desires {

			new_loc := pilot.LocationAfterMove(desire)

			if self.Booker(new_loc) == nil {
				pilot.Ship.Move(desire)
				self.SetBook(pilot, new_loc)
				break
			}
		}
	}

	factory := self.Game.MyFactory()
	willing := true

	if self.InitialGroundHalite / (self.Game.GroundHalite() + 1) >= 2 {		// remember int division, also div-by-zero
		willing = false
	}

	if self.Game.Turn() >= self.Game.Constants.MAX_TURNS / 2 {
		willing = false
	}

	if self.Game.MyBudget() >= 1000 && self.Booker(factory) == nil && willing {
		self.Game.SetGenerate(true)
	}

	self.SanityCheck()

	self.Flog()
	return
}

func (self *Overmind) UpdatePilots() {

	// Step 1: add new pilots...

	known_ships := make(map[int]bool)		// Ships that have a pilot already

	for _, pilot := range self.Pilots {
		known_ships[pilot.Sid] = true
	}

	for _, ship := range self.Game.MyShips() {

		if known_ships[ship.Sid] == false {

			pilot := new(Pilot)
			pilot.Game = self.Game
			pilot.Overmind = self
			pilot.Ship = ship
			pilot.Sid = ship.Sid
			pilot.Target = ship.Box()

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

func (self *Overmind) Booker(pos hal.XYer) *Pilot {

	x := hal.Mod(pos.GetX(), self.Game.Width())
	y := hal.Mod(pos.GetY(), self.Game.Height())

	return self.Book[x][y]
}

func (self *Overmind) SetBook(pilot *Pilot, pos hal.XYer) {

	x := hal.Mod(pos.GetX(), self.Game.Width())
	y := hal.Mod(pos.GetY(), self.Game.Height())

	self.Book[x][y] = pilot
}

func (self *Overmind) SanityCheck() {

	targets := make(map[*hal.Box]int)

	for _, pilot := range self.Pilots {
		targetter_sid, ok := targets[pilot.Target]
		if ok && pilot.TargetIsDropoff() == false {
			self.Game.Log("Ships %d and %d looking at same target!", pilot.Sid, targetter_sid)
		} else {
			targets[pilot.Target] = pilot.Sid
		}
	}
}

func (self *Overmind) Flog() {
	for _, pilot := range self.Pilots {
		pilot.Flog()
	}
}

func (self *Overmind) FindSwaps() {

	for cycle := 0; cycle < 4; cycle++ {

		swap_count := 0

		for i, pilot_a := range self.Pilots {

			if pilot_a.TargetIsDropoff() {
				continue
			}

			for _, pilot_b := range self.Pilots[i + 1:] {

				if pilot_b.TargetIsDropoff() {
					continue
				}

				swap_dist := pilot_a.Dist(pilot_b.Target) + pilot_b.Dist(pilot_a.Target)
				curr_dist := pilot_a.Dist(pilot_a.Target) + pilot_b.Dist(pilot_b.Target)

				if swap_dist < curr_dist {
					pilot_a.Target, pilot_b.Target = pilot_b.Target, pilot_a.Target
					self.Game.Log("Swapped targets for pilots %d, %d", pilot_a.Sid, pilot_b.Sid)
					swap_count++
				}
			}
		}

		if swap_count == 0 {
			return
		}
	}
}
