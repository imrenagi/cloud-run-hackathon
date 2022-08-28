package main

import (
	"context"

	"github.com/rs/zerolog/log"
)

func NewChasing(p *Player, target *Player) State {
	return &Chasing{
		Player:              p,
		Target:              target,
		ExplorationStrategy: &TargetedEnemy{},
	}
}

type Chasing struct {
	Player *Player
	Target *Player
	ExplorationStrategy Exploration
}

func (c *Chasing) Play(ctx context.Context) Move {
	ctx, span := tracer.Start(ctx, "Chasing.Play")
	defer span.End()

	if c.Player.CanHitPoint(ctx, c.Target.GetPosition()) {
		log.Debug().Msg("attack")
		return Throw
	}
	return c.ExplorationStrategy.Explore(ctx, c.Player)
}
