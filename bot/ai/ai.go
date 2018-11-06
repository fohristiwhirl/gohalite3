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
	SimTest					bool
}

type Overmind struct {

	Pid						int

	Config					*Config
	Frame					*hal.Frame				// Needs to be updated every turn
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

func NewOvermind(frame *hal.Frame, config *Config, pid int) *Overmind {

	// At this point, frame has already been pre-pre-parsed and pre-parsed, so the map data exists.

	o := new(Overmind)
	o.Pid = pid
	o.Frame = frame

	o.Config = config
	o.InitialGroundHalite = frame.GroundHalite()

	o.WealthMap = NewWealthMap(o, frame)
	o.DistMap = NewDistMap(o, frame)
	o.EnemyDistMap = NewEnemyDistMap(o, frame)
	o.DropoffDistMap = NewDropoffDistMap(o, frame)
	o.ContestMap = NewContestMap(o, frame)

	return o
}

func (self *Overmind) Step(frame *hal.Frame) {

	// Various calls rely on these two things happening...

	self.Frame = frame
	self.Frame.SetPid(self.Pid)

	rand.Seed(int64(self.Frame.MyBudget() + self.Pid))

	self.WealthMap.Update()
	self.DistMap.Update()
	self.EnemyDistMap.Update()
	self.DropoffDistMap.Update()
	self.ContestMap.Update(self.DistMap, self.EnemyDistMap)

/*
	if self.Frame.Turn() % 100 == 0 {
		self.Frame.Log("Flogging all stats at turn %d", self.Frame.Turn())
		self.WealthMap.Flog()
		self.DistMap.Flog()
		self.EnemyDistMap.Flog()
		self.DropoffDistMap.Flog()
		self.ContestMap.Flog()
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
			self.Frame.Log("Couldn't find a safe move for ship %d (first desire was %s)", pilot.Sid, pilot.Desires[0])
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

	budget := self.Frame.MyBudget()

	factory := self.Frame.MyFactory()
	willing := true

	if self.InitialGroundHalite / (self.Frame.GroundHalite() + 1) >= 2 {		// remember int division, also div-by-zero
		willing = false
	}

	if self.Frame.Turn() >= self.Frame.Constants.MAX_TURNS / 2 {
		willing = false
	}

	if budget >= self.Frame.Constants.NEW_ENTITY_ENERGY_COST && self.MoveBooker(factory) == nil && willing {
		self.Frame.SetGenerate(true)
		budget -= self.Frame.Constants.NEW_ENTITY_ENERGY_COST
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

		if self.Frame.HaliteAt(pilot) == 0 {		// Cheap way to avoid building on enemy dropoff / factory
			continue
		}

		possible_constructs = append(possible_constructs, pilot)
	}

	sort.Slice(possible_constructs, func (a, b int) bool {

		return	self.WealthMap.Values[possible_constructs[a].GetX()][possible_constructs[a].GetY()] <
				self.WealthMap.Values[possible_constructs[b].GetX()][possible_constructs[b].GetY()]
	})

	for _, pilot := range possible_constructs {
		if pilot.Ship.Halite + self.Frame.HaliteAt(pilot) + budget >= self.Frame.Constants.DROPOFF_COST {
			pilot.Ship.Command = "c"
			self.Frame.Log("Ship %d building dropoff (wmap: %d)", pilot.Sid, self.WealthMap.Values[pilot.Ship.X][pilot.Ship.Y])
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

					// self.Frame.Log("Swapped targets for pilots %d, %d (cycle %d)", pilot_a.Sid, pilot_b.Sid, cycle)
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

	for _, ship := range self.Frame.MyShips() {

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

	for x := 0; x < self.Frame.Width(); x++ {
		for y := 0; y < self.Frame.Height(); y++ {
			current_ground_halite += self.Frame.HaliteAtFast(x, y)
		}
	}

	avg_ground_halite := current_ground_halite / (self.Frame.Width() * self.Frame.Height())

	self.HappyThreshold = avg_ground_halite / 2			// Above this, ground is sticky
	self.IgnoreThreshold = avg_ground_halite * 2 / 3	// Less than this not counted for targeting
}

func (self *Overmind) ClearBooks() {

	self.MoveBook = make([][]*Pilot, self.Frame.Width())
	self.TargetBook = make([][]bool, self.Frame.Width())

	for x := 0; x < self.Frame.Width(); x++ {
		self.MoveBook[x] = make([]*Pilot, self.Frame.Height())
		self.TargetBook[x] = make([]bool, self.Frame.Height())
	}
}

func (self *Overmind) MoveBooker(pos hal.XYer) *Pilot {

	x := hal.Mod(pos.GetX(), self.Frame.Width())
	y := hal.Mod(pos.GetY(), self.Frame.Height())

	return self.MoveBook[x][y]
}

func (self *Overmind) SetMoveBook(pilot *Pilot, pos hal.XYer) {

	x := hal.Mod(pos.GetX(), self.Frame.Width())
	y := hal.Mod(pos.GetY(), self.Frame.Height())

	self.MoveBook[x][y] = pilot
}

func (self *Overmind) SameTargetCheck() {

	targets := make(map[hal.Point]int)

	for _, pilot := range self.Pilots {
		targetter_sid, ok := targets[pilot.Target]
		if ok && pilot.TargetIsDropoff() == false {
			self.Frame.Log("Ships %d and %d looking at same target: %d %d", pilot.Sid, targetter_sid, pilot.Target.X, pilot.Target.Y)
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

	pos_halite := self.Frame.HaliteAt(pos)
	tar_halite := self.Frame.HaliteAt(tar)

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
