package main

import "context"

func NewSemiBrutalChasing(target *Player) *BrutalChasingStrategy {
	return &BrutalChasingStrategy{Target: target}
}

// SemiBrutalChasingStrategy chases the player and ignore any attack.
// This fits for kroco implementation
type SemiBrutalChasingStrategy struct {
	Target *Player
}

func (cs *SemiBrutalChasingStrategy) Play(ctx context.Context, p *Player) Move {
	ctx, span := tracer.Start(ctx, "SemiBrutalChasingStrategy.Play")
	defer span.End()
	p.ChangeState(&Chasing{
		Player:              p,
		Target:              cs.Target,
		ExplorationStrategy: &TargetedEnemy{Target: cs.Target},
	})
	// TODO check why it is not trying to escape with it is trapped and hit multiple times
	if p.WasHit {
		p.ChangeState(&Escape{Player: p})
	}
	return p.State.Play(ctx)
}

func NewBrutalChasing(target *Player) *BrutalChasingStrategy {
	return &BrutalChasingStrategy{Target: target}
}

// BrutalChasingStrategy chases the player and ignore any attack.
// This fits for kroco implementation
type BrutalChasingStrategy struct {
	Target *Player
}

func (cs *BrutalChasingStrategy) Play(ctx context.Context, p *Player) Move {
	ctx, span := tracer.Start(ctx, "BrutalChasingStrategy.Play")
	defer span.End()
	p.ChangeState(&Chasing{
		Player:              p,
		Target:              cs.Target,
		ExplorationStrategy: &TargetedEnemy{Target: cs.Target},
	})
	return p.State.Play(ctx)
}

func NewSafeChasing(target *Player) *SafeChasingStrategy {
	return &SafeChasingStrategy{Target: target}
}

// SafeChasingStategy attack normally but when exploring try to search for
// target
type SafeChasingStrategy struct {
	Target *Player
}

func (cs *SafeChasingStrategy) Play(ctx context.Context, p *Player) Move {
	ctx, span := tracer.Start(ctx, "SafeChasingStrategy.Play")
	defer span.End()
	p.ChangeState(&Attack{
		Player:              p,
		ExplorationStrategy: &TargetedEnemy{Target: cs.Target},
	})
	if p.WasHit {
		p.ChangeState(&Escape{Player: p})
	}
	return p.State.Play(ctx)
}
