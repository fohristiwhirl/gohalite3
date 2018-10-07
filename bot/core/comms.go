package core

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
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

	self.constants_json = self.token_parser.Str()
	self.players = self.token_parser.Int()
	self.pid = self.token_parser.Int()
}

func (self *Game) PreParse() {

	self.factories = make([]Point, self.players)

	for n := 0; n < self.players; n++ {

		pid := self.token_parser.Int()
		x := self.token_parser.Int()
		y := self.token_parser.Int()

		self.factories[pid] = Point{x, y}
	}

	self.width = self.token_parser.Int()
	self.height = self.token_parser.Int()

	self.halite = make([]int, self.width * self.height)

	for y := 0; y < self.height; y++ {
		for x := 0; x < self.width; x++ {
			self.halite[y * self.width + x] = self.token_parser.Int()
		}
	}
}

func (self *Game) Parse() {

	self.generate = false

	// Hold onto the sid lookup map so we can find
	// the entities while still creating a new map...

	old_ship_id_lookup := self.ship_id_lookup

	self.budgets = make([]int, self.players)

	self.ships = nil
	self.dropoffs = make([][]Point, self.players)

	self.ship_xy_lookup = make(map[Point]*Ship)
	self.ship_id_lookup = make(map[int]*Ship)

	// ------------------------------------------------

	self.turn = self.token_parser.Int()

	for n := 0; n < self.players; n++ {

		pid := self.token_parser.Int()
		ships := self.token_parser.Int()
		dropoffs := self.token_parser.Int()

		self.budgets[pid] = self.token_parser.Int()

		for i := 0; i < ships; i++ {

			// Either update the entity or create it if needed.
			// In any case, it ends up placed in the new maps.

			sid := self.token_parser.Int()

			ship, ok := old_ship_id_lookup[sid]

			if ok == false {
				ship = new(Ship)
				ship.game = self
				ship.Owner = pid
			}

			ship.Command = ""

			ship.Id = sid
			ship.X = self.token_parser.Int()
			ship.Y = self.token_parser.Int()
			ship.Halite = self.token_parser.Int()

			self.ship_xy_lookup[Point{ship.X, ship.Y}] = ship
			self.ship_id_lookup[ship.Id] = ship

			self.ships = append(self.ships, ship)
		}

		for i := 0; i < dropoffs; i++ {

			_ = self.token_parser.Int()		// sid (not needed)
			x := self.token_parser.Int()
			y := self.token_parser.Int()

			self.dropoffs[pid] = append(self.dropoffs[pid], Point{x, y})
		}
	}

	cell_update_count := self.token_parser.Int()

	for n := 0; n < cell_update_count; n++ {

		x := self.token_parser.Int()
		y := self.token_parser.Int()
		val := self.token_parser.Int()

		self.halite[y * self.width + x] = val
	}

	// ------------------------------------------------
	// Some cleanup...

	sort.Slice(self.ships, func(a, b int) bool {
		return self.ships[a].Id < self.ships[b].Id
	})

	return
}

func (self *Game) SetGenerate(x bool) {
	self.generate = x
}

func (self *Game) Send() {

	var commands []string

	if self.generate {
		commands = append(commands, "g")
	}

	for _, ship := range self.ships {
		if ship.Owner == self.pid && ship.Command != "" {
			commands = append(commands, ship.Command)
		}
	}

	output := strings.Join(commands, " ")
	fmt.Printf(output)
	fmt.Printf("\n")
	return
}

/*
	Example Parse() input for 2 players

	4				- turn

	0 2 1 3000		- pid, ships, dropoffs, budget
	1 28 28 0		- sid, x, y, energy
	0 27 28 22      - sid, x, y, energy
	2 10 10			- dropoff id, x, y

	1 1 0 3000		- pid, ships, dropoffs, budget
	3 15 17 0		- sid, x, y, energy

	2               - cell update count
	27 28 63        - x y val
	28 28 0         - x y val
*/
