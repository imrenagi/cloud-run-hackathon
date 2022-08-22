package main

type Escape struct {
	// TODO should pass pointer
	Player *Player
}

const maxHitWhenTrapped int = 3

func (e *Escape) Play() Move {
	front := len(e.Player.FindShooterOnDirection(e.Player.GetDirection()))
	// back := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Backward()))
	left := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Left()))
	right := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Right()))

	emptyAdjacents := e.Player.Game.Arena.GetAdjacent(e.Player.GetPosition(), WithEmptyAdjacent())
	if len(emptyAdjacents) == 0 {
		if front > 0 {
			e.Player.trappedCount++
			if e.Player.trappedCount > maxHitWhenTrapped {
				e.Player.trappedCount = 0
				if left > 0 { return TurnLeft }
				if right > 0 { return TurnRight }
			}
			return Throw
		} else if left > 0 {
			return TurnLeft
		} else if right > 0 {
			return TurnRight
		} else {
			return e.Player.Walk()
		}
	}

	e.Player.trappedCount = 0

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
