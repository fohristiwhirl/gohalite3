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

func (self *Game) PrePreParse() {

	// Very early parsing that has to be done before log is opened
	// so that we can open the right log name.

	constants_json := self.token_parser.Str()
	err := json.Unmarshal([]byte(constants_json), &self.Constants)

	if err != nil {
		panic("Couldn't load initial JSON line.")
	}

	self.players = self.token_parser.Int()
	self.pid = self.token_parser.Int()
}

func (self *Game) PreParse() {

	for n := 0; n < self.players; n++ {
		self.dropoffs = append(self.dropoffs, nil)
	}

	for n := 0; n < self.players; n++ {

		pid := self.token_parser.Int()
		x := self.token_parser.Int()
		y := self.token_parser.Int()

		self.dropoffs[pid] = &Dropoff{
			Game:		self,
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

	self.boxes = make([][]*Box, self.width)
	for x := 0; x < self.width; x++ {
		self.boxes[x] = make([]*Box, self.height)
	}

	for y := 0; y < self.height; y++ {
		for x := 0; x < self.width; x++ {
			self.boxes[x][y] = &Box{
				Game: self,
				X: x,
				Y: y,
				Halite: self.token_parser.Int(),
			}
		}
	}
}

func (self *Game) Parse() {

	self.generate = false

	// Set all ships as dead (for stale ref detection by the AI).
	// Also clear all commands...

	for _, ship := range self.ships {
		ship.Alive = false
		ship.Command = ""
	}

	// Hold onto the sid lookup map so we can find
	// the entities while still creating a new map...

	old_ship_id_lookup := self.ship_id_lookup

	// Clear some slices and maps (not dropoffs; old dropoffs are always OK)...

	self.budgets = make([]int, self.players)
	self.ships = nil
	self.ship_xy_lookup = make(map[Point]*Ship)
	self.ship_id_lookup = make(map[int]*Ship)

	// ------------------------------------------------

	self.turn = self.token_parser.Int() - 1			// Out by 1 correction

	self.ParseTime = time.Now()						// Must come after the first read

	for n := 0; n < self.players; n++ {

		pid := self.token_parser.Int()
		ships := self.token_parser.Int()
		dropoffs := self.token_parser.Int()

		self.budgets[pid] = self.token_parser.Int()

		for i := 0; i < ships; i++ {

			// Either update the ship or create it if needed.
			// In any case, it ends up placed in the new maps.

			sid := self.token_parser.Int()

			ship, ok := old_ship_id_lookup[sid]

			if ok == false {
				ship = new(Ship)
				ship.Game = self
			}

			ship.Alive = true
			ship.Inspired = false			// Will detect later

			ship.Owner = pid
			ship.Sid = sid
			ship.X = self.token_parser.Int()
			ship.Y = self.token_parser.Int()
			ship.Halite = self.token_parser.Int()

			self.ships = append(self.ships, ship)
			self.ship_xy_lookup[Point{ship.X, ship.Y}] = ship
			self.ship_id_lookup[ship.Sid] = ship
		}

		for i := 0; i < dropoffs; i++ {

			_ = self.token_parser.Int()					// sid (not needed)
			x := self.token_parser.Int()
			y := self.token_parser.Int()

			known := false

			for _, dropoff := range self.dropoffs {
				if dropoff.X == x && dropoff.Y == y {
					known = true
					break
				}
			}

			if known == false {

				new_dropoff := &Dropoff{
					Game:		self,
					Factory:	false,
					Owner:		pid,
					X:			x,
					Y:			y,
				}

				self.dropoffs = append(self.dropoffs, new_dropoff)
			}
		}
	}

	cell_update_count := self.token_parser.Int()

	for n := 0; n < cell_update_count; n++ {

		x := self.token_parser.Int()
		y := self.token_parser.Int()
		val := self.token_parser.Int()

		self.boxes[x][y].Halite = val
	}

	// ------------------------------------------------
	// Some cleanup...

	sort.Slice(self.ships, func(a, b int) bool {
		return self.ships[a].Sid < self.ships[b].Sid
	})

	self.fix_inspiration()
	self.set_hash()

	return
}

// ---------------------------------------

func (self *Game) SetGenerate(x bool) {
	self.generate = x
}

func (self *Game) Send() {

	var commands []string

	budget_left := self.MyBudget()

	if self.generate {
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
				required -= ship.Box().Halite

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
