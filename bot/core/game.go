package core

type Point struct {
	X							int
	Y							int
}

// ------------------------------------------------------------

type Game struct {
	turn						int
	players						int
	pid							int
	width						int
	height						int

	constants_json				string

	budgets						[]int
	halite						[]int
	ships						[]*Ship

	factories					[]Point
	dropoffs					[][]Point

	ship_xy_lookup				map[Point]*Ship
	ship_id_lookup				map[int]*Ship

	logfile						*Logfile
	token_parser				*TokenParser

	generate					bool
}

func NewGame() *Game {

	game := new(Game)
	game.turn = -1
	game.token_parser = NewTokenParser()

	game.ship_xy_lookup = make(map[Point]*Ship)
	game.ship_id_lookup = make(map[int]*Ship)

	return game
}

func (self *Game) HaliteAt(x, y int) int {

	// Translate out-of-bounds coordinates...
	// Use a special function since % doesn't work for negative.

	x = mod(x, self.width)
	y = mod(y, self.height)

	return self.halite[y * self.width + x]
}

func (self *Game) ShipAt(x, y int) (*Ship, bool) {
	ret, ok := self.ship_xy_lookup[Point{x, y}]
	return ret, ok
}

func (self *Game) Pid() int {
	return self.pid
}

func (self *Game) Turn() int {
	return self.turn
}

func (self *Game) Width() int {
	return self.width
}

func (self *Game) Height() int {
	return self.height
}

func (self *Game) NeighbourPoints(x, y int) []Point {
	return []Point{
		Point{mod(x - 1, self.width), y},
		Point{x,                      mod(y - 1, self.height)},
		Point{mod(x + 1, self.width), y},
		Point{x,                      mod(y + 1, self.height)},
	}
}

func (self *Game) FactoryXY(pid int) (int, int) {
	factory := self.factories[pid]
	return factory.X, factory.Y
}

func (self *Game) DropoffPoints(pid int) []Point {

	var ret []Point

	dropoffs := self.dropoffs[pid]

	for _, point := range dropoffs {
		ret = append(ret, Point{point.X, point.Y})
	}

	return ret
}

func (self *Game) ReturnPoints(pid int) []Point {

	// So-called "dropoff points", plus factory.

	ret := self.DropoffPoints(pid)
	factory := self.factories[pid]
	ret = append(ret, Point{factory.X, factory.Y})
	return ret
}

func (self *Game) MyShips() []*Ship {
	return self.PlayerShips(self.pid)
}

func (self *Game) PlayerShips(pid int) []*Ship {

	var ret []*Ship

	for _, ship := range self.ships {
		if ship.Owner == pid {
			ret = append(ret, ship)
		}
	}

	return ret
}

func (self *Game) MyBudget() int {
	return self.PlayerBudget(self.pid)
}

func (self *Game) PlayerBudget(pid int) int {
	return self.budgets[pid]
}
