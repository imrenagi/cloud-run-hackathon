package main

type Attack struct {
	Player *Player
}

func (a *Attack) Play() Move {
	playersInFront := a.Player.GetPlayersInRange(a.Player.GetDirection(), 3)
	playersInLeft := a.Player.GetPlayersInRange(a.Player.GetDirection().Left(), 3)
	playersInRight := a.Player.GetPlayersInRange(a.Player.GetDirection().Right(), 3)

	if len(playersInFront) > 0 {
		return Throw
	} else if len(playersInLeft) > 0 {
		// TODO check whether opponent is already targeting us
		return TurnLeft
	} else if len(playersInRight) > 0 {
		// TODO check whether opponent is already targeting us
		return TurnRight
	} else {

		// return a.Player.Walk()
		targets := a.Player.FindClosestPlayers()
		if len(targets) == 0 {
			return Throw
		}

		var path Path
		for _, target := range targets {
			aStar := NewAStar(a.Player.Game.Arena,
				WithIsUnblockFn(CheckTargetAttackRangeFn(target)),
			)
			var err error
			path, err = aStar.SearchPath(a.Player.GetPosition(), target.GetPosition())
			if err != nil {
				continue
			}
			if len(path) > 0 {
				break
			}
		}

		if len(path) == 0 {
			return a.Player.Walk()
		}

		moves := a.Player.RequiredMoves(path, WithOnlyNextMove())
		if len(moves) > 0 {
			return moves[0]
		} else {
			return a.Player.Walk()
		}
	}
}

// CheckTargetAttackRangeFn checks whether a target and its facing opponent can hit p
func CheckTargetAttackRangeFn(target Player) IsUnblockFn {
	return func(p Point) bool {
		canHit := target.CanAttack(p)
		if !target.WasHit {
			return !canHit
		}

		right := target.FindShooterOnDirection(target.GetDirection().Right())
		if len(right) > 0 {
			rp := right[0]
			canHit = canHit || rp.CanAttack(p)
		}

		back := target.FindShooterOnDirection(target.GetDirection().Backward())
		if len(back) > 0 {
			bp := back[0]
			canHit = canHit || bp.CanAttack(p)
		}

		// TODO add test
		front := target.FindShooterOnDirection(target.GetDirection())
		if len(front) > 0 {
			fp := front[0]
			canHit = canHit || fp.CanAttack(p)
		}

		// TODO add test
		left := target.FindShooterOnDirection(target.GetDirection().Left())
		if len(left) > 0 {
			lp := left[0]
			canHit = canHit || lp.CanAttack(p)
		}
		return !canHit
	}
}
