package main


type PlayerState struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	WasHit    bool   `json:"wasHit"`
	Score     int    `json:"score"`
}

func (p PlayerState) Play(g Game) Decision {
	var state State = Attack{Player: p}
	if p.WasHit {
		state = Escape{Player: p}
	}
	return state.Play(g)
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

func (p PlayerState) Walk(g Game) Decision {
	destination := p.GetPosition().MoveWithDirection(1, p.GetDirection())
	if destination.X < 0 || destination.X > g.Arena.Width-1 || destination.Y < 0 || destination.Y > g.Arena.Height-1 {
		return TurnRight
	}
	// check other player
	players := p.GetPlayersInRange(g, p.GetDirection(), 1)
	if len(players) > 0 {
		return TurnRight
	}
	return MoveForward
}

const attackRange = 3

// FindShooterOnDirection return other players which are in attach range and heading toward the player
func (p PlayerState) FindShooterOnDirection(g Game, direction Direction) []PlayerState {
	var filtered []PlayerState
	opponents := p.GetPlayersInRange(g, direction, attackRange)
	for _, opponent := range opponents {
		// exclude if they are not heading toward the player
		if p.isHeadingTowardMe(g, opponent) {
			filtered = append(filtered, opponent)
		}
	}
	return filtered
}

func (p PlayerState) isMe(p2 PlayerState) bool {
	// TODO Compare with url instead
	return p2.GetPosition().Equal(p.GetPosition())
}

func (p PlayerState) isHeadingTowardMe(g Game, p2 PlayerState) bool {
	players := p2.GetPlayersInRange(g, p2.GetDirection(), attackRange)
	for i, player := range players {
		probablyIsAttackingMe := p.isMe(player) && i == 0
		if probablyIsAttackingMe {
			return true
		}
	}
	return false
}

// Test Cases:
// * persis disebelah
// * ada user ditengah, mestinya ini gak lookingatme

func (p PlayerState) GetPlayersInRange(g Game, direction Direction, distance int) []PlayerState {
	var playersInRange []PlayerState
	var ptA = p.GetPosition()
	var ptB = p.GetPosition().MoveWithDirection(distance, direction)

	if ptB.X > g.Arena.Width-1 {
		ptB.X = g.Arena.Width - 1
	}
	if ptB.Y > g.Arena.Height-1 {
		ptB.Y = g.Arena.Height - 1
	}
	if ptB.X < 0 {
		ptB.X = 0
	}
	if ptB.Y < 0 {
		ptB.Y = 0
	}

	for i := 1; i < (distance + 1); i++ {
		npt := ptA.MoveWithDirection(i, direction)
		if npt.X > g.Arena.Width-1 || npt.X < 0 {
			break
		}
		if npt.Y > g.Arena.Height-1 || npt.Y < 0 {
			break
		}

		if player, ok := g.PlayersByPosition[npt.String()]; ok {
			playersInRange = append(playersInRange, player)
		}
	}
	return playersInRange
}

