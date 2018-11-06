package core

import (
	"encoding/json"
)

type Constants struct {			// This could conceivably change in the engine in future

	CAPTURE_ENABLED				bool
	CAPTURE_RADIUS				int
	DEFAULT_MAP_HEIGHT			int
	DEFAULT_MAP_WIDTH			int
	DROPOFF_COST				int
	DROPOFF_PENALTY_RATIO		int
	EXTRACT_RATIO				int
	FACTOR_EXP_1				float64
	FACTOR_EXP_2				float64
	INITIAL_ENERGY				int
	INSPIRATION_ENABLED			bool
	INSPIRATION_RADIUS			int
	INSPIRATION_SHIP_COUNT		int
	INSPIRED_BONUS_MULTIPLIER	float64
	INSPIRED_EXTRACT_RATIO		int
	INSPIRED_MOVE_COST_RATIO	int
	MAX_CELL_PRODUCTION			int
	MAX_ENERGY					int
	MAX_PLAYERS					int
	MAX_TURNS					int
	MAX_TURN_THRESHOLD			int
	MIN_CELL_PRODUCTION			int
	MIN_TURNS					int
	MIN_TURN_THRESHOLD			int
	MOVE_COST_RATIO				int
	NEW_ENTITY_ENERGY_COST		int
	PERSISTENCE					float64
	SHIPS_ABOVE_FOR_CAPTURE		int
	STRICT_ERRORS				bool

	GameSeed					int64		`json:"game_seed"`
}

func (self *Frame) LogConstants() {
	s, err := json.MarshalIndent(self.Constants, "", "    ")
	if err != nil {
		self.LogWithoutTurn("%v", err)
	} else {
		self.LogWithoutTurn(string(s))
	}
}
