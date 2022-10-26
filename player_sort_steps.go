package main

import (
	"context"
	"sort"
	"sync"
)

type StepSorter struct {
	p Player
}

func (ss StepSorter) Sort(ctx context.Context) []Player {
	ctx, span := tracer.Start(ctx, "Player.SortPlayerByStepsRequired")
	defer span.End()

	var wg sync.WaitGroup
	pairChan := make(chan dPair, len(ss.p.Game.LeaderBoard))

	// TODO Instead of iterating over all players, use bfs??
	for _, ps := range ss.p.Game.LeaderBoard {
		wg.Add(1)

		go func(ps PlayerState) {
			defer wg.Done()
			otherPlayerPt := Point{ps.X, ps.Y}
			if ss.p.GetPosition().Equal(otherPlayerPt) {
				return
			}
			aStar := NewAStar(ss.p.Game.Arena)
			path, err := aStar.SearchPath(ctx, ss.p.GetPosition(), otherPlayerPt)
			if err != nil {
				return
			}
			moves := ss.p.RequiredMoves(ctx, path)

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
		cp := ss.p.Game.GetPlayerByPosition(Point{X: closestPlayer.X, Y: closestPlayer.Y})
		closestPlayers = append(closestPlayers, *cp)
	}
	return closestPlayers
}

