package main

type PlayerState struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	WasHit    bool   `json:"wasHit"`
	Score     int    `json:"score"`
	Game      Game   `json:"-"`
}

func (p PlayerState) Play() Decision {
	var state State = Attack{Player: p}
	if p.WasHit {
		state = Escape{Player: p}
	}
	return state.Play()
}

func (p PlayerState) GetDirection() Direction {
	switch p.Direction {
	case "N":
		return North
	case "W":
		return West
	case "E":
		return East
	default:
		return South
	}
}

func (p PlayerState) GetPosition() Point {
	return Point{
		X: p.X,
		Y: p.Y,
	}
}

func (p PlayerState) Walk() Decision {
	destination := p.GetPosition().TranslateToDirection(1, p.GetDirection())
	if !p.Game.Arena.IsValid(destination) {
		return TurnRight
	}
	// check other player
	players := p.GetPlayersInRange(p.GetDirection(), 1)
	if len(players) > 0 {
		return TurnRight
	}
	return MoveForward
}

func (p *PlayerState) RotateLeft() {
	p.Direction = p.GetDirection().Left().Name
}

func (p *PlayerState) RotateRight() {
	p.Direction = p.GetDirection().Right().Name
}

func (p *PlayerState) setDirection(d Direction) {
	p.Direction = d.Name
}

const attackRange = 3

// FindShooterOnDirection return other players which are in attach range and heading toward the player
func (p PlayerState) FindShooterOnDirection(direction Direction) []PlayerState {
	var filtered []PlayerState
	opponents := p.GetPlayersInRange(direction, attackRange)
	for _, opponent := range opponents {
		// exclude if they are not heading toward the player
		if p.canBeAttackedBy(opponent) {
			filtered = append(filtered, opponent)
		}
	}
	return filtered
}

func (p PlayerState) isMe(p2 PlayerState) bool {
	// TODO Compare with url instead
	return p2.GetPosition().Equal(p.GetPosition())
}

func (p PlayerState) canBeAttackedBy(p2 PlayerState) bool {
	players := p2.GetPlayersInRange(p2.GetDirection(), attackRange)
	for i, player := range players {
		probablyIsAttackingMe := p.isMe(player) && i == 0
		if probablyIsAttackingMe {
			return true
		}
	}
	return false
}

func (p PlayerState) GetPlayersInRange(direction Direction, distance int) []PlayerState {
	var playersInRange []PlayerState
	var ptA = p.GetPosition()
	var ptB = p.GetPosition().TranslateToDirection(distance, direction)

	if ptB.X > p.Game.Arena.Width-1 {
		ptB.X = p.Game.Arena.Width - 1
	}
	if ptB.Y > p.Game.Arena.Height-1 {
		ptB.Y = p.Game.Arena.Height - 1
	}
	if ptB.X < 0 {
		ptB.X = 0
	}
	if ptB.Y < 0 {
		ptB.Y = 0
	}

	for i := 1; i < (distance + 1); i++ {
		npt := ptA.TranslateToDirection(i, direction)
		if !p.Game.Arena.IsValid(npt) {
			break
		}

		if player, ok := p.Game.GetPlayerByPosition(npt); ok {
			playersInRange = append(playersInRange, player)
		}
	}
	return playersInRange
}