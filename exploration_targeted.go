package main

import "github.com/rs/zerolog/log"

type TargetedEnemy struct {
	Target *Player
}

func (t *TargetedEnemy) Explore(p *Player) Move {

	// TODO kalau target di serang sama semua orang, cari aja target lain karena gak perlu kita serang lagi. biar orang lain.

	if t.Target == nil {
		// TODO what happened when target is nil
		return Throw
	}
	log.Debug().
		Str("name", t.Target.Name).
		Int("x", t.Target.X).
		Int("y", t.Target.Y).
		Msgf("target")

	var path Path
	aStar := NewAStar(p.Game.Arena,
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

func ObstacleMapFn(player *Player) IsUnblockFn {
	return func(p Point) bool {
		if !player.Game.Arena.IsValid(p) {
			return false
		}
		obstacleMap := player.Game.ObstacleMap()
		return obstacleMap[p.Y][p.X]
	}
}
