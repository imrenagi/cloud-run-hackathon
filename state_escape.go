package main

type Escape struct {
	// TODO should pass pointer
	Player *Player
}

const maxHitWhenTrapped int = 3

func (e *Escape) Play() Move {
	// front := len(e.Player.FindShooterOnDirection(e.Player.GetDirection()))
	// // back := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Backward()))
	// left := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Left()))
	// right := len(e.Player.FindShooterOnDirection(e.Player.GetDirection().Right()))

	front := e.Player.FindShooterOnDirection(e.Player.GetDirection())
	// back := e.Player.FindShooterOnDirection(e.Player.GetDirection().Backward())
	left := e.Player.FindShooterOnDirection(e.Player.GetDirection().Left())
	right := e.Player.FindShooterOnDirection(e.Player.GetDirection().Right())

	emptyAdjacents := e.Player.Game.Arena.GetAdjacent(e.Player.GetPosition(), WithEmptyAdjacent())
	if len(emptyAdjacents) == 0 {
		if len(front) > 0 {
			e.Player.trappedCount++
			if e.Player.trappedCount > maxHitWhenTrapped {
				e.Player.trappedCount = 0
				if len(left) > 0 {
					return TurnLeft
				}
				if len(right) > 0 {
					return TurnRight
				}
			}
			return Throw
		} else if len(left) > 0 {
			return TurnLeft
		} else if len(right) > 0 {
			return TurnRight
		} else {
			return e.Player.Walk()
		}
	}

	if len(front) > 0 {
		for _, fp := range front {
			for idx, adj := range emptyAdjacents {
				if fp.CanAttack(adj) {
					emptyAdjacents = append(emptyAdjacents[:idx], emptyAdjacents[idx+1:]...)
					break
				}
			}
		}
	}

	e.Player.trappedCount = 0

	// TODO cari adjacent dengan movement paling minimal
	// isOnAttackRange := make([]bool, len(emptyAdjacents))
	requiredMoves := make([][]Move, len(emptyAdjacents))
	for idx, adj := range emptyAdjacents {
		decisions, err := e.Player.MoveNeededToReachAdjacent(adj)
		if err != nil {
			continue
		}
		requiredMoves[idx] = decisions
	}

	// TODO bug: kalau hadap2an, dia malah ngejar. harusnya kabur. ini karena adjacent di depan
	// available dan satisfy the minimum step juga lebih kecil.

	// TODO check whether adjacent is on opponents attack range?
	// kalau gak ada ambil yg ada aja.

	mostEfficientMoves := requiredMoves[0]
	for _, sc := range requiredMoves {
		if len(sc) < len(mostEfficientMoves) {
			mostEfficientMoves = sc
		}
	}

	return mostEfficientMoves[0]
}
