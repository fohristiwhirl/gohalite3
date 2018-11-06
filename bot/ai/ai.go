package ai

import (
	"math/rand"
	"sort"

	hal "../core"
	maps "../maps"
)

const (
	DROPOFF_SPACING = 12
	NICE_THRESHOLD = 8000
)

type Config struct {
	Crash					bool
	RemakeTest				bool
	SimTest					bool
}

type Overmind struct {

	Pid						int
	Config					*Config

	Frame					*hal.Frame				// Needs to be updated every turn

	// Stategic stats:

	WealthMap				*maps.WealthMap			// Other maps are available in /maps

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

	o.WealthMap = maps.NewWealthMap(frame)

	return o
}

func (self *Overmind) Step(frame *hal.Frame) {

	// Various calls rely on these two things happening......

	self.Frame = frame
	self.Frame.SetPid(self.Pid)

	// Various other initialisation..........................

	rand.Seed(int64(self.Frame.MyBudget() + self.Pid))
	self.WealthMap.Update(frame)
	self.SetTurnParameters()

	// What each ship wants to do............................

	my_ships := self.Frame.MyShips()

	for _, ship := range my_ships {
		self.NewTurn(ship)
	}

	target_book := hal.Make2dBoolArray(frame.Width(), frame.Height())		// What points are targets. Updated for each ship.

	for _, ship := range my_ships {
		self.SetTarget(ship, target_book)
	}

	self.TargetSwaps(my_ships)

	for _, ship := range my_ships {
		self.SetDesires(ship)
	}

	// Resolve the desired moves.............................

	move_book := Resolve(frame, my_ships)

	// Other.................................................

	self.MaybeBuild(my_ships, move_book)

	for _, ship := range my_ships {
		self.FlogTarget(ship)
	}

	self.SameTargetCheck(my_ships)		// Just logs
	return
}

func (self *Overmind) MaybeBuild(my_ships []*hal.Ship, move_book *MoveBook) {

	budget := self.Frame.MyBudget()

	factory := self.Frame.MyFactory()
	willing := true

	if self.InitialGroundHalite / (self.Frame.GroundHalite() + 1) >= 2 {		// remember int division, also div-by-zero
		willing = false
	}

	if self.Frame.Turn() >= self.Frame.Constants.MAX_TURNS / 2 {
		willing = false
	}

	if budget >= self.Frame.Constants.NEW_ENTITY_ENERGY_COST && move_book.Booker(factory) == nil && willing {
		self.Frame.SetGenerate(true)
		budget -= self.Frame.Constants.NEW_ENTITY_ENERGY_COST
	}

	// -------------------------------------------

	var possible_constructs []*hal.Ship

	for _, ship := range my_ships {

		if ship.Dist(ship.NearestDropoff()) < DROPOFF_SPACING {
			continue
		}

		if self.WealthMap.Values[ship.X][ship.Y] < NICE_THRESHOLD {
			continue
		}

		if self.Frame.HaliteAtFast(ship.X, ship.Y) == 0 {		// Cheap way to avoid building on enemy dropoff / factory
			continue
		}

		possible_constructs = append(possible_constructs, ship)
	}

	sort.Slice(possible_constructs, func (a, b int) bool {

		return	self.WealthMap.Values[possible_constructs[a].X][possible_constructs[a].Y] <
				self.WealthMap.Values[possible_constructs[b].X][possible_constructs[b].Y]
	})

	for _, ship := range possible_constructs {
		if ship.Halite + self.Frame.HaliteAtFast(ship.X, ship.Y) + budget >= self.Frame.Constants.DROPOFF_COST {
			ship.Command = "c"
			self.Frame.Log("Ship %d building dropoff (wmap: %d)", ship.Sid, self.WealthMap.Values[ship.X][ship.Y])
			break
		}
	}
}

func (self *Overmind) TargetSwaps(my_ships []*hal.Ship) {

	for cycle := 0; cycle < 4; cycle++ {

		swap_count := 0

		for i, ship_a := range my_ships {

			if ship_a.TargetIsDropoff() {
				continue
			}

			for _, ship_b := range my_ships[i + 1:] {

				if ship_b.TargetIsDropoff() {
					continue
				}

				a_dist_b := ship_a.Dist(ship_b.Target)
				b_dist_a := ship_b.Dist(ship_a.Target)

				alt_score_a := halite_dist_score(ship_b.TargetHalite(), a_dist_b)
				alt_score_b := halite_dist_score(ship_a.TargetHalite(), b_dist_a)

				if alt_score_a + alt_score_b > ship_a.Score + ship_b.Score {

					ship_a.Target, ship_b.Target = ship_b.Target, ship_a.Target

					ship_a.Score = alt_score_a
					ship_b.Score = alt_score_b

					// self.Frame.Log("Swapped targets for pilots %d, %d (cycle %d)", ship_a.Sid, ship_b.Sid, cycle)
					swap_count++
				}
			}
		}

		if swap_count == 0 {
			return
		}
	}
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

func (self *Overmind) SameTargetCheck(my_ships []*hal.Ship) {

	targets := make(map[hal.Point]int)

	for _, ship := range my_ships {
		targetter_sid, ok := targets[ship.Target]
		if ok && ship.TargetIsDropoff() == false {
			self.Frame.Log("Ships %d and %d looking at same target: %d %d", ship.Sid, targetter_sid, ship.Target.X, ship.Target.Y)
		} else {
			targets[ship.Target] = ship.Sid
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
