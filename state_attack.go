package main

func DefaultAttack(p *Player) State {
	return &Attack{
		Player:              p,
		ExplorationStrategy: &ClosestEnemy{},
	}
}

type Attack struct {
	Player              *Player
	ExplorationStrategy Exploration
}

func (a *Attack) Play() Move {
	playersInFront := a.Player.GetPlayersInRange(a.Player.GetDirection(), 3)
	playersInLeft := a.Player.GetPlayersInRange(a.Player.GetDirection().Left(), 3)
	playersInRight := a.Player.GetPlayersInRange(a.Player.GetDirection().Right(), 3)

	if len(playersInFront) > 0 {
		return Throw
	} else if len(playersInLeft) > 0 {
		// TODO check whether opponent is already targeting us
		return TurnLeft
	} else if len(playersInRight) > 0 {
		// TODO check whether opponent is already targeting us
		return TurnRight
	} else {
		return a.ExplorationStrategy.Explore(a.Player)
	}
}
