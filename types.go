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

var (
	North = Direction{"N", 180} //
	West  = Direction{"W", 270}
	South = Direction{"S", 0}
	East  = Direction{"E", 90}

	DirectionMap = map[int]Direction{
		-90: West,
		0:   South,
		90:  East,
		270: West,
		180: North,
		360: South,
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

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", int(p.X), int(p.Y))
}

const attackRange = 3

func (p PlayerState) GetPlayersInDirection(g Game, direction Direction) int {
	var enemyInAttackRange []Point
	var ptA = p.Position()
	var ptB = NewPointAfterMove(p.Position(), direction, attackRange)

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

	for i := 1; i < 4; i++ {
		npt := NewPointAfterMove(ptA, direction, i)
		if npt.X > g.Dimension[0]-1 || npt.X < 0 {
			break
		}
		if npt.Y > g.Dimension[1]-1 || npt.Y < 0 {
			break
		}

		if player, ok := g.PlayersByPosition[npt.String()]; ok {
			enemyInAttackRange = append(enemyInAttackRange, Point{
				X: player.X,
				Y: player.Y,
			})
		}
	}

	// if ptA.X == ptB.X {
	// 	// iterate over y
	// 	minY := ptA.Y
	// 	maxY := ptB.Y
	// 	if ptA.X > ptB.X {
	// 		maxY = ptA.Y
	// 		minY = ptB.Y
	// 	}
	//
	// 	for y := minY + 1; y <= maxY; y++ {
	// 		pt := Point{X: ptA.X, Y: y}
	// 		if player, ok := g.PlayersByPosition[pt.String()]; ok {
	// 			enemyInAttackRange = append(enemyInAttackRange, Point{
	// 				X: player.X,
	// 				Y: player.Y,
	// 			})
	// 		}
	// 	}
	// } else if ptA.Y == ptB.Y {
	// 	// iterate over X
	// 	minX := ptA.X
	// 	maxX := ptB.X
	// 	if ptA.X > ptB.X {
	// 		maxX = ptA.X
	// 		minX = ptB.X
	// 	}
	//
	// 	for x := minX + 1; x <= maxX; x++ {
	// 		pt := Point{X: x, Y: ptA.Y}
	// 		if player, ok := g.PlayersByPosition[pt.String()]; ok {
	// 			enemyInAttackRange = append(enemyInAttackRange, Point{
	// 				X: player.X,
	// 				Y: player.Y,
	// 			})
	// 		}
	// 	}
	// }

	return len(enemyInAttackRange)
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
