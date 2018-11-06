package ai

import (
	"fmt"
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
	o.Config = config

	o.WealthMap = maps.NewWealthMap(frame)
	o.InitialGroundHalite = frame.GroundHalite()

	return o
}

func (self *Overmind) Step(frame *hal.Frame) {

	// Various calls rely on this happening..................
	// Always have this first.

	frame.SetPid(self.Pid)

	// Various other initialisation..........................

	rand.Seed(int64(frame.MyBudget() + self.Pid))
	self.WealthMap.Update(frame)
	self.SetTurnParameters(frame)

	// What each ship wants to do............................

	my_ships := frame.MyShips()

	for _, ship := range my_ships {
		NewTurn(ship)
	}

	target_book := hal.Make2dBoolArray(frame.Width(), frame.Height())		// What points are targets. Updated for each ship.

	for _, ship := range my_ships {
		self.SetTarget(ship, target_book)
	}

	TargetSwaps(my_ships)

	for _, ship := range my_ships {
		self.SetDesires(ship)
	}

	// Resolve the desired moves.............................

	move_book := Resolve(frame, my_ships)

	// Other.................................................

	self.MaybeBuild(frame, my_ships, move_book)

	for _, ship := range my_ships {
		FlogTarget(ship)
	}

	for _, report := range SameTargetReports(my_ships) {
		frame.Log(report)
	}

	return
}

func (self *Overmind) MaybeBuild(frame *hal.Frame, my_ships []*hal.Ship, move_book *MoveBook) {

	budget := frame.MyBudget()

	factory := frame.MyFactory()
	willing := true

	if self.InitialGroundHalite / (frame.GroundHalite() + 1) >= 2 {		// remember int division, also div-by-zero
		willing = false
	}

	if frame.Turn() >= frame.Constants.MAX_TURNS / 2 {
		willing = false
	}

	if budget >= frame.Constants.NEW_ENTITY_ENERGY_COST && move_book.Booker(factory) == nil && willing {
		frame.SetGenerate(true)
		budget -= frame.Constants.NEW_ENTITY_ENERGY_COST
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

		if frame.HaliteAtFast(ship.X, ship.Y) == 0 {		// Cheap way to avoid building on enemy dropoff / factory
			continue
		}

		possible_constructs = append(possible_constructs, ship)
	}

	sort.Slice(possible_constructs, func (a, b int) bool {

		return	self.WealthMap.Values[possible_constructs[a].X][possible_constructs[a].Y] <
				self.WealthMap.Values[possible_constructs[b].X][possible_constructs[b].Y]
	})

	for _, ship := range possible_constructs {
		if ship.Halite + frame.HaliteAtFast(ship.X, ship.Y) + budget >= frame.Constants.DROPOFF_COST {
			ship.Command = "c"
			frame.Log("Ship %d building dropoff (wmap: %d)", ship.Sid, self.WealthMap.Values[ship.X][ship.Y])
			break
		}
	}
}

func (self *Overmind) SetTurnParameters(frame *hal.Frame) {

	current_ground_halite := 0

	for x := 0; x < frame.Width(); x++ {
		for y := 0; y < frame.Height(); y++ {
			current_ground_halite += frame.HaliteAtFast(x, y)
		}
	}

	avg_ground_halite := current_ground_halite / (frame.Width() * frame.Height())

	self.HappyThreshold = avg_ground_halite / 2			// Above this, ground is sticky
	self.IgnoreThreshold = avg_ground_halite * 2 / 3	// Less than this not counted for targeting
}

func (self *Overmind) ShouldMine(frame *hal.Frame, halite_carried int, pos, tar hal.XYer) bool {

	// Whether a ship -- if it were carrying n halite, at pos, with specified target -- would stop to mine.

	if halite_carried >= 800 {
		return false
	}

	pos_halite := frame.HaliteAt(pos)
	tar_halite := frame.HaliteAt(tar)

	if pos_halite > self.HappyThreshold {
		if pos_halite > tar_halite / 3 {			// This is a bit odd since the test even happens when target is dropoff.
			return true
		}
	}

	return false
}

func HaliteDistScore(halite, dist int) float32 {
	return float32(halite) / float32((dist + 1) * (dist + 1))	// Avoid div-by-zero
}

func TargetSwaps(my_ships []*hal.Ship) {

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

				alt_score_a := HaliteDistScore(ship_b.TargetHalite(), a_dist_b)
				alt_score_b := HaliteDistScore(ship_a.TargetHalite(), b_dist_a)

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

func SameTargetReports(my_ships []*hal.Ship) []string {

	var ret []string

	targets := make(map[hal.Point]int)

	for _, ship := range my_ships {
		targetter_sid, ok := targets[ship.Target]
		if ok && ship.TargetIsDropoff() == false {
			ret = append(ret, fmt.Sprintf("Ships %d and %d looking at same target: %d %d", ship.Sid, targetter_sid, ship.Target.X, ship.Target.Y))
		} else {
			targets[ship.Target] = ship.Sid
		}
	}

	return ret
}
