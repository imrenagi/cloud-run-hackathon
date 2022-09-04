package main

import (
	"context"
	"math"
	"sort"
)

type Escape struct {
	Player *Player
}

const maxHitWhenTrapped int = 3

func (e *Escape) Play(ctx context.Context) Move {
	ctx, span := tracer.Start(ctx, "Escape.Play")
	defer span.End()

	front := e.Player.FindShooterOnDirection(ctx, e.Player.GetDirection())
	// back := e.Player.FindShooterOnDirection(ctx, e.Player.GetDirection().Backward())
	left := e.Player.FindShooterOnDirection(ctx, e.Player.GetDirection().Left())
	right := e.Player.FindShooterOnDirection(ctx, e.Player.GetDirection().Right())

	// totalShoots := 0
	//
	// // TODO escape state (brave mode) tapi kalau cuma ada satu orang yg nembak, arahin ke dia terus tembak balik
	// if len(front) > 0 {
	// 	totalShoots++
	// }
	// if len(back) > 0 {
	// 	totalShoots++
	// }
	// if len(left) > 0 {
	// 	totalShoots++
	// }
	// if len(right) > 0 {
	// 	totalShoots++
	// }
	//
	// if totalShoots <= 1 {
	// 	// escape
	// 	// TODO redirect to the enemy
	// 	// if already facing, shot
	// }
	//
	// TODO hindari escape ke arah orang lagi perang

	var paths []Path // list of possible path
	validAdjacent := e.Player.Game.Arena.GetAdjacent(ctx, e.Player.GetPosition(), WithDiagonalAdjacents(), WithEmptyAdjacent())

	if front != nil {
		var newAdjacent []Point
		for _, adj := range validAdjacent {
			if !front.CanHitPoint(ctx, adj) {
				newAdjacent = append(newAdjacent, adj)
			}
		}
		validAdjacent = newAdjacent
	}
	if left != nil {
		var newAdjacent []Point
		for _, adj := range validAdjacent {
			if !left.CanHitPoint(ctx, adj) {
				newAdjacent = append(newAdjacent, adj)
			}
		}
		validAdjacent = newAdjacent
	}
	if right != nil {
		var newAdjacent []Point
		for _, adj := range validAdjacent {
			if !right.CanHitPoint(ctx, adj) {
				newAdjacent = append(newAdjacent, adj)
			}
		}
		validAdjacent = newAdjacent
	}
	for _, adj := range validAdjacent {
		aStar := NewAStar(e.Player.Game.Arena)
		path, err := aStar.SearchPath(ctx, e.Player.GetPosition(), adj)
		if err == ErrPathNotFound {
			continue
		}
		paths = append(paths, path)
	}

	// when there is no escape route
	if len(paths) == 0 {
		if front != nil {
			e.Player.trappedCount++
			if e.Player.trappedCount > maxHitWhenTrapped {
				e.Player.trappedCount = 0
				if left != nil {
					return TurnLeft
				}
				if right != nil {
					return TurnRight
				}
			}
			return Throw
		} else if left != nil {
			return TurnLeft
		} else if right != nil {
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

type byPathLength []Path

func (a byPathLength) Len() int           { return len(a) }
func (a byPathLength) Less(i, j int) bool { return len(a[i]) < len(a[j]) }
func (a byPathLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
