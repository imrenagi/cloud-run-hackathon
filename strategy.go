package main

type Strategy interface {
	Play() Move
}

func DefaultStrategy(p *Player) Strategy {
	return &NormalStrategy{
		Player: p,
	}
}

type NormalStrategy struct {
	Player *Player
}

func (ns *NormalStrategy) Play() Move {
	var state State = &Attack{Player: ns.Player}
	if ns.Player.WasHit {
		state = &Escape{Player: ns.Player}
	}
	return state.Play()
}

type ChasingStrategy struct {
	Player *Player
	Target *Player
}

func (cs *ChasingStrategy) Play() Move {
	var state State
	if cs.Player.WasHit {
		state = &Escape{Player: cs.Player}
	}
	return state.Play()
}

// lol just idea if we play with swarm bot
type SurroundedStrategy struct {

}