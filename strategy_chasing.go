package main

func NewBrutalChasing(target *Player) *BrutalChasingStrategy {
	return &BrutalChasingStrategy{Target: target}
}

// BrutalChasingStrategy chases the player and ignore any attack.
// This fits for kroco implementation
type BrutalChasingStrategy struct {
	Target *Player
}

func (cs *BrutalChasingStrategy) Play(p *Player) Move {
	p.ChangeState(&Chasing{
		Player:              p,
		Target:              cs.Target,
		ExplorationStrategy: &TargetedEnemy{Target: cs.Target},
	})
	return p.State.Play()
}

func NewSafeChasing(target *Player) *SafeChasingStrategy {
	return &SafeChasingStrategy{Target: target}
}

// SafeChasingStategy attack normally but when exploring try to search for
// target
type SafeChasingStrategy struct {
	Target *Player
}

func (cs *SafeChasingStrategy) Play(p *Player) Move {
	p.ChangeState(&Attack{
		Player:              p,
		ExplorationStrategy: &TargetedEnemy{Target: cs.Target},
	})
	if p.WasHit {
		p.ChangeState(&Escape{Player: p})
	}
	return p.State.Play()
}
