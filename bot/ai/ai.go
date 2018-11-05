package ai

import (
	"math/rand"
	"sort"

	hal "../core"
)

const (
	DROPOFF_SPACING = 12
	NICE_THRESHOLD = 8000
)

type Config struct {
	Crash					bool
}

type Overmind struct {

	Pid						int

	Config					*Config
	Game					*hal.Game				// Needs to be updated every turn for sims
	Pilots					[]*Pilot

	// ATC stuff:

	TargetBook				[][]bool
	MoveBook				[][]*Pilot

	// Stategic stats:

	WealthMap				*WealthMap
	DistMap					*DistMap
	EnemyDistMap			*EnemyDistMap
	DropoffDistMap			*DropoffDistMap
	ContestMap				*ContestMap

	InitialGroundHalite		int
	HappyThreshold			int
	IgnoreThreshold			int
}

func NewOvermind(game *hal.Game, config *Config, pid int) *Overmind {

	// At this point, game has already been pre-pre-parsed and pre-parsed, so the map data exists.

	o := new(Overmind)
	o.Pid = pid
	o.Game = game

	o.Config = config
	o.InitialGroundHalite = game.GroundHalite()

	o.WealthMap = NewWealthMap(game)
	o.DistMap = NewDistMap(game)
	o.EnemyDistMap = NewEnemyDistMap(game)
	o.DropoffDistMap = NewDropoffDistMap(game)
	o.ContestMap = NewContestMap(game)

	return o
}

func (self *Overmind) Step(game *hal.Game) {

	self.Game = game
	self.Game.SetPid(self.Pid)

	rand.Seed(int64(self.Game.MyBudget() + self.Pid))

	self.WealthMap.Update(game)
	self.DistMap.Update(game)
	self.EnemyDistMap.Update(game)
	self.DropoffDistMap.Update(game)
	self.ContestMap.Update(game, self.DistMap, self.EnemyDistMap)

/*
	if self.Game.Turn() % 100 == 0 {
		self.Game.Log("Flogging all stats at turn %d", self.Game.Turn())
		self.WealthMap.Flog(game)
		self.DistMap.Flog(game)
		self.EnemyDistMap.Flog(game)
		self.DropoffDistMap.Flog(game)
		self.ContestMap.Flog(game)
	}
*/

	self.SetTurnParameters()
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

	for _, pilot := range self.Pilots {
		pilot.FlogTarget()
	}

	self.SameTargetCheck()		// Just logs
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

		if self.WealthMap.Values[pilot.GetX()][pilot.GetY()] < NICE_THRESHOLD {
			continue
		}

		if self.Game.HaliteAt(pilot) == 0 {		// Cheap way to avoid building on enemy dropoff / factory
			continue
		}

		possible_constructs = append(possible_constructs, pilot)
	}

	sort.Slice(possible_constructs, func (a, b int) bool {

		return	self.WealthMap.Values[possible_constructs[a].GetX()][possible_constructs[a].GetY()] <
				self.WealthMap.Values[possible_constructs[b].GetX()][possible_constructs[b].GetY()]
	})

	for _, pilot := range possible_constructs {
		if pilot.Ship.Halite + self.Game.HaliteAt(pilot) + budget >= self.Game.Constants.DROPOFF_COST {
			pilot.Ship.Command = "c"
			self.Game.Log("Ship %d building dropoff (wmap: %d)", pilot.Sid, self.WealthMap.Values[pilot.Ship.X][pilot.Ship.Y])
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

				alt_score_a := halite_dist_score(pilot_b.TargetHalite(), a_dist_b)
				alt_score_b := halite_dist_score(pilot_a.TargetHalite(), b_dist_a)

				if alt_score_a + alt_score_b > pilot_a.Score + pilot_b.Score {

					pilot_a.Target, pilot_b.Target = pilot_b.Target, pilot_a.Target

					pilot_a.Score = alt_score_a
					pilot_b.Score = alt_score_b

					// self.Game.Log("Swapped targets for pilots %d, %d (cycle %d)", pilot_a.Sid, pilot_b.Sid, cycle)
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

	var live_pilots []*Pilot		// Pilots with valid ships this turn

	sid_pilot_map := make(map[int]*Pilot)

	for _, pilot := range self.Pilots {
		sid_pilot_map[pilot.Sid] = pilot
	}

	for _, ship := range self.Game.MyShips() {

		pilot, ok := sid_pilot_map[ship.Sid]

		if ok {

			pilot.Ship = ship

			live_pilots = append(live_pilots, pilot)

		} else {

			pilot = new(Pilot)
			pilot.Overmind = self
			pilot.Ship = ship
			pilot.Sid = ship.Sid

			live_pilots = append(live_pilots, pilot)
		}
	}

	for _, pilot := range live_pilots {
		pilot.NewTurn()
	}

	self.Pilots = live_pilots
}

func (self *Overmind) SetTurnParameters() {

	current_ground_halite := 0

	for x := 0; x < self.Game.Width(); x++ {
		for y := 0; y < self.Game.Height(); y++ {
			current_ground_halite += self.Game.HaliteAtFast(x, y)
		}
	}

	avg_ground_halite := current_ground_halite / (self.Game.Width() * self.Game.Height())

	self.HappyThreshold = avg_ground_halite / 2			// Above this, ground is sticky
	self.IgnoreThreshold = avg_ground_halite * 2 / 3	// Less than this not counted for targeting
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

func (self *Overmind) SameTargetCheck() {

	targets := make(map[hal.Point]int)

	for _, pilot := range self.Pilots {
		targetter_sid, ok := targets[pilot.Target]
		if ok && pilot.TargetIsDropoff() == false {
			self.Game.Log("Ships %d and %d looking at same target: %d %d", pilot.Sid, targetter_sid, pilot.Target.X, pilot.Target.Y)
		} else {
			targets[pilot.Target] = pilot.Sid
		}
	}
}

func (self *Overmind) ShouldMine(halite_carried int, pos, tar hal.XYer) bool {

	// Whether a ship -- if it were carrying n halite, at pos, with specified target -- would stop to mine.

	if halite_carried >= 800 {
		return false
	}

	pos_halite := self.Game.HaliteAt(pos)
	tar_halite := self.Game.HaliteAt(tar)

	if pos_halite > self.HappyThreshold {
		if pos_halite > tar_halite / 3 {			// This is a bit odd since the test even happens when target is dropoff.
			return true
		}
	}

	return false
}

func halite_dist_score(halite, dist int) float32 {
	return float32(halite) / float32((dist + 1) * (dist + 1))	// Avoid div-by-zero
}
