package core

import (
	"crypto/sha1"
	"fmt"
)

func Mod(x, n int) int {

	// Works for negative x
	// https://dev.to/maurobringolf/a-neat-trick-to-compute-modulo-of-negative-numbers-111e

	return (x % n + n) % n
}

func Abs(a int) int {
	if a < 0 { return -a }
	return a
}

func StringToDxDy(s string) (int, int) {

	switch s {

	case "e":
		return 1, 0
	case "w":
		return -1, 0
	case "s":
		return 0, 1
	case "n":
		return 0, -1
	case "c":
		return 0, 0
	case "o":
		return 0, 0
	case "":
		return 0, 0
	}

	panic("StringToDxDy() got illegal string")
}

func HashFromString(datastring string) string {
	data := []byte(datastring)
	sum := sha1.Sum(data)
	return fmt.Sprintf("%x", sum)
}

var fluorine_colours = []string{"#c5ec98", "#ff9999", "#ffbe00", "#66cccc"}

func FluorineColour(pid int) string {
	if pid >= 0 && pid < len(fluorine_colours) {
		return fluorine_colours[pid]
	}
	return "#ffffff"
}

func Make2dIntArray(width, height int) [][]int {
	ret := make([][]int, width)
	for x := 0; x < width; x++ {
		ret[x] = make([]int, height)
	}
	return ret
}
