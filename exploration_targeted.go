package main

import (
	"context"
)

type TargetedEnemy struct {
	Target *Player
}

func (t *TargetedEnemy) Explore(ctx context.Context, p *Player) Move {
	ctx, span := tracer.Start(ctx, "TargetedEnemy.Explore")
	defer span.End()

	// TODO kalau target di serang sama semua orang, cari aja target lain karena gak perlu kita serang lagi. biar orang lain.

	if t.Target == nil {
		return Throw
	}

	var path Path
	aStar := NewAStar(p.Game.Arena,
		WithIsUnblockFn(CheckTargetSurroundingAttackRangeFn(*t.Target)),
	)
	var err error
	path, err = aStar.SearchPath(ctx, p.GetPosition(), t.Target.GetPosition())
	if err != nil {
		// TODO what happened when path not found
		// continue
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