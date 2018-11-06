package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------

type TokenParser struct {
	scanner			*bufio.Scanner
	count			int
}

func NewTokenParser() *TokenParser {
	ret := new(TokenParser)
	ret.scanner = bufio.NewScanner(os.Stdin)
	ret.scanner.Split(bufio.ScanWords)
	return ret
}

func (self *TokenParser) Int() int {
	bl := self.scanner.Scan()
	if bl == false {
		err := self.scanner.Err()
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		} else {
			panic(fmt.Sprintf("End of input."))
		}
	}
	ret, err := strconv.Atoi(self.scanner.Text())
	if err != nil {
		panic(fmt.Sprintf("TokenReader.Int(): Atoi failed at token %d: \"%s\"", self.count, self.scanner.Text()))
	}

	self.count++
	return ret
}

func (self *TokenParser) Str() string {
	bl := self.scanner.Scan()
	if bl == false {
		err := self.scanner.Err()
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		} else {
			panic(fmt.Sprintf("End of input."))
		}
	}
	return self.scanner.Text()
}

// ---------------------------------------

func (self *Frame) PrePreParse() {

	// Very early parsing that has to be done before log is opened
	// so that we can open the right log name.

	constants_json := self.token_parser.Str()
	err := json.Unmarshal([]byte(constants_json), &self.Constants)

	if err != nil {
		panic("Couldn't load initial JSON line.")
	}

	self.players = self.token_parser.Int()
	self.pid = self.token_parser.Int()
	self.__true_pid = self.pid
}

func (self *Frame) PreParse() {

	for n := 0; n < self.players; n++ {
		self.dropoffs = append(self.dropoffs, nil)
	}

	for n := 0; n < self.players; n++ {

		pid := self.token_parser.Int()
		x := self.token_parser.Int()
		y := self.token_parser.Int()

		self.dropoffs[pid] = &Dropoff{
			Frame:		self,
			Factory:	true,
			Owner:		pid,
			X:			x,
			Y:			y,
		}
	}

	// The factories are stored in the dropoffs
	// list at the very start, in player order.

	sort.Slice(self.dropoffs, func(a, b int) bool {
		return self.dropoffs[a].Owner < self.dropoffs[b].Owner
	})

	self.width = self.token_parser.Int()
	self.height = self.token_parser.Int()

	self.halite = Make2dIntArray(self.width, self.height)

	for y := 0; y < self.height; y++ {
		for x := 0; x < self.width; x++ {
			val := self.token_parser.Int()
			self.halite[x][y] = val
		}
	}

	// Now put the frame into a valid Turn 0 state, mostly for sim purposes...

	self.turn = 0
	self.set_hash()

	self.budgets = make([]int, self.players)

	for pid := 0; pid < self.players; pid++ {
		self.budgets[pid] = self.Constants.INITIAL_ENERGY
	}

	self.ship_xy_lookup = make(map[Point]*Ship)
	self.ship_id_lookup = make(map[int]*Ship)
	self.box_deltas = make(map[Point]int)
	self.generate = make(map[int]bool)
}

