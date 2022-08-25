package main

import (
	"fmt"
	"sort"
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
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	WasHit    bool   `json:"wasHit"`
	Score     int    `json:"score"`
	Game      Game   `json:"-"`

	Strategy     Strategy `json:"-"`
	trappedCount int
}

func (p Player) Clone() Player {
	return Player{
		Name:         p.Name,
		X:            p.X,
		Y:            p.Y,
		Direction:    p.Direction,
		WasHit:       p.WasHit,
		Score:        p.Score,
		Game:         p.Game,
		Strategy:     p.Strategy,
		trappedCount: p.trappedCount,
	}
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

// CanAttack check whether can attack a player in pt
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

func (p *Player) moveForward() {
	newPt := p.GetPosition().TranslateToDirection(1, p.GetDirection())
	if p.Game.Arena.IsValid(newPt) {
		p.setLocation(newPt)
	}
}

func (p *Player) setDirection(d Direction) {
	p.Direction = d.Name
}

func (p *Player) setLocation(pt Point) {
	p.X = pt.X
	p.Y = pt.Y
}

type MoveOption func(o *MoveOptions)

type MoveOptions struct {
	NextMoveOnly bool
}

func WithOnlyNextMove() MoveOption {
	return func(o *MoveOptions) {
		o.NextMoveOnly = true
	}
}

// RequiredMoves return array of moves that should be taken to follow path
func (p Player) RequiredMoves(forPath Path, opts ...MoveOption) []Move {
	options := &MoveOptions{}
	for _, o := range opts {
		o(options)
	}
	var fMoves []Move
	pc := p.Clone()
	for _, pt := range forPath {
		// skip the source path
		if pt.Equal(p.GetPosition()) {
			continue
		}
		moves, err := pc.MoveToAdjacent(pt)
		if err != nil {
			break
		}
		for _, move := range moves {
			fMoves = append(fMoves, move)
			if options.NextMoveOnly && len(fMoves) == 1 {
				return fMoves
			}
			pc.Apply(move)
		}
	}
	return fMoves
}

func (p *Player) Apply(m Move) {
	switch m {
	case WalkForward:
		p.moveForward()
	case TurnRight:
		p.rotateClockwise()
	case TurnLeft:
		p.rotateCounterClockwise()
	}
}

var ErrDestNotFound = fmt.Errorf("target not found")

// MoveToAdjacent return array of moves to reach adjacent cell.
// This only return non empty moves if the toPt is adjacent cell on north, east,
// west, or south.
func (p Player) MoveToAdjacent(toPt Point) ([]Move, error) {
	const distance = 1
	var cCount, ccCount int // clockwise and counter clockwise counter

	p1 := p.Clone()
	p2 := p.Clone()
	var p1Move, p2Move []Move

	var found = false
	for i := 0; i < 4; i++ {
		ptInFront := p1.GetPosition().TranslateToDirection(distance, p1.GetDirection())
		if ptInFront.Equal(toPt) {
			found = true
			break
		}
		p1.rotateCounterClockwise()
		p1Move = append(p1Move, TurnLeft)
		ccCount++
	}

	if !found {
		return nil, ErrDestNotFound
	}

	for i := 0; i < 4; i++ {
		ptInFront := p2.GetPosition().TranslateToDirection(distance, p2.GetDirection())
		if ptInFront.Equal(toPt) {
			break
		}
		p2.rotateClockwise()
		p2Move = append(p2Move, TurnRight)
		cCount++
	}

	if ccCount <= cCount {
		p1Move = append(p1Move, WalkForward)
		return p1Move, nil
	} else {
		p2Move = append(p2Move, WalkForward)
		return p2Move, nil
	}
}

func (p Player) FindClosestPlayers() []Player {
	distanceCalculator := EuclideanDistance{}
	var dPairs []dPair

	for _, ps := range p.Game.Players {
		otherPlayerPt := Point{ps.X, ps.Y}
		if p.GetPosition().Equal(otherPlayerPt) {
			continue
		}

		d := distanceCalculator.Distance(p.GetPosition(), otherPlayerPt)
		dPairs = append(dPairs, dPair{
			distance: d,
			player:   ps,
		})
	}
	if len(dPairs) == 0 {
		return nil
	}

	sort.Sort(byDistance(dPairs))
	var closestPlayers []Player
	for _, dp := range dPairs {
		closestPlayer := dp.player
		cp := p.Game.GetPlayerByPosition(Point{X: closestPlayer.X, Y: closestPlayer.Y})
		closestPlayers = append(closestPlayers, *cp)
	}
	return closestPlayers
}

type dPair struct {
	distance float64
	player   PlayerState
}

type byDistance []dPair

func (a byDistance) Len() int           { return len(a) }
func (a byDistance) Less(i, j int) bool {
	if a[i].distance != a[j].distance {
		return a[i].distance < a[j].distance
	}
	if a[i].player.Y != a[j].player.Y {
		return a[i].player.Y < a[j].player.Y
	}
	return a[i].player.X < a[j].player.X

}
func (a byDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
