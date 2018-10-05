package core

type Point struct {
	X							int
	Y							int
}

// ------------------------------------------------------------

type Ship struct {
	game						*Game

	X							int
	Y							int
	Owner						int			// Player ID
	Id							int			// Ship ID
	Halite						int
}

func (self *Ship) HaliteAt() int {
	return self.game.HaliteAt(self.X, self.Y)
}

func (self *Ship) Neighbours() []Point {
	return self.game.Neighbours(self.X, self.Y)
}

// ------------------------------------------------------------

type Game struct {
	turn						int
	pid							int
	width						int
	height						int

	halite						[]int
	ships						[]*Ship

	factories					map[int]Point
	dropoffs					map[int][]Point

	ship_xy_lookup				map[Point]*Ship
	ship_id_lookup				map[int]*Ship

	logfile						*Logfile
	token_parser				*TokenParser
}

func NewGame() *Game {

	game := new(Game)
	game.turn = -1
	game.token_parser = NewTokenParser()

	game.ship_xy_lookup = make(map[Point]*Ship)
	game.ship_id_lookup = make(map[int]*Ship)

	// FIXME: Do a parse of whatever info the game sends pre-game.

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

func (self *Game) Neighbours(x, y int) []Point {
	return []Point{
		Point{mod(x - 1, self.width), y},
		Point{x,                      mod(y - 1, self.height)},
		Point{mod(x + 1, self.width), y},
		Point{x,                      mod(y + 1, self.height)},
	}
}

func (self *Game) Factory(pid int) (int, int) {
	factory := self.factories[pid]
	return factory.X, factory.Y
}

func (self *Game) Dropoffs(pid int) []Point {

	var ret []Point

	dropoffs := self.dropoffs[pid]

	for _, point := range dropoffs {
		ret = append(ret, Point{point.X, point.Y})
	}

	return ret
}
