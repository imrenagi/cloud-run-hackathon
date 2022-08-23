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

}

// lol just idea if we play with swarm bot
type SurroundedStrategy struct {

}