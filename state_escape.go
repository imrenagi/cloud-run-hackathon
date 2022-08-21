package main

type Escape struct {
	Player Player
}

func (e Escape) Play() Move {
	front := len(e.Player.FindShooterOnDirection(e.Player.GetDirection()))
	// back := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Backward()))
	left := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Left()))
	right := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Right()))


	// TODO Logic Kabur Perlu di update. Tembak 3 kali user di depan. Kalau masih belum, tembak adjacent lain.
	emptyAdjacents := e.Player.Game.Arena.GetAdjacent(e.Player.GetPosition(), WithEmptyAdjacent())
	if len(emptyAdjacents) == 0 {
		if front > 0 {
			return Throw
		} else if left > 0 {
			return TurnLeft
		} else if right > 0 {
			return TurnRight
		} else {
			return e.Player.Walk()
		}
	}

	// TODO cari adjacent dengan movement paling minimal
	scores := make([][]Move, len(emptyAdjacents))
	for idx, adj := range emptyAdjacents {
		decisions, err := e.Player.MoveNeededToReachAdjacent(adj)
		if err != nil {
			continue
		}
		scores[idx] = decisions
	}

	mostEfficientDecision := scores[0]
	for _, sc := range scores {
		if len(sc) < len(mostEfficientDecision) {
			mostEfficientDecision = sc
		}
	}

	return mostEfficientDecision[0]
}
