package config

import (
	"flag"
)

var Crash bool
var NoAntiEnemyCollision bool
var NoAntiSelfCollision bool
var RemakeTest bool
var SimTest bool

var gen_min_arg float64
var GenMin float64 = 0.5

func ParseCommandLine() {

	flag.BoolVar(&Crash, "crash", false, "randomly crash")
	flag.BoolVar(&NoAntiEnemyCollision, "noantienemycollision", false, "disable anti-enemy-collision")
	flag.BoolVar(&NoAntiSelfCollision, "noantiselfcollision", false, "disable recursive anti-self-collision")
	flag.BoolVar(&RemakeTest, "remaketest", false, "test the frame remaker")
	flag.BoolVar(&SimTest, "simtest", false, "test the simulator")

	flag.Float64Var(&gen_min_arg, "genmin", -1, "halite required to generate a ship")

	flag.Parse()
}

func SetGenMin(size int, players int) {

	if gen_min_arg >= 0 {
		GenMin = gen_min_arg
		return
	}

	if players == 2 {
		GenMin = 0.5
		return
	}

	if size <= 32 {
		GenMin = 0.35
	} else if size <= 40 {
		GenMin = 0.4
	} else if size <= 48 {
		GenMin = 0.4
	} else if size <= 56 {
		GenMin = 0.47
	} else {
		GenMin = 0.5
	}
}
