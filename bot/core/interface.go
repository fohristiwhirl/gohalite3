package core

type FrameXYer interface {
	GetFrame()	*Frame
	GetX()		int
	GetY()		int
}

type XYer interface {
	GetX()		int
	GetY()		int
}

func DxDy(self FrameXYer, other XYer) (int, int) {

	// How to get from (x1, y1) to (x2, y2)

	frame := self.GetFrame()
	width := frame.Width()
	height := frame.Height()

	x1 := Mod(self.GetX(), width)
	y1 := Mod(self.GetY(), height)

	x2 := Mod(other.GetX(), width)
	y2 := Mod(other.GetY(), height)

	// Naive result:

	dx, dy := x2 - x1, y2 - y1

	// Change for x wrap...

	if x1 < x2 {					// Naive is positive (right)
		x3 := x1 + width
		if x3 - x2 < x2 - x1 {
			dx = x2 - x3			// But correct is negative (left)
		}
	} else if x2 < x1 {				// Naive is negative (left)
		x0 := x1 - width
		if x2 - x0 < x1 - x2 {
			dx = x2 - x0			// But correct is positive (right)
		}
	}

	// Likewise for y wrap...

	if y1 < y2 {
		y3 := y1 + height
		if y3 - y2 < y2 - y1 {
			dy = y2 - y3
		}
	} else if y2 < y1 {
		y0 := y1 - height
		if y2 - y0 < y1 - y2 {
			dy = y2 - y0
		}
	}

	return dx, dy
}

func Dist(self FrameXYer, other XYer) int {
	dx, dy := DxDy(self, other)
	return Abs(dx) + Abs(dy)
}

func SamePlace(self FrameXYer, other XYer) bool {

	frame := self.GetFrame()
	width := frame.Width()
	height := frame.Height()

	x1 := Mod(self.GetX(), width)
	y1 := Mod(self.GetY(), height)

	x2 := Mod(other.GetX(), width)
	y2 := Mod(other.GetY(), height)

	if x1 == x2 && y1 == y2 {
		return true
	}

	return false
}
