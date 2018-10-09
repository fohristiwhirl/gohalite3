package ai

func mod(x, n int) int {

	// Works for negative x
	// https://dev.to/maurobringolf/a-neat-trick-to-compute-modulo-of-negative-numbers-111e

	return (x % n + n) % n
}

func abs(a int) int {
	if a < 0 { return -a }
	return a
}
