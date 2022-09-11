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

	front := a.Player.FindTargetOnDirection(ctx, a.Player.GetDirection())
	if front != nil {
		return Throw
	}

	left := a.Player.FindTargetOnDirection(ctx, a.Player.GetDirection().Left())
	right := a.Player.FindTargetOnDirection(ctx, a.Player.GetDirection().Right())
	if left != nil && right != nil {
		if left.Score > right.Score {
			return TurnLeft
		} else {
			return TurnRight
		}
	}

	if left != nil {
		// TODO check whether opponent is already targeting us
		return TurnLeft
	}

	if right != nil {
		// TODO check whether opponent is already targeting us
		return TurnRight
	}

	// TODO attack should be able to use closest/targeted
	// TODO attack sambil maju satu langkah biar lawan terjepit
	return a.ExplorationStrategy.Explore(ctx, a.Player)
}
