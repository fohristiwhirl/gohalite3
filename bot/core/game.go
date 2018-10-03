package core

type Game struct {
	inited						bool
	turn						int
	pid							int
	width						int
	height						int

	logfile						*Logfile
	token_parser				*TokenParser
}

func NewGame() *Game {
	game := new(Game)
	game.turn = -1
	game.token_parser = NewTokenParser()

	// FIXME: Do a parse of whatever info the game sends pre-game.

	game.inited = true		// Just means Parse() will increment the turn value before parsing.
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
