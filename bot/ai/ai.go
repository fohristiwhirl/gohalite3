package ai

import (
	"sort"

	hal "../core"
)

const (
	DROPOFF_SPACING = 12
	NICE_THRESHOLD = 8000
)

type Config struct {
	Timeseed				bool
}

type Overmind struct {

	Config					*Config
	Game					*hal.Game
	Pilots					[]*Pilot

	// ATC stuff:

	TargetBook				[][]bool
	MoveBook				[][]*Pilot

	// Stategic stats:

	NiceMap					*NiceMap

	InitialGroundHalite		int
	HappyThreshold			int
}

func NewOvermind(game *hal.Game, config *Config) *Overmind {

	// At this point, game has already been pre-pre-parsed and pre-parsed, so the map data exists.

	o := new(Overmind)
	o.Game = game

	o.Config = config
	o.InitialGroundHalite = game.GroundHalite()
	o.NiceMap = NewNiceMap(game)

	return o
}

func (self *Overmind) Step() {

	self.NiceMap.Update()
	self.InspectGround()
	self.ClearBooks()
	self.UpdatePilots()

	// What each ship wants to do............................

	for _, pilot := range self.Pilots {
		pilot.SetTarget()
	}

	self.TargetSwaps()

	for _, pilot := range self.Pilots {
		pilot.SetDesires()
	}

	// Resolve the desired moves.............................

	for _, pilot := range self.Pilots {
		if pilot.Desires[0] == "o" {
			pilot.Ship.Move("o")
			self.SetMoveBook(pilot, pilot)
		}
	}

	for cycle := 0; cycle < 5; cycle++ {

		for _, pilot := range self.Pilots {

			if pilot.Ship.Command != "" {
				continue
			}

			// Special case: if ship is next to a dropoff and is in its mad dash, always move.
			// And don't set the book, it can only confuse matters...

			if pilot.TargetIsDropoff() && pilot.Dist(pilot.Target) == 1 && pilot.FinalDash() {
				pilot.Ship.Move(pilot.Desires[0])
				continue
			}

			// Normal case...

			for _, desire := range pilot.Desires {

				new_loc := pilot.LocationAfterMove(desire)
				booker := self.MoveBooker(new_loc)

				if booker == nil {
					pilot.Ship.Move(desire)
					self.SetMoveBook(pilot, new_loc)
					break
				} else {
					if booker.Ship.Command == "o" {		// Never clear a booking made by a stationary ship
						continue
					}
					if booker.Ship.Halite < pilot.Ship.Halite {
						pilot.Ship.Move(desire)
						self.SetMoveBook(pilot, new_loc)
						booker.Ship.ClearMove()
						break
					}
				}
			}
		}
	}

	for _, pilot := range self.Pilots {
		if pilot.Ship.Command == "" {
			self.Game.Log("Couldn't find a safe move for ship %d (first desire was %s)", pilot.Sid, pilot.Desires[0])
		}
	}

	// Other.................................................

	self.MaybeBuild()

	// FIXME: re-add the sanity checks.

	self.Flog()
	return
}

func (self *Overmind) MaybeBuild() {

	budget := self.Game.MyBudget()

	factory := self.Game.MyFactory()
	willing := true

	if self.InitialGroundHalite / (self.Game.GroundHalite() + 1) >= 2 {		// remember int division, also div-by-zero
		willing = false
	}

	if self.Game.Turn() >= self.Game.Constants.MAX_TURNS / 2 {
		willing = false
	}

	if budget >= self.Game.Constants.NEW_ENTITY_ENERGY_COST && self.MoveBooker(factory) == nil && willing {
		self.Game.SetGenerate(true)
		budget -= self.Game.Constants.NEW_ENTITY_ENERGY_COST
	}

	// -------------------------------------------

	var possible_constructs []*Pilot

	for _, pilot := range self.Pilots {

		if pilot.Dist(pilot.NearestDropoff()) < DROPOFF_SPACING {
			continue
		}

		if self.NiceMap.Values[pilot.GetX()][pilot.GetY()] < NICE_THRESHOLD {
			continue
		}

		if pilot.Box().Halite == 0 {		// Cheap way to avoid building on enemy dropoff / factory
			continue
		}

		possible_constructs = append(possible_constructs, pilot)
	}

	sort.Slice(possible_constructs, func (a, b int) bool {

		return	self.NiceMap.Values[possible_constructs[a].GetX()][possible_constructs[a].GetY()] <
				self.NiceMap.Values[possible_constructs[b].GetX()][possible_constructs[b].GetY()]
	})

	for _, pilot := range possible_constructs {
		if pilot.Ship.Halite + pilot.Box().Halite + budget >= self.Game.Constants.DROPOFF_COST {
			pilot.Ship.Command = "c"
			break
		}
	}
}

func (self *Overmind) TargetSwaps() {

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

				a_dist_b := pilot_a.Dist(pilot_b.Target)
				b_dist_a := pilot_b.Dist(pilot_a.Target)

				alt_score_a := halite_dist_score(pilot_b.Target.Halite, a_dist_b)
				alt_score_b := halite_dist_score(pilot_a.Target.Halite, b_dist_a)

				if alt_score_a + alt_score_b > pilot_a.Score + pilot_b.Score {

					pilot_a.Target, pilot_b.Target = pilot_b.Target, pilot_a.Target

					pilot_a.Score = alt_score_a
					pilot_b.Score = alt_score_b

					self.Game.Log("Swapped targets for pilots %d, %d (cycle %d)", pilot_a.Sid, pilot_b.Sid, cycle)
					swap_count++
				}
			}
		}

		if swap_count == 0 {
			return
		}
	}
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
		pilot.Target = pilot.Box()
		pilot.Score = 0
	}
}

func (self *Overmind) InspectGround() {		// Old behaviour used to be simply self.HappyThreshold = 50

	current_ground_halite := 0

	for x := 0; x < self.Game.Width(); x++ {
		for y := 0; y < self.Game.Height(); y++ {
			current_ground_halite += self.Game.BoxAtFast(x, y).Halite
		}
	}

	avg_ground_halite := current_ground_halite / (self.Game.Width() * self.Game.Height())

	self.HappyThreshold = avg_ground_halite / 2
}

func (self *Overmind) ClearBooks() {

	self.MoveBook = make([][]*Pilot, self.Game.Width())
	self.TargetBook = make([][]bool, self.Game.Width())

	for x := 0; x < self.Game.Width(); x++ {
		self.MoveBook[x] = make([]*Pilot, self.Game.Height())
		self.TargetBook[x] = make([]bool, self.Game.Height())
	}
}

func (self *Overmind) MoveBooker(pos hal.XYer) *Pilot {

	x := hal.Mod(pos.GetX(), self.Game.Width())
	y := hal.Mod(pos.GetY(), self.Game.Height())

	return self.MoveBook[x][y]
}

func (self *Overmind) SetMoveBook(pilot *Pilot, pos hal.XYer) {

	x := hal.Mod(pos.GetX(), self.Game.Width())
	y := hal.Mod(pos.GetY(), self.Game.Height())

	self.MoveBook[x][y] = pilot
}

func (self *Overmind) Flog() {
	for _, pilot := range self.Pilots {
		pilot.Flog()
	}
}

func halite_dist_score(halite, dist int) float32 {
	return float32(halite) / float32((dist + 1) * (dist + 1))	// Avoid div-by-zero
}
