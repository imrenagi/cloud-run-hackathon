package main

import "context"

// DefaultAttack attack closest players
func DefaultAttack(p *Player) State {
	return &Attack{
		Player:              p,
		ExplorationStrategy: &ClosestEnemy{},
	}
}

// TargetedAttack should attack in normal cases, but when exploring it tries to search for the target
func TargetedAttack(p *Player) State {
	return &Attack{
		Player:              p,
		ExplorationStrategy: &TargetedEnemy{},
	}
}

type Attack struct {
	Player              *Player
	ExplorationStrategy Exploration
}

func (a *Attack) Play(ctx context.Context) Move {
	ctx, span := tracer.Start(ctx, "Attack.Play")
	defer span.End()

	playersInFront := a.Player.GetPlayersInRange(ctx, a.Player.GetDirection(), attackRange)
	playersInLeft := a.Player.GetPlayersInRange(ctx, a.Player.GetDirection().Left(), attackRange)
	playersInRight := a.Player.GetPlayersInRange(ctx, a.Player.GetDirection().Right(), attackRange)

	if len(playersInFront) > 0 {
		return Throw
	} else if len(playersInLeft) > 0 {
		// TODO check whether opponent is already targeting us
		return TurnLeft
	} else if len(playersInRight) > 0 {
		// TODO check whether opponent is already targeting us
		return TurnRight
	} else {
		// TODO attack should be able to use closest/targeted
		return a.ExplorationStrategy.Explore(ctx, a.Player)
	}
	// TODO attack sambil maju satu langkah
}
