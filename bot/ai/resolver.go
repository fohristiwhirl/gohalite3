package ai

// Given ships with desired moves, resolve them without collisions, as far as reasonable...

import (
	"../config"
	hal "../core"
)

type MoveBook struct {
	width				int
	height				int
	book				[][]*hal.Ship
}

func NewMoveBook(width, height int) *MoveBook {
	o := new(MoveBook)
	o.width = width
	o.height = height
	o.book = make([][]*hal.Ship, o.width)
	for x := 0; x < o.width; x++ {
		o.book[x] = make([]*hal.Ship, o.height)
	}
	return o
}

func (self *MoveBook) Booker(pos hal.XYer) *hal.Ship {
	x := hal.Mod(pos.GetX(), self.width)
	y := hal.Mod(pos.GetY(), self.height)
	return self.book[x][y]
}

func (self *MoveBook) SetBook(ship *hal.Ship, pos hal.XYer) {
	x := hal.Mod(pos.GetX(), self.width)
	y := hal.Mod(pos.GetY(), self.height)
	self.book[x][y] = ship
}

func Resolve(frame *hal.Frame, my_ships []*hal.Ship) *MoveBook {

	// Resolve the moves by setting the ships' actual .Commands
	// and return the MoveBook.

	book := NewMoveBook(frame.Width(), frame.Height())

	for _, ship := range my_ships {
		if ship.Desires[0] == "o" {
			ship.Move("o")
			book.SetBook(ship, ship)
		}
	}

	for cycle := 0; cycle < 5; cycle++ {

		for _, ship := range my_ships {

			if ship.Command != "" {
				continue
			}

			// Special case: if ship is next to a dropoff and is in its mad dash, always move.
			// And don't set the book, it can only confuse matters...

			if ship.TargetIsDropoff() && ship.Dist(ship.Target()) == 1 && ship.FinalDash {
				ship.Move(ship.Desires[0])
				continue
			}

			// Normal case...

			for _, desire := range ship.Desires {

				new_loc := ship.LocationAfterMove(desire)
				booker := book.Booker(new_loc)

				if booker == nil {
					ship.Move(desire)
					book.SetBook(ship, new_loc)
					break
				} else {
					if booker.Command == "o" {		// Never clear a booking made by a stationary ship
						continue
					}
					if booker.Halite < ship.Halite {
						ship.Move(desire)
						book.SetBook(ship, new_loc)
						booker.ClearMove()
						break
					}
				}
			}
		}
	}

	if config.NoAC == false {
		for _, ship := range my_ships {
			if ship.Command == "" {
				PreventCollision(ship, book)
			}
		}
	}

	return book
}

func PreventCollision(innocent *hal.Ship, book *MoveBook) {

	// If called correctly, the innocent ship is motionless but has
	// some incoming ship that's going to collide with it.

	if innocent.OnDropoff() && innocent.FinalDash {
		innocent.Frame.Log("PreventCollision(Ship %d) -- ship was on dropoff and in final dash mode", innocent.Sid)
		return
	}

	if innocent.Command != "" && innocent.Command != "o" {
		innocent.Frame.Log("PreventCollision(Ship %d) -- this ship was moving", innocent.Sid)
		return
	}

	villain := book.Booker(innocent)

	if villain == nil {
		innocent.Frame.Log("PreventCollision(Ship %d) -- no incoming ship noted in the book", innocent.Sid)
		return
	}

	if villain == innocent {
		innocent.Frame.Log("PreventCollision(Ship %d) -- this ship was already the booker", innocent.Sid)
		return
	}

	innocent.Move("o")						// It either had this already, or had ""
	book.SetBook(innocent, innocent)

	villain.Move("o")

	// Do we need to recurse?

	if book.Booker(villain) != nil {
		PreventCollision(villain, book)		// Will set the book for the villain (which will be considered the innocent by the recursed function)
	} else {
		book.SetBook(villain, villain)
	}
}
