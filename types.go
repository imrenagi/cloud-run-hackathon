package main

import (
	"fmt"
	"math"
)

type Direction struct {
	Name   string
	Degree int
}

func (d Direction) Left() Direction {
	degree := d.Degree + 90
	return DirectionMap[degree]
}

func (d Direction) Right() Direction {
	degree := d.Degree - 90
	return DirectionMap[degree]
}

func (d Direction) Opposite() Direction {
	degree := d.Degree + 180
	return DirectionMap[degree]
}

var (
	North = Direction{"N", 180} //
	West  = Direction{"W", 270}
	South = Direction{"S", 0}
	East  = Direction{"E", 90}

	DirectionMap = map[int]Direction{
		-90: West,
		0:   South,
		90:  East,
		180: North,
		270: West,
		360: South,
		450: East,
		540: North,
	}
)

type Game struct {
	Dimension         []int
	PlayersByPosition map[string]PlayerState
}

type ArenaUpdate struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Arena struct {
		Dimensions []int                  `json:"dims"`
		State      map[string]PlayerState `json:"state"`
	} `json:"arena"`
}

func NewGame(a ArenaUpdate) Game {

	playersByPosition := make(map[string]PlayerState)
	for _, v := range a.Arena.State {
		playersByPosition[v.Position().String()] = v
	}

	return Game{
		Dimension:         a.Arena.Dimensions,
		PlayersByPosition: playersByPosition,
	}
}

func (a ArenaUpdate) GetSelf() PlayerState {
	return a.Arena.State[a.Links.Self.Href]
}

type PlayerState struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	WasHit    bool   `json:"wasHit"`
	Score     int    `json:"score"`
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

func (p PlayerState) Position() Point {
	return Point{
		X: p.X,
		Y: p.Y,
	}
}

func (p PlayerState) SearchOpponent() {

}
//
// func (p PlayerState) TurnRight(g Game) Decision {
// 	destination := NewPointAfterMove(p.Position(), p.GetDirection(), 1)
// 	if destination.X < 0 || destination.X > g.Dimension[0] - 1 || destination.Y < 0 || destination.Y > g.Dimension[1] - 1 {
// 		return TurnRight
// 	}
//
// 	//check other player
// 	players := p.GetPlayersInRange(g, p.GetDirection(), 1)
// 	if len(players) > 0 {
// 		return TurnRight
// 	}
//
// 	return TurnRight
// }

func (p PlayerState) MoveForward(g Game) Decision {
	destination := NewPointAfterMove(p.Position(), p.GetDirection(), 1)
	if destination.X < 0 || destination.X > g.Dimension[0] - 1 || destination.Y < 0 || destination.Y > g.Dimension[1] - 1 {
		return TurnRight
	}

	//check other player
	players := p.GetPlayersInRange(g, p.GetDirection(), 1)
	if len(players) > 0 {
		return TurnRight
	}

	return MoveForward
}

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", int(p.X), int(p.Y))
}

func (p Point) Equal(anotherPoint Point) bool {
	return p.X == anotherPoint.X && p.Y == anotherPoint.Y
}

const attackRange = 3

func (p PlayerState) FindShooterFromDirection(g Game, direction Direction) []PlayerState {
	var filtered []PlayerState
	opponents := p.GetPlayersInRange(g, direction, attackRange)
	for _, opponent := range opponents {
		if p.LookingAtMe(g, opponent) {
			filtered = append(filtered, opponent)
		}
	}
	return filtered
}

func (p PlayerState) Escape(g Game) Decision {
	front := p.FindShooterFromDirection(g, p.GetDirection())
	if len(front) > 0 {
		return TurnRight
	}
	back := p.FindShooterFromDirection(g, p.GetDirection().Opposite())
	if len(back) > 0 {
		return TurnRight
	}

	// TODO bug ditembak dari south tapi cuma right terus
	// TODO fix cara cari lawan

	left := p.FindShooterFromDirection(g, p.GetDirection().Left())
	if len(left) > 0 {
		return p.MoveForward(g)
	}
	right := p.FindShooterFromDirection(g, p.GetDirection().Right())
	if len(right) > 0 {
		return p.MoveForward(g)
	}
	return p.MoveForward(g)

}

func (p PlayerState) IsMe(op PlayerState) bool {
	// TODO Compare with url instead
	return op.Position().Equal(p.Position())
}

func (p PlayerState) LookingAtMe(g Game, op PlayerState) bool {
	players := op.GetPlayersInRange(g, op.GetDirection(), attackRange)
	for i, player := range players {
		probablyIsAttackingMe := p.IsMe(player) && i == 0
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
	var ptA = p.Position()
	var ptB = NewPointAfterMove(p.Position(), direction, distance)

	if ptB.X > g.Dimension[0]-1 {
		ptB.X = g.Dimension[0] - 1
	}
	if ptB.Y > g.Dimension[1]-1 {
		ptB.Y = g.Dimension[1] - 1
	}
	if ptB.X < 0 {
		ptB.X = 0
	}
	if ptB.Y < 0 {
		ptB.Y = 0
	}

	for i := 1; i < (distance + 1); i++ {
		npt := NewPointAfterMove(ptA, direction, i)
		if npt.X > g.Dimension[0]-1 || npt.X < 0 {
			break
		}
		if npt.Y > g.Dimension[1]-1 || npt.Y < 0 {
			break
		}

		if player, ok := g.PlayersByPosition[npt.String()]; ok {
			playersInRange = append(playersInRange, player)
		}
	}
	return playersInRange
}

// attack mode
// if there is enemy in attack range, return throw
// else if there is enemy on the left or right, turn
// else return F

type Decision string

const (
	MoveForward Decision = "F"
	TurnRight   Decision = "R"
	TurnLeft    Decision = "L"
	Fight       Decision = "T"
)

type State interface {
	GetDecision(input ArenaUpdate) Decision
}

type Player struct {
}

type Attack interface {
}

func NewPointAfterMove(p Point, direction Direction, distance int) Point {
	dir := float64(direction.Degree) * math.Pi / 180
	return Point{
		X: p.X + distance*int(math.Round(math.Sin(dir))),
		Y: p.Y + distance*int(math.Round(math.Cos(dir))),
	}
}
