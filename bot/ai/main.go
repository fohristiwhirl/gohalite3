package ai

import (
	hal "../core"
)

type Config struct {
}

type Overmind struct {
	Config					*Config
	Game					*hal.Game
}

func NewOvermind(game *hal.Game, config *Config) *Overmind {
	ret := new(Overmind)
	ret.Game = game
	ret.Config = config

	return ret
}

func (self *Overmind) Step() {
	return
}
