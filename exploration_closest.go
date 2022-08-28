package main


type ClosestEnemy struct {
}

func (a *ClosestEnemy) Explore(p *Player) Move {
	targets := p.FindClosestPlayers()
	if len(targets) == 0 {
		return Throw
	}

	var path Path
	for _, target := range targets {
		aStar := NewAStar(p.Game.Arena,
			WithIsUnblockFn(CheckTargetSurroundingAttackRangeFn(target)),
		)
		var err error
		path, err = aStar.SearchPath(p.GetPosition(), target.GetPosition())
		if err != nil {
			continue
		}
		if len(path) > 0 {
			break
		}
	}

	if len(path) == 0 {
		return p.Walk()
	}

	moves := p.RequiredMoves(path, WithOnlyNextMove())
	if len(moves) > 0 {
		return moves[0]
	} else {
		return p.Walk()
	}
}


// CheckTargetSurroundingAttackRangeFn checks whether a target and its facing opponent can hit p
func CheckTargetSurroundingAttackRangeFn(target Player) IsUnblockFn {
	return func(p Point) bool {
		player := target.Game.GetPlayerByPosition(p)
		if player != nil {
			return false
		}

		canHit := target.CanHitPoint(p)
		if !target.WasHit {
			return !canHit
		}

		right := target.FindShooterOnDirection(target.GetDirection().Right())
		if len(right) > 0 {
			rp := right[0]
			canHit = canHit || rp.CanHitPoint(p)
		}

		back := target.FindShooterOnDirection(target.GetDirection().Backward())
		if len(back) > 0 {
			bp := back[0]
			canHit = canHit || bp.CanHitPoint(p)
		}

		front := target.FindShooterOnDirection(target.GetDirection())
		if len(front) > 0 {
			fp := front[0]
			canHit = canHit || fp.CanHitPoint(p)
		}

		left := target.FindShooterOnDirection(target.GetDirection().Left())
		if len(left) > 0 {
			lp := left[0]
			canHit = canHit || lp.CanHitPoint(p)
		}
		return !canHit
	}
}