func (self *Frame) Parse() {

	// Note: creates brand new objects for literally everything;
	// No holding onto the old ones.

	// Save some things we will need later...

	old_dropoffs := self.dropoffs
	old_halite := self.halite
	old_ship_id_lookup := self.ship_id_lookup

	// Clear all the things...

	self.budgets = make([]int, self.players)
	self.halite = Make2dIntArray(self.width, self.height)
	self.ships = nil
	self.dropoffs = nil
	self.ship_xy_lookup = make(map[Point]*Ship)
	self.ship_id_lookup = make(map[int]*Ship)
	self.box_deltas = make(map[Point]int)
	self.generate = make(map[int]bool)

	// Remake the factories...

	for _, factory := range old_dropoffs[0:self.players] {
		remade := *factory
		self.dropoffs = append(self.dropoffs, &remade)
	}

	// ------------------------------------------------

	self.turn = self.token_parser.Int() - 1			// Out by 1 correction

	self.ParseTime = time.Now()						// Must come after the first read

	for n := 0; n < self.players; n++ {

		pid := self.token_parser.Int()
		ships := self.token_parser.Int()
		dropoffs := self.token_parser.Int()

		self.budgets[pid] = self.token_parser.Int()

		for i := 0; i < ships; i++ {

			ship := new(Ship)

			sid := self.token_parser.Int()

			old_ship, ok := old_ship_id_lookup[sid]
			if ok {
				*ship = *old_ship					// (Shallow) copy all the AI stuff, if available. Everything else is replaced below.
			}										// If the AI stuff is not available, the zeroed vars must work.

			ship.X = self.token_parser.Int()
			ship.Y = self.token_parser.Int()
			ship.Halite = self.token_parser.Int()

			ship.Frame = self
			ship.Sid = sid
			ship.Inspired = false					// Will set this correctly later
			ship.Owner = pid
			ship.Command = ""

			self.ships = append(self.ships, ship)
			self.ship_xy_lookup[Point{ship.X, ship.Y}] = ship
			self.ship_id_lookup[ship.Sid] = ship

			if ship.Sid > self.highest_sid_seen {
				self.highest_sid_seen = ship.Sid
			}
		}

		for i := 0; i < dropoffs; i++ {

			_ = self.token_parser.Int()				// sid (not needed)

			dropoff := new(Dropoff)
			dropoff.Frame = self

			dropoff.X = self.token_parser.Int()
			dropoff.Y = self.token_parser.Int()

			dropoff.Factory = false
			dropoff.Owner = pid

			self.dropoffs = append(self.dropoffs, dropoff)
		}
	}

	for x := 0; x < self.width; x++ {
		for y := 0; y < self.height; y++ {
			self.halite[x][y] = old_halite[x][y]
		}
	}

	cell_update_count := self.token_parser.Int()

	for n := 0; n < cell_update_count; n++ {

		x := self.token_parser.Int()
		y := self.token_parser.Int()

		val := self.token_parser.Int()
		old_val := old_halite[x][y]

		if val != old_val {
			self.halite[x][y] = val
			self.box_deltas[Point{x, y}] = val - old_val
		}
	}

	// ------------------------------------------------
	// Some cleanup...

	sort.Slice(self.ships, func(a, b int) bool {
		return self.ships[a].Sid < self.ships[b].Sid
	})

	self.fix_inspiration()
	self.set_hash()

	// self.Log("Parsing took %v", time.Now().Sub(self.ParseTime))

	return
}

// ---------------------------------------

func (self *Frame) SetGenerate(val bool) {
	self.generate[self.pid] = val
}

func (self *Frame) Send() {

	self.pid = self.__true_pid		// In case any simulating has been going on.

	var commands []string

	budget_left := self.MyBudget()

	if self.generate[self.pid] {
		if budget_left >= self.Constants.NEW_ENTITY_ENERGY_COST {
			commands = append(commands, "g")
			budget_left -= self.Constants.NEW_ENTITY_ENERGY_COST
		} else {
			self.Log("Warning: GENERATE command blocked due to lack of resources!")
		}
	}

	for _, ship := range self.ships {
		if ship.Owner == self.pid && ship.Command != "" {
			if ship.Command == "c" {

				required := self.Constants.DROPOFF_COST
				required -= ship.Halite
				required -= self.HaliteAt(ship)

				if budget_left >= required {
					commands = append(commands, fmt.Sprintf("c %d", ship.Sid))
					budget_left -= required
				} else {
					self.Log("Warning: CONSTRUCT command blocked due to lack of resources!")
				}
			} else {
				commands = append(commands, fmt.Sprintf("m %d %s", ship.Sid, ship.Command))
			}
		}
	}

	output := strings.Join(commands, " ")
	fmt.Printf(output)
	fmt.Printf("\n")
	return
}
