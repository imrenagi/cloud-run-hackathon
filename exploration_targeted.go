package main

type TargetedEnemy struct {
	Target *Player
}

func (t *TargetedEnemy) Explore(p *Player) Move {

	// kalau target di serang sama semua orang, cari aja target lain

	if t.Target == nil {
		// TODO what happened when target is nil
		return Throw
	}

	var path Path
	aStar := NewAStar(p.Game.Arena,
		// TODO this must try to find the safest path as possible
		WithIsUnblockFn(CheckTargetSurroundingAttackRangeFn(*t.Target)),
	)
	var err error
	path, err = aStar.SearchPath(p.GetPosition(), t.Target.GetPosition())
	if err != nil {
		// TODO what happened when path not found
		// continue
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
