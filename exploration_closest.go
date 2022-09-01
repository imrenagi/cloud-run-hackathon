package main

import "context"

type ClosestEnemy struct {
}

func (a *ClosestEnemy) Explore(ctx context.Context, p *Player) Move {
	ctx, span := tracer.Start(ctx, "ClosestEnemy.Explore")
	defer span.End()

	targets := p.FindClosestPlayers(ctx)
	if len(targets) == 0 {
		return Throw
	}

	var path Path
	for _, target := range targets {
		aStar := NewAStar(p.Game.Arena,
			WithIsUnblockFn(CheckTargetSurroundingAttackRangeFn(target)),
		)
		var err error
		path, err = aStar.SearchPath(ctx, p.GetPosition(), target.GetPosition())
		if err != nil {
			continue
		}
		if len(path) > 0 {
			break
		}
	}

	if len(path) == 0 {
		return p.Walk(ctx)
	}

	moves := p.RequiredMoves(ctx, path, WithOnlyNextMove())
	if len(moves) > 0 {
		return moves[0]
	} else {
		return p.Walk(ctx)
	}
}


// CheckTargetSurroundingAttackRangeFn checks whether a target and its facing opponent can hit p
func CheckTargetSurroundingAttackRangeFn(target Player) IsUnblockFn {
	return func(ctx context.Context, p Point) bool {
		player := target.Game.GetPlayerByPosition(p)
		if player != nil {
			return false
		}

		canHit := target.CanHitPoint(ctx, p)
		if !target.WasHit {
			return !canHit
		}

		right := target.FindShooterOnDirection(ctx, target.GetDirection().Right())
		if len(right) > 0 {
			rp := right[0]
			canHit = canHit || rp.CanHitPoint(ctx, p)
		}

		back := target.FindShooterOnDirection(ctx, target.GetDirection().Backward())
		if len(back) > 0 {
			bp := back[0]
			canHit = canHit || bp.CanHitPoint(ctx, p)
		}

		front := target.FindShooterOnDirection(ctx, target.GetDirection())
		if len(front) > 0 {
			fp := front[0]
			canHit = canHit || fp.CanHitPoint(ctx, p)
		}

		left := target.FindShooterOnDirection(ctx, target.GetDirection().Left())
		if len(left) > 0 {
			lp := left[0]
			canHit = canHit || lp.CanHitPoint(ctx, p)
		}
		return !canHit
	}
}
