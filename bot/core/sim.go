package core

func (self *Game) SetPid(pid int) {

	// For simulation purposes, it's simplest just to have
	// each AI set its PID at the start of its turn...

	self.pid = pid
}

func (self *Game) SimGen() *Game {

	duplicate := *self
	g := &duplicate

	g.logfile = nil
	g.flogfile = nil
	g.token_parser = nil

	g.turn += 1
	g.hash = ""

	g.budgets = make([]int, g.players)
	g.halite = Make2dIntArray(g.width, g.height)
	g.ships	= nil
	g.dropoffs = nil
	g.ship_xy_lookup = make(map[Point]*Ship)
	g.ship_id_lookup = make(map[int]*Ship)
	g.box_deltas = make(map[Point]int)			// FIXME - do
	g.generate = make(map[int]bool)

	// Remake some things. No objects are reused.
	// Every pointer is to a new object...

	for pid, val := range self.budgets {
		g.budgets[pid] = val
	}

	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			g.halite[x][y] = self.halite[x][y]
		}
	}

	for _, ship := range self.ships {
		remade := *ship
		remade.Game = g
		g.ships = append(g.ships, &remade)
	}

	for _, dropoff := range self.dropoffs {
		remade := *dropoff
		remade.Game = g
		g.dropoffs = append(g.dropoffs, &remade)
	}

	// Adjust budgets...

	for pid := 0; pid < g.players; pid++ {
		if self.generate[pid] {
			g.budgets[pid] -= g.Constants.NEW_ENTITY_ENERGY_COST
		}
	}

	for _, ship := range g.ships {

		if ship.Command == "c" {

			pid := ship.Owner

			g.budgets[pid] -= g.Constants.DROPOFF_COST
			g.budgets[pid] += ship.Halite
			g.budgets[pid] += g.halite[ship.X][ship.Y]
		}
	}

	// Skip the check for going over budget, we'll trust our AI

	// Make dropoffs...

	for n, ship := range g.ships {

		if ship == nil {
			continue
		}

		if ship.Command != "c" {
			continue
		}

		dropoff := &Dropoff{
			Game:		g,
			Factory:	false,
			Owner: 		ship.Sid,
			X:			ship.X,
			Y:			ship.Y,
		}

		g.dropoffs = append(g.dropoffs, dropoff)
		g.halite[ship.X][ship.Y] = 0

		g.ships[n] = nil
	}

	// Move ships...

	ship_positions := make(map[Point][]*Ship)

	for _, ship := range g.ships {

		if ship == nil {
			continue
		}

		mcr := g.Constants.MOVE_COST_RATIO
		if ship.Inspired { mcr = g.Constants.INSPIRED_MOVE_COST_RATIO }

		if ship.Halite > g.halite[ship.X][ship.Y] / mcr {

			if ship.Command != "" && ship.Command != "o" && ship.Command != "c" {

				ship.Halite -= g.halite[ship.X][ship.Y] / mcr

				dx, dy := StringToDxDy(ship.Command)

				ship.X += dx
				ship.Y += dy
				ship.X = Mod(ship.X, g.width)
				ship.Y = Mod(ship.Y, g.height)
			}
		}

		ship_positions[Point{ship.X, ship.Y}] = append(ship_positions[Point{ship.X, ship.Y}], ship)
	}

	// Find places that want to spawn, so we can check for collisions...

	attempted_spawn_points := make(map[Point]bool)

	for pid := 0; pid < g.players; pid++ {

		if self.generate[pid] {

			factory := g.dropoffs[pid]
			x := factory.X
			y := factory.Y

			attempted_spawn_points[Point{x, y}] = true
		}
	}

	// Delete ships that collide...

	collision_points := make(map[Point]bool)

	all_wrecked_sids := make(map[int]bool)			// Used to delete dead ships

	for point, ships_here := range ship_positions {

		x, y := point.X, point.Y

		if len(ships_here) == 1 && attempted_spawn_points[Point{x, y}] == false {
			continue
		}

		// Collision...

		collision_points[Point{x, y}] = true

		for _, ship := range ships_here {
			g.halite[x][y] += ship.Halite			// Dump the halite on the ground
			all_wrecked_sids[ship.Sid] = true
		}
	}

	for n, ship := range g.ships {
		if ship != nil && all_wrecked_sids[ship.Sid] {
			g.ships[n] = nil
		}
	}

	// Deliveries...

	for _, dropoff := range g.dropoffs {

		pid := dropoff.Owner
		x := dropoff.X
		y := dropoff.Y
		ships_here := ship_positions[Point{x, y}]

		halite_on_ground := g.halite[x][y]

		if halite_on_ground > 0 {
			g.budgets[pid] += halite_on_ground
			g.halite[x][y] = 0
		}

		if len(ships_here) == 1 && collision_points[Point{x, y}] == false {
			if ships_here[0].Owner == pid {
				g.budgets[pid] += ships_here[0].Halite
				ships_here[0].Halite = 0
			}
		}
	}

	// Gen...

	for pid := 0; pid < g.players; pid++ {

		if self.generate[pid] {

			factory := g.dropoffs[pid]
			x := factory.X
			y := factory.Y

			if len(ship_positions[Point{x, y}]) == 1 {
				continue	// i.e. cancel spawn
			}

			sid := g.highest_sid_seen + 1
			g.highest_sid_seen = sid

			ship := &Ship{
				Game:			g,
				Owner:			pid,
				Sid:			sid,
				X:				x,
				Y:				y,
				Halite:			0,
				Inspired:		false,
				Command:		"",
			}

			g.ships = append(g.ships, ship)
		}
	}

	// Mining...

	ibm := int(g.Constants.INSPIRED_BONUS_MULTIPLIER)

	for _, ship := range g.ships {

		if ship == nil {
			continue
		}

		old_ship := self.ship_id_lookup[ship.Sid]

		if old_ship == nil {
			continue
		}

		if old_ship.X == ship.X && old_ship.Y == ship.Y {

			exrat := g.Constants.EXTRACT_RATIO
			if ship.Inspired { exrat = g.Constants.INSPIRED_EXTRACT_RATIO }

			amount_to_mine := (g.halite[ship.X][ship.Y] + exrat - 1) / exrat

			if amount_to_mine + ship.Halite >= g.Constants.MAX_ENERGY {
				amount_to_mine = g.Constants.MAX_ENERGY - ship.Halite
			}

			ship.Halite += amount_to_mine
			g.halite[ship.X][ship.Y] -= amount_to_mine

			if ship.Inspired {

				inspired_bonus := amount_to_mine * ibm

				if inspired_bonus + ship.Halite >= g.Constants.MAX_ENERGY {
					inspired_bonus = g.Constants.MAX_ENERGY - ship.Halite
				}

				ship.Halite += inspired_bonus
			}
		}
	}

	// Clear nils from the ship list...

	var live_ships []*Ship

	for _, ship := range g.ships {
		if ship != nil {
			live_ships = append(live_ships, ship)
		}
	}

	g.ships = live_ships

	// Final cleanups...

	for _, ship := range g.ships {
		ship.Command = ""
		g.ship_xy_lookup[Point{ship.X, ship.Y}] = ship
		g.ship_id_lookup[ship.Sid] = ship
	}

	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			if g.halite[x][y] != self.halite[x][y] {
				g.box_deltas[Point{x, y}] = g.halite[x][y] - self.halite[x][y]
			}
		}
	}

	g.fix_inspiration()

	return g
}
