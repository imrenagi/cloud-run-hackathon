package main

import "context"


// TODO exploration can check whether it is worth it to go to that point
type Explorer interface {
	Explore(ctx context.Context, p *Player) Move
}

type opponentSorter interface {
	Sort(ctx context.Context) []Player
}

type Exploration struct {
	p *Player
	sorter opponentSorter
}

func (a *Exploration) Explore(ctx context.Context, p *Player) Move {
	ctx, span := tracer.Start(ctx, "Exploration.Explore")
	defer span.End()

	targets := a.sorter.Sort(ctx)
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
		if right != nil {
			canHit = canHit || right.CanHitPoint(ctx, p)
		}

		back := target.FindShooterOnDirection(ctx, target.GetDirection().Backward())
		if back != nil {
			canHit = canHit || back.CanHitPoint(ctx, p)
		}

		front := target.FindShooterOnDirection(ctx, target.GetDirection())
		if front != nil {
			canHit = canHit || front.CanHitPoint(ctx, p)
		}

		left := target.FindShooterOnDirection(ctx, target.GetDirection().Left())
		if left != nil {
			canHit = canHit || left.CanHitPoint(ctx, p)
		}
		return !canHit
	}
}