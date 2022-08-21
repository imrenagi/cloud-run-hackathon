package main

import (
	"fmt"
	"math"
)

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

func (p PlayerState) PlanDecision(path Path) []Decision {
	if path == nil || len(path) == 0 {
		return nil
	}




	return nil
}

func (p *PlayerState) rotateCounterClockwise() {
	p.Direction = p.GetDirection().Left().Name
}

func (p *PlayerState) rotateClockwise() {
	p.Direction = p.GetDirection().Right().Name
}

func (p *PlayerState) setDirection(d Direction) {
	p.Direction = d.Name
}

var ErrDestNotFound = fmt.Errorf("target not found")

// GetShortestRotation return decision to turn to change direction to toPt
func (p PlayerState) GetShortestRotation(toPt Point) ([]Decision, error) {

	// myPt := Point{X: p.X, Y: p.Y}
	// const distance = 1
	// ptInFront := myPt.TranslateToDirection(distance, p.GetDirection())
	//
	// vec1 := NewVector(myPt, ptInFront)
	// vec2 := NewVector(myPt, toPt)

	myPt := Point{X: p.X, Y: p.Y}
	const distance = 1
	var cCount, ccCount int // clockwise and counter clockwise counter
	initialDirection := p.GetDirection()

	var found = false
	for i := 0; i<4; i++ {
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

	for i := 0; i<4; i++ {
		ptInFront := myPt.TranslateToDirection(distance, p.GetDirection())
		if ptInFront.Equal(toPt) {
			break
		}
		p.rotateClockwise()
		cCount++
	}
	p.setDirection(initialDirection)

	var rotationDecision []Decision
	minRotationCount := ccCount
	for i := 0; i<ccCount; i++ {
		p.rotateCounterClockwise()
		rotationDecision = append(rotationDecision, TurnLeft)
	}

	if minRotationCount > cCount {
		rotationDecision = []Decision{}
		minRotationCount = cCount
		p.setDirection(initialDirection)
		for i := 0; i<cCount; i++ {
			rotationDecision = append(rotationDecision, TurnRight)
			p.rotateClockwise()
		}
	}

	return rotationDecision, nil
}

func NewVector(p1, p2 Point) Vector {
	return Vector{
		X: float64(p2.X - p1.X),
		Y: float64(p2.Y - p1.Y),
	}
}

type Vector struct {
	X, Y float64
}

func (v Vector) Angle(v2 Vector) float64 {
	nom := v.X*v2.Y - v2.X*v.Y
	denom := math.Sqrt(v.X*v.X + v.Y*v.Y) * math.Sqrt(v2.X*v2.X + v2.Y*v2.Y)
	angleInRad := math.Asin(nom/denom)
	return 180 * angleInRad / math.Pi
}


// TODO fix obstacle logic in a star
// TODO translate shortest path to Decision
