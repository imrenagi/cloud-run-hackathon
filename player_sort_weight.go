package main

import (
	"context"
	"sort"

	"github.com/rs/zerolog/log"
)

type WeightedSorter struct {
	p Player
}

func (ws WeightedSorter) Sort(ctx context.Context) []Player {
	ctx, span := tracer.Start(ctx, "WeightedSorter.Sort")
	defer span.End()

	ws.p.Game.LeaderBoard.Sort()
	maxScore := ws.p.Game.LeaderBoard[0].Score
	minScore := ws.p.Game.LeaderBoard[len(ws.p.Game.LeaderBoard)-1].Score
	minDistance := float64(0)
	maxDistance := ws.p.Game.Arena.Diagonal()

	distanceCalculator := EuclideanDistance{}

	var wps []weightedPlayer

	tmpMaxNormalizedScore := float64(-1)
	tmpMaxNormalizedDistance := float64(-1)
	for _, ps := range ws.p.Game.LeaderBoard {
		pp := ws.p.Game.GetPlayerByPosition(Point{ps.X, ps.Y})
		if pp == nil {
			log.Warn().Msgf("player not found")
			continue
		}
		if ws.p.IsMe(pp) {
			continue
		}

		score := float64(pp.Score - minScore) / float64(maxScore - minScore)
		if score > tmpMaxNormalizedScore {
			tmpMaxNormalizedScore = score
		}

		playerDistance := distanceCalculator.Distance(ws.p.GetPosition(), pp.GetPosition())
		distance := 1 / ((playerDistance - minDistance) / (maxDistance - minDistance))
		if distance > tmpMaxNormalizedDistance {
			tmpMaxNormalizedDistance = distance
		}

		wp := weightedPlayer{
			Player:   pp,
			distance: distance,
			score:    score,
		}

		wps = append(wps, wp)
	}

	for idx, wp := range wps {
		wps[idx].distance = wp.distance / tmpMaxNormalizedDistance
		wps[idx].score = wp.score / tmpMaxNormalizedScore
	}

	sort.Sort(byWeight(wps))
	var players []Player

	for _, wp := range wps {
		players = append(players, *wp.Player)
	}
	return players
}

type weightedPlayer struct {
	*Player

	distance float64
	score    float64
}

func (w weightedPlayer) Weight() float64 {
	return w.distance + w.score
}

type byWeight []weightedPlayer

func (a byWeight) Len() int { return len(a) }
func (a byWeight) Less(i, j int) bool {
	return a[i].Weight() > a[j].Weight()
}
func (a byWeight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

