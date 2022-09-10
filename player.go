package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

func NewPlayerWithUrl(url string, state PlayerState) *Player {
	p := NewPlayer(state)
	p.Name = url
	p.Strategy = NewNormalStrategy()
	return p
}

func NewPlayer(state PlayerState) *Player {

	whitelisted := make(map[string]string)
	urls := os.Getenv("WHITELISTED_URLS")
	for _, url := range strings.Split(urls, ",") {
		whitelisted[url] = url
	}

	p := &Player{
		X:           state.X,
		Y:           state.Y,
		Direction:   state.Direction,
		WasHit:      state.WasHit,
		Score:       state.Score,
		Whitelisted: whitelisted,
	}
	p.Strategy = NewNormalStrategy()
	return p
}

type Player struct {
	Name        string
	X           int      `json:"x"`
	Y           int      `json:"y"`
	Direction   string   `json:"direction"`
	WasHit      bool     `json:"wasHit"`
	Score       int      `json:"score"`
	Game        Game     `json:"-"`
	State       State    `json:"-"`
	Strategy    Strategy `json:"-"`
	Whitelisted map[string]string

	trappedCount        int
	consecutiveHitCount int
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

func (p *Player) Play(ctx context.Context) Move {
	ctx, span := tracer.Start(ctx, "Player.Play")
	defer span.End()

	// TODO Calculate priority whether to attack or to chase
	return p.Strategy.Play(ctx, p)
}

func (p *Player) ChangeState(s State) {
	p.State = s
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

func (p Player) Walk(ctx context.Context) Move {
	ctx, span := tracer.Start(ctx, "Player.Walk")
	defer span.End()

	destination := p.GetPosition().TranslateToDirection(1, p.GetDirection())
	if !p.Game.Arena.IsValid(destination) {
		return TurnRight
	}
	// check other player
	players := p.GetPlayersInRange(ctx, p.GetDirection(), 1)
	if len(players) > 0 {
		return TurnRight
	}
	return WalkForward
}

const attackRange = 3

// GetLowestRank returned players that has lowest rank, exclude whitelisted
func (p Player) GetLowestRank(ctx context.Context) *Player {
	ctx, span := tracer.Start(ctx, "Player.GetHighestRank")
	defer span.End()

	for idx := len(p.Game.LeaderBoard) - 1; idx >= 0; idx-- {
		ps := p.Game.LeaderBoard[idx]

		target := p.Game.GetPlayerByPosition(Point{ps.X, ps.Y})
		if target == nil {
			continue
		}
		_, ok := p.Whitelisted[ps.URL]
		if ok {
			continue
		}

		if p.IsMe(target) {
			continue
		}

		if target != nil {
			return target
		}
	}
	return nil
}

func (p Player) IsMe(ap *Player) bool {
	return p.Name == ap.Name
}

// GetHighestRank returned players that has highest rank, exclude whitelisted
func (p Player) GetHighestRank(ctx context.Context) *Player {
	ctx, span := tracer.Start(ctx, "Player.GetHighestRank")
	defer span.End()

	for _, ps := range p.Game.LeaderBoard {
		target := p.Game.GetPlayerByPosition(Point{ps.X, ps.Y})
		if target == nil {
			continue
		}
		_, ok := p.Whitelisted[ps.URL]
		if ok {
			continue
		}
		return target
	}
	return nil
}

func (p Player) GetPlayerOnNextPodium(ctx context.Context) *Player {
	ctx, span := tracer.Start(ctx, "Player.GetPlayerOnNextPodium")
	defer span.End()

	myRank := p.Game.LeaderBoard.GetRank(p)
	if myRank == 0 {
		return nil
	}
	ps := p.Game.LeaderBoard.GetPlayerByRank(myRank - 1)
	return p.Game.GetPlayerByPosition(Point{ps.X, ps.Y})
}

// FindShooterOnDirection return other players which are in attach range and heading toward the player
func (p Player) FindShooterOnDirection(ctx context.Context, direction Direction) *Player {
	ctx, span := tracer.Start(ctx, "Player.FindShooterOnDirection")
	defer span.End()

	var shooter *Player
	opponents := p.GetPlayersInRange(ctx, direction, attackRange)
	for _, opponent := range opponents {
		// assumming the first opponent is the closest one
		if opponent.CanHit(ctx, p) {
			// filtered = append(filtered, opponent)
			shooter = &opponent
			break
		}
	}
	return shooter
}

type HitOption func(*HitOptions)

type HitOptions struct {
	IgnorePlayer bool
}

func WithIgnorePlayer() HitOption {
	return func(options *HitOptions) {
		options.IgnorePlayer = true
	}
}

// CanHitPoint check whether can attack a player in pt
func (p Player) CanHitPoint(ctx context.Context, pt Point, opts ...HitOption) bool {
	ctx, span := tracer.Start(ctx, "Player.CanHitPoint")
	defer span.End()

	options := &HitOptions{}
	for _, o := range opts {
		o(options)
	}

	var ptA = p.GetPosition()
	for i := 1; i < (attackRange + 1); i++ {
		npt := ptA.TranslateToDirection(i, p.GetDirection())
		if !p.Game.Arena.IsValid(npt) {
			break
		}
		if npt.X == pt.X && npt.Y == pt.Y {
			return true
		}

		if !options.IgnorePlayer {
			pl := p.Game.GetPlayerByPosition(npt)
			if pl != nil {
				return false
			}
		}

	}
	return false
}

func (p Player) CanHit(ctx context.Context, p2 Player, opts ...HitOption) bool {
	ctx, span := tracer.Start(ctx, "Player.CanHit")
	defer span.End()
	return p.CanHitPoint(ctx, p2.GetPosition(), opts...)
}

func (p Player) GetRank() int {
	return p.Game.LeaderBoard.GetRank(p)
}

func (p Player) FindTargetOnDirection(ctx context.Context, direction Direction) *Player {
	players := p.GetPlayersInRange(ctx, direction, attackRange)
	if len(players) == 0 {
		return nil
	}
	target := players[0]
	return &target
}

func (p Player) GetPlayersInRange(ctx context.Context, direction Direction, distance int) []Player {
	ctx, span := tracer.Start(ctx, "Player.GetPlayersInRange")
	defer span.End()

	var playersInRange []Player
	var ptA = p.GetPosition()
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
	p.setDirection(p.GetDirection().Left())
}

func (p *Player) rotateClockwise() {
	p.setDirection(p.GetDirection().Right())
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
func (p Player) RequiredMoves(ctx context.Context, forPath Path, opts ...MoveOption) []Move {
	ctx, span := tracer.Start(ctx, "Player.RequiredMoves")
	defer span.End()

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

func (p Player) FindClosestPlayers2(ctx context.Context) []Player {
	ctx, span := tracer.Start(ctx, "Player.FindClosestPlayers")
	defer span.End()

	distanceCalculator := EuclideanDistance{}
	var dPairs []dPair

	for _, ps := range p.Game.LeaderBoard {
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

func (p Player) FindClosestPlayers(ctx context.Context) []Player {
	ctx, span := tracer.Start(ctx, "Player.FindClosestPlayers2")
	defer span.End()

	var wg sync.WaitGroup
	pairChan := make(chan dPair, len(p.Game.LeaderBoard))

	for _, ps := range p.Game.LeaderBoard {
		wg.Add(1)

		go func(ps PlayerState) {
			defer wg.Done()
			otherPlayerPt := Point{ps.X, ps.Y}
			if p.GetPosition().Equal(otherPlayerPt) {
				return
			}
			aStar := NewAStar(p.Game.Arena)
			path, err := aStar.SearchPath(ctx, p.GetPosition(), otherPlayerPt)
			if err != nil {
				return
			}
			moves := p.RequiredMoves(ctx, path)

			pairChan <- dPair{
				distance: float64(len(moves)),
				player:   ps,
			}
		}(ps)
	}
	wg.Wait()
	close(pairChan)

	var dPairs []dPair
	for elem := range pairChan {
		dPairs = append(dPairs, elem)
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

func (p *Player) UpdateHitCount() {
	if p.WasHit {
		p.consecutiveHitCount++
	} else {
		p.consecutiveHitCount = 0
	}
}

type dPair struct {
	distance float64
	player   PlayerState
}

type byDistance []dPair

func (a byDistance) Len() int { return len(a) }
func (a byDistance) Less(i, j int) bool {
	if a[i].distance != a[j].distance {
		return a[i].distance < a[j].distance
	}
	if a[i].player.Y != a[j].player.Y {
		return a[i].player.Y < a[j].player.Y
	}
	return a[i].player.X < a[j].player.X
}
func (a byDistance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
