package core

func mod(x, n int) int {

	// Works for negative x
	// https://dev.to/maurobringolf/a-neat-trick-to-compute-modulo-of-negative-numbers-111e

	return (x % n + n) % n
}

func (self *Game) DxDy(x1, y1, x2, y2 int) (int, int) {

	// How to get from (x1, y1) to (x2, y2)

	x1 = mod(x1, self.width)
	y1 = mod(y1, self.height)

	x2 = mod(x2, self.width)
	y2 = mod(y2, self.height)

	// Naive result:

	foo := Vector{x2 - x1, y2 - y1}

	// Change for x wrap...

	if x1 < x2 {					// Naive is positive (right)
		x3 := x1 + self.width
		if x3 - x2 < x2 - x1 {
			foo.X = x2 - x3			// But correct is negative (left)
		}
	} else if x2 < x1 {				// Naive is negative (left)
		x0 := x1 - self.width
		if x2 - x0 < x1 - x2 {
			foo.X = x2 - x0			// But correct is positive (right)
		}
	}

	// Likewise for y wrap...

	if y1 < y2 {
		y3 := y1 + self.height
		if y3 - y2 < y2 - y1 {
			foo.Y = y2 - y3
		}
	} else if y2 < y1 {
		y0 := y1 - self.height
		if y2 - y0 < y1 - y2 {
			foo.Y = y2 - y0
		}
	}

	return foo.X, foo.Y
}
