package core

type Point struct {
	x							int
	y							int
}

// ------------------------------------------------------------

type Box struct {
	game						*Game
	x							int
	y							int
	halite						int
}

// ------------------------------------------------------------

type Ship struct {
	game						*Game
	x							int
	y							int
	pid							int			// Player ID
	sid							int			// Ship ID
}

func (self *Ship) Box() *Box {
	return self.game.Box(self.x, self.y)
}

// ------------------------------------------------------------

type Game struct {
	turn						int
	pid							int
	width						int
	height						int

	boxes						[]*Box
	ships						[]*Ship

	ship_xy_lookup				map[Point]*Ship
	ship_id_lookup				map[int]*Ship

	logfile						*Logfile
	token_parser				*TokenParser
}

func NewGame() *Game {
	game := new(Game)
	game.turn = -1
	game.token_parser = NewTokenParser()

	self.ship_xy_lookup = make(map[Point]*Ship)
	self.ship_id_lookup = make(map[int]*Ship)

	// FIXME: Do a parse of whatever info the game sends pre-game.

	return game
}

func (self *Game) Box(x, y int) *Box {

	// Translate out-of-bounds coordinates...
	// Use a special function since % doesn't work for negative.

	x = mod(x, self.width)
	y = mod(y, self.height)

	return self.boxes[y * self.width + x]
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
