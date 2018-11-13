package config

import (
	"flag"
)

var Crash bool
var NoAntiEnemyCollision bool
var NoAntiSelfCollision bool
var RemakeTest bool
var SimTest bool

func ParseCommandLine() {
	flag.BoolVar(&Crash, "crash", false, "randomly crash")
	flag.BoolVar(&NoAntiEnemyCollision, "noantienemycollision", false, "disable anti-enemy-collision")
	flag.BoolVar(&NoAntiSelfCollision, "noantiselfcollision", false, "disable recursive anti-self-collision")
	flag.BoolVar(&RemakeTest, "remaketest", false, "test the frame remaker")
	flag.BoolVar(&SimTest, "simtest", false, "test the simulator")
	flag.Parse()
}
