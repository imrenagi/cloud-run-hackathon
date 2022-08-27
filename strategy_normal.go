package main

func DefaultStrategy() *NormalStrategy {
	return &NormalStrategy{}
}

type NormalStrategy struct {
}

func (ns *NormalStrategy) Play(p *Player) Move {
	p.ChangeState(DefaultAttack(p))
	if p.WasHit {
		p.ChangeState(&Escape{Player: p})
	}
	return p.State.Play()
}
