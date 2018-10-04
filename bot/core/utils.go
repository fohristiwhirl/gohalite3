package core

func mod(x, n int) int {

	// Works for negative x
	// https://dev.to/maurobringolf/a-neat-trick-to-compute-modulo-of-negative-numbers-111e

	return (x % n + n) % n
}
