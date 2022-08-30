package main

import "sort"

type LeaderBoard []PlayerState

func (l LeaderBoard) GetPlayerByRank(rank int) *PlayerState {
	if rank > len(l) || rank < 0 {
		return nil
	}
	return &l[rank]
}

func (l LeaderBoard) GetRank(p Player) int {
	for rank, ps := range l {
		if ps.URL == p.Name {
			return rank
		}
	}
	return -1
}

func (l LeaderBoard) Sort() {
	sort.Sort(byScore(l))
}

// ByAge implements sort.Interface based on the Age field.
type byScore []PlayerState

func (a byScore) Len() int           { return len(a) }
func (a byScore) Less(i, j int) bool {
	if a[i].Score != a[j].Score {
		return a[i].Score > a[j].Score
	}
	return a[i].URL < a[j].URL
}
func (a byScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

