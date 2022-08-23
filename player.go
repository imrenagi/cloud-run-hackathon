package main

import (
	"fmt"
)

func NewPlayerWithUrl(url string, state PlayerState) *Player {
	p := NewPlayer(state)
	p.Name = url
	p.Strategy = DefaultStrategy(p)
	return p
}

func NewPlayer(state PlayerState) *Player {
	p := &Player{
		X:         state.X, // TODO ubah jadi location
		Y:         state.Y,
		Direction: state.Direction, // TODO ubah jadi direction
		WasHit:    state.WasHit,
		Score:     state.Score,
	}
	p.Strategy = DefaultStrategy(p)
	return p
}

type Player struct {
	Name      string
	X         int      `json:"x"`
	Y         int      `json:"y"`
	Direction string   `json:"direction"`
	WasHit    bool     `json:"wasHit"`
	Score     int      `json:"score"`
	Game      Game     `json:"-"`

	Strategy  Strategy `json:"-"`
	trappedCount int
}

func (p *Player) Play() Move {
	return p.Strategy.Play()
}

func (p Player) GetDirection() Direction {
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

func (p Player) GetPosition() Point {
	return Point{
		X: p.X,
		Y: p.Y,
	}
}

func (p Player) Walk() Move {
	destination := p.GetPosition().TranslateToDirection(1, p.GetDirection())
	if !p.Game.Arena.IsValid(destination) {
		return TurnRight
	}
	// check other player
	players := p.GetPlayersInRange(p.GetDirection(), 1)
	if len(players) > 0 {
		return TurnRight
	}
	return WalkForward
}

const attackRange = 3

// FindShooterOnDirection return other players which are in attach range and heading toward the player
func (p Player) FindShooterOnDirection(direction Direction) []Player {
	var filtered []Player
	opponents := p.GetPlayersInRange(direction, attackRange)
	for _, opponent := range opponents {
		// exclude if they are not heading toward the player
		if p.canBeAttackedBy(opponent) {
			filtered = append(filtered, opponent)
		}
	}
	return filtered
}

func (p Player) CanAttack(pt Point) bool {
	var ptA = p.GetPosition()
	var ptB = p.GetPosition().TranslateToDirection(attackRange, p.GetDirection())

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

	for i := 1; i < (attackRange + 1); i++ {
		npt := ptA.TranslateToDirection(i, p.GetDirection())
		if !p.Game.Arena.IsValid(npt) {
			break
		}

		if npt.X == pt.X && npt.Y == pt.Y {
			return true
		}
	}
	return false
}

func (p Player) isMe(p2 Player) bool {
	// TODO Compare with url instead
	return p2.GetPosition().Equal(p.GetPosition())
}

func (p Player) canBeAttackedBy(p2 Player) bool {
	players := p2.GetPlayersInRange(p2.GetDirection(), attackRange)
	for i, player := range players {
		probablyIsAttackingMe := p.isMe(player) && i == 0
		if probablyIsAttackingMe {
			return true
		}
	}
	return false
}

func (p Player) GetPlayersInRange(direction Direction, distance int) []Player {
	var playersInRange []Player
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

		if player := p.Game.GetPlayerByPosition(npt); player != nil {
			playersInRange = append(playersInRange, *player)
		}
	}
	return playersInRange
}

func (p *Player) rotateCounterClockwise() {
	p.Direction = p.GetDirection().Left().Name
}

func (p *Player) rotateClockwise() {
	p.Direction = p.GetDirection().Right().Name
}

func (p *Player) setDirection(d Direction) {
	p.Direction = d.Name
}

var ErrDestNotFound = fmt.Errorf("target not found")

// MoveNeededToReachAdjacent return array of moves to reach adjacent cell
func (p Player) MoveNeededToReachAdjacent(toPt Point) ([]Move, error) {
	myPt := Point{X: p.X, Y: p.Y}
	const distance = 1
	var cCount, ccCount int // clockwise and counter clockwise counter
	initialDirection := p.GetDirection()

	var found = false
	for i := 0; i < 4; i++ {
		ptInFront := myPt.TranslateToDirection(distance, p.GetDirection())
		if ptInFront.Equal(toPt) {
			found = true
			break
		}
		p.rotateCounterClockwise()
		ccCount++
	}
	p.setDirection(initialDirection)

	if !found {
		return nil, ErrDestNotFound
	}

	for i := 0; i < 4; i++ {
		ptInFront := myPt.TranslateToDirection(distance, p.GetDirection())
		if ptInFront.Equal(toPt) {
			break
		}
		p.rotateClockwise()
		cCount++
	}
	p.setDirection(initialDirection)

	var rotationDecision []Move
	minRotationCount := ccCount
	for i := 0; i < ccCount; i++ {
		p.rotateCounterClockwise()
		rotationDecision = append(rotationDecision, TurnLeft)
	}

	if minRotationCount > cCount {
		rotationDecision = []Move{}
		minRotationCount = cCount
		p.setDirection(initialDirection)
		for i := 0; i < cCount; i++ {
			rotationDecision = append(rotationDecision, TurnRight)
			p.rotateClockwise()
		}
	}

	// add move forward
	rotationDecision = append(rotationDecision, WalkForward)

	return rotationDecision, nil
}
