package main

import "context"

func NewNormalStrategy() *NormalStrategy {
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

func NewBraveStrategy() *BraveStrategy {
	return &BraveStrategy{}
}

type BraveStrategy struct {}

func (ns *BraveStrategy) Play(ctx context.Context, p *Player) Move {
	ctx, span := tracer.Start(ctx, "BraveStrategy.Play")
	defer span.End()

	p.ChangeState(DefaultAttack(p))
	if p.WasHit {
		p.ChangeState(&BraveEscapeDecorator{Escaper: &Escape{Player: p}})
	}
	return p.State.Play(ctx)
}