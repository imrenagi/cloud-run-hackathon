package main

import "context"

func DefaultStrategy() *NormalStrategy {
	return &NormalStrategy{}
}

type NormalStrategy struct {
}

func (ns *NormalStrategy) Play(ctx context.Context, p *Player) Move {
	ctx, span := tracer.Start(ctx, "NormalStrategy.Play")
	defer span.End()

	p.ChangeState(DefaultAttack(p))
	if p.WasHit {
		p.ChangeState(&Escape{Player: p})
	}
	return p.State.Play(ctx)
}
