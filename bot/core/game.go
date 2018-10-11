package core

// ------------------------------------------------------------

type Game struct {

	Constants

	turn						int
	players						int
	pid							int
	width						int
	height						int

	budgets						[]int
	boxes						[][]*Box
	ships						[]*Ship			// Each ship contains a command field for the AI to set
	dropoffs					[]*Dropoff		// The first <player_count> items are always the factories

	ship_xy_lookup				map[Point]*Ship
	ship_id_lookup				map[int]*Ship

	logfile						*Logfile
	flogfile					*Flogfile
	token_parser				*TokenParser

	generate					bool			// Whether the AI wants to send a "g" command
}

func NewGame() *Game {

	game := new(Game)
	game.turn = -1
	game.token_parser = NewTokenParser()

	game.ship_xy_lookup = make(map[Point]*Ship)
	game.ship_id_lookup = make(map[int]*Ship)

	return game
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

func (self *Game) Players() int {
	return self.players
}
