package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func (self *TokenParser) Float() float64 {
	bl := self.scanner.Scan()
	if bl == false {
		err := self.scanner.Err()
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		} else {
			panic(fmt.Sprintf("End of input."))
		}
	}
	ret, err := strconv.ParseFloat(self.scanner.Text(), 64)
	if err != nil {
		panic(fmt.Sprintf("TokenReader.Float(): ParseFloat failed at token %d: \"%s\"", self.count, self.scanner.Text()))
	}

	self.count++
	return ret
}

func (self *TokenParser) Bool() bool {
	bl := self.scanner.Scan()
	if bl == false {
		err := self.scanner.Err()
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		} else {
			panic(fmt.Sprintf("End of input."))
		}
	}
	val, err := strconv.Atoi(self.scanner.Text())
	if err != nil {
		panic(fmt.Sprintf("TokenReader.Bool(): Atoi failed at token %d: \"%s\"", self.count, self.scanner.Text()))
	}
	if val != 0 && val != 1 {
		panic(fmt.Sprintf("TokenReader.Bool(): Non-bool at token %d: \"%s\"", self.count, self.scanner.Text()))
	}

	self.count++
	if val == 0 {
		return false
	} else {
		return true
	}
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

	// FIXME: check I have these the right way round...

	for y := 0; y < self.height; y++ {
		for x := 0; x < self.width; x++ {
			self.halite[y * self.width + x] = self.token_parser.Int()
		}
	}
}

func pretend(...interface{}) {
	return
}

func (self *Game) Parse() {

	// Hold onto the sid lookup map so we can find
	// the entities while still creating a new map...

	old_ship_id_lookup := self.ship_id_lookup

	self.budgets = make([]int, self.players)
	self.ship_xy_lookup = make(map[Point]*Ship)
	self.ship_id_lookup = make(map[int]*Ship)

	self.ships = nil

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

			ship.Id = sid
			ship.X = self.token_parser.Int()
			ship.Y = self.token_parser.Int()
			ship.Halite = self.token_parser.Int()

			self.ship_xy_lookup[Point{ship.X, ship.Y}] = ship
			self.ship_id_lookup[ship.Id] = ship

			self.ships = append(self.ships, ship)

		}

		// Dropoffs.
		// The following is not known to be correct...
		// FIXME: update self.dropoffs

		for i := 0; i < dropoffs; i++ {

			sid := self.token_parser.Int()
			x := self.token_parser.Int()
			y := self.token_parser.Int()

			pretend(sid, x, y)
		}
	}

	cell_update_count := self.token_parser.Int()

	for n := 0; n < cell_update_count; n++ {

		x := self.token_parser.Int()
		y := self.token_parser.Int()
		val := self.token_parser.Int()

		self.halite[y * self.width + x] = val
	}

	return
}

func (self *Game) Send() {
	fmt.Printf("\n")
	return
}




	/*
		4				- turn
		0 2 0 3000		- pid, ships, dropoff count (??), budget
		1 28 28 0		- sid, x, y, energy
		0 27 28 22      - sid, x, y, energy
		2               - cell update count
		27 28 63        - x y val
		28 28 0         - x y val

		// Probably dropoff info after
	*/
