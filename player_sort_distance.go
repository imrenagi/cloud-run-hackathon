package main

import (
	"context"
	"sort"
)

type DistanceSorter struct {
	p Player
}

func (ds DistanceSorter) Sort(ctx context.Context) []Player {
	ctx, span := tracer.Start(ctx, "Player.SortPlayerByDistance")
	defer span.End()

	distanceCalculator := EuclideanDistance{}
	var dPairs []dPair

	for _, ps := range ds.p.Game.LeaderBoard {
		otherPlayerPt := Point{ps.X, ps.Y}
		if ds.p.GetPosition().Equal(otherPlayerPt) {
			continue
		}

		d := distanceCalculator.Distance(ds.p.GetPosition(), otherPlayerPt)
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
		cp := ds.p.Game.GetPlayerByPosition(Point{X: closestPlayer.X, Y: closestPlayer.Y})
		closestPlayers = append(closestPlayers, *cp)
	}
	return closestPlayers
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


