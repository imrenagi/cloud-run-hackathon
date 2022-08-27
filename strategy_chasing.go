package main

func NewChasingStrategy(target *Player) *ChasingStrategy {
	return &ChasingStrategy{Target: target}
}

type ChasingStrategy struct {
	Target *Player
}

func (cs *ChasingStrategy) Play(p *Player) Move {

	// TODO attacknya normal, tapi explore target.
	p.ChangeState(&Chasing{
		Player:              p,
		Target:              cs.Target,
		ExplorationStrategy: &TargetedEnemy{Target: cs.Target},
	})
	// if p.WasHit {
	// 	p.ChangeState(&Escape{Player: p})
	// }
	return p.State.Play()
}
