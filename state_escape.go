package main

import (
	"context"
	"math"
	"sort"
	"sync"
)

type Escape struct {
	Player *Player
}

const maxHitWhenTrapped int = 3

type attackers []*Player

func (a attackers) Front() *Player {
	return a[0]
}

func (a attackers) Left() *Player {
	return a[1]
}

func (a attackers) Right() *Player {
	return a[2]
}

func (a attackers) Back() *Player {
	return a[3]
}

func (e *Escape) Play(ctx context.Context) Move {
	ctx, span := tracer.Start(ctx, "Escape.Play")
	defer span.End()

	opponents := attackers{
		e.Player.FindShooterOnDirection(ctx, e.Player.GetDirection()),
		e.Player.FindShooterOnDirection(ctx, e.Player.GetDirection().Left()),
		e.Player.FindShooterOnDirection(ctx, e.Player.GetDirection().Right()),
		e.Player.FindShooterOnDirection(ctx, e.Player.GetDirection().Backward()),
	}

	adjacents := e.Player.Game.Arena.GetAdjacent(ctx, e.Player.GetPosition(), WithDiagonalAdjacents(), WithEmptyAdjacent())

	var wg sync.WaitGroup
	pathChan := make(chan Path, len(adjacents))

	for _, adjx := range adjacents {
		wg.Add(1)

		go func(adj Point, opponents []*Player) {
			defer wg.Done()

			var canHit bool
			for _, opponent := range opponents {
				if opponent == nil {
					continue
				}
				if opponent.CanHitPoint(ctx, adj) {
					canHit = true
					break
				}
			}
			if canHit == true {
				return
			}

			// TODO (PENTING) hindari escape ke arah orang lagi perang
			aStar := NewAStar(e.Player.Game.Arena)

			path, err := aStar.SearchPath(ctx, e.Player.GetPosition(), adj)
			if err == ErrPathNotFound {
				return
			}
			pathChan <- path
		}(adjx, opponents)
	}
	wg.Wait()
	close(pathChan)

	var paths []Path
	for elem := range pathChan {
		paths = append(paths, elem)
	}

	// when there is no escape route
	if len(paths) == 0 {
		if opponents.Front() != nil {
			e.Player.trappedCount++
			if e.Player.trappedCount > maxHitWhenTrapped {
				e.Player.trappedCount = 0
				if opponents.Left() != nil {
					return TurnLeft
				}
				if opponents.Right() != nil {
					return TurnRight
				}
			}
			return Throw
		} else if opponents.Left() != nil {
			return TurnLeft
		} else if opponents.Right() != nil {
			return TurnRight
		} else {
			return e.Player.Walk(ctx)
		}
	}

	e.Player.trappedCount = 0
	sort.Sort(byPathLength(paths))

	minPathLength := math.MaxInt
	var shortestPaths []Path
	for _, path := range paths {
		if minPathLength > len(path) {
			minPathLength = len(path)
		}
	}

	for _, path := range paths {
		if len(path) <= minPathLength {
			shortestPaths = append(shortestPaths, path)
		}
	}

	requiredMoves := make([][]Move, len(shortestPaths))
	for idx, aPath := range shortestPaths {
		nextPt := aPath[1]
		moves, err := e.Player.MoveToAdjacent(nextPt)
		if err != nil {
			continue
		}
		requiredMoves[idx] = moves
	}

	if len(requiredMoves) == 0 {
		return e.Player.Walk(ctx)
	}

	mostEfficientMoves := requiredMoves[0]
	for _, sc := range requiredMoves {
		if len(sc) < len(mostEfficientMoves) {
			mostEfficientMoves = sc
		}
	}

	return mostEfficientMoves[0]
}

func (e Escape) GetPlayer() *Player {
	return e.Player
}

type byPathLength []Path

func (a byPathLength) Len() int           { return len(a) }
func (a byPathLength) Less(i, j int) bool { return len(a[i]) < len(a[j]) }
func (a byPathLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type PlayerRetriever interface {
	GetPlayer() *Player
}

type Escaper interface {
	State
	PlayerRetriever
}

// BraveEscapeDecorator use escape logic. But it should be able to attack other component in case of
// only one opponent is attacking
type BraveEscapeDecorator struct {
	Escaper Escaper
}

const maxConsecutiveHitToEscape = 3

func (e *BraveEscapeDecorator) Play(ctx context.Context) Move {
	ctx, span := tracer.Start(ctx, "BraveEscapeDecorator.Play")
	defer span.End()

	player := e.Escaper.GetPlayer()

	if player.consecutiveHitCount > maxConsecutiveHitToEscape {
		return e.Escaper.Play(ctx)
	}

	// TODO latency akan tinggi disini karena harus ngitungin semua direction
	front := player.FindShooterOnDirection(ctx, player.GetDirection())
	back := player.FindShooterOnDirection(ctx, player.GetDirection().Backward())
	left := player.FindShooterOnDirection(ctx, player.GetDirection().Left())
	right := player.FindShooterOnDirection(ctx, player.GetDirection().Right())

	totalShoots := 0

	// TODO kalau gak brave2 bgt, escape setelah 3 kali tembak
	if front != nil{
		totalShoots++
	}
	if left != nil {
		totalShoots++
	}
	if right != nil {
		totalShoots++
	}
	if back != nil {
		totalShoots++
	}

	if (totalShoots == 1 && back != nil) || totalShoots > 1 {
		return e.Escaper.Play(ctx)
	}

	if front != nil {
		return Throw
	} else if left != nil {
		return TurnLeft
	} else if right != nil {
		return TurnRight
	} else {
		return player.Walk(ctx)
	}
}

