package main

func NewChasing(p *Player, target *Player) State {
	return &Chasing{
		Player:              p,
		Target:              target,
		ExplorationStrategy: &TargetedEnemy{},
	}
}

type Chasing struct {
	Player *Player
	Target *Player
	ExplorationStrategy Exploration
}

func (c *Chasing) Play() Move {
	if c.Player.CanHitPoint(c.Target.GetPosition()) {
		return Throw
	}
	return c.ExplorationStrategy.Explore(c.Player)
}

