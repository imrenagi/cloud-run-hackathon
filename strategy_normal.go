package main

import "context"

// NewNormalStrategy attack normally, and immediately escape when it get hits
func NewNormalStrategy() *NormalStrategy {
	return &NormalStrategy{}
}

type NormalStrategy struct {
}

func (ns *NormalStrategy) Play(ctx context.Context, p *Player) Move {
	ctx, span := tracer.Start(ctx, "NormalStrategy.Play")
	defer span.End()

	if p.WasHit {
		p.ChangeState(&Escape{Player: p})
	} else {
		p.ChangeState(DefaultAttack(p))
	}
	return p.State.Play(ctx)
}

// NewBraveStrategy attacks normally, but when it was hit, it tried to do counter attack
// TODO add counter for hit threshold until it needs to escape
func NewBraveStrategy() *BraveStrategy {
	return &BraveStrategy{}
}

type BraveStrategy struct {}

func (ns *BraveStrategy) Play(ctx context.Context, p *Player) Move {
	ctx, span := tracer.Start(ctx, "BraveStrategy.Play")
	defer span.End()

	if p.WasHit {
		p.ChangeState(&BraveEscapeDecorator{Escaper: &Escape{Player: p}})
	} else {
		p.ChangeState(DefaultAttack(p))
	}
	return p.State.Play(ctx)
}