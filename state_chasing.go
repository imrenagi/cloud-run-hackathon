package main

import (
	"context"
)

func NewChasing(p *Player, target *Player) State {
	return &Chasing{
		Player:              p,
		Target:              target,
		ExplorationStrategy: &TargetedEnemy{},
	}
}

type Chasing struct {
	Player              *Player
	Target              *Player
	ExplorationStrategy Explorer
}

func (c *Chasing) Play(ctx context.Context) Move {
	ctx, span := tracer.Start(ctx, "Chasing.Play")
	defer span.End()

	if c.Player.CanHitPoint(ctx, c.Target.GetPosition()) {
		return Throw
	}
	return c.ExplorationStrategy.Explore(ctx, c.Player)
}
