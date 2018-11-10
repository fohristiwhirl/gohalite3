package config

import (
	"flag"
)

var Crash bool
var NoAC bool
var RemakeTest bool
var SimTest bool

func ParseCommandLine() {
	flag.BoolVar(&Crash, "crash", false, "randomly crash")
	flag.BoolVar(&NoAC, "noac", false, "disable recursive anti-collision")
	flag.BoolVar(&RemakeTest, "remaketest", false, "test the frame remaker")
	flag.BoolVar(&SimTest, "simtest", false, "test the simulator")
	flag.Parse()
}
