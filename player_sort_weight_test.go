package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayer_Weight(t *testing.T) {
	type fields struct {
		Name                string
		X                   int
		Y                   int
		Direction           string
		WasHit              bool
		Score               int
		Game                Game
		State               State
		Strategy            Strategy
		Whitelisted         map[string]string
		trappedCount        int
		consecutiveHitCount int
	}
	tests := []struct {
		name   string
		fields fields
		want   []Player
	}{
		{
			name: "higher score and closer should be weighted higher",
			fields: fields{
				Name: "3",
				X:    0, Y: 0, Direction: "E",
				Score: 3,
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "E", Score: 3}}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 1, Direction: "E", Score: 5}}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 2, Direction: "E", Score: 2}}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{
							URL:       "3",
							X:         0,
							Y:         0,
							Direction: "E",
							Score:     3,
						},
						{
							URL:       "2",
							X:         0,
							Y:         2,
							Direction: "E",
							Score:     2,
						},
						{
							URL:       "1",
							X:         0,
							Y:         1,
							Direction: "E",
							Score:     5,
						},
					},
				},
			},
			want: []Player{
				{
					X:         0,
					Y:         1,
					Direction: "E",
				},
				{
					X:         0,
					Y:         2,
					Direction: "E",
				},
			},
		},
		{
			name: "higher score but bigger distance, can have higher priority",
			fields: fields{
				Name: "3",
				X:    0, Y: 2, Direction: "E",
				Score: 2,
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 1, Direction: "E", Score: 3}}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 2, Direction: "E", Score: 2}}, {}, {Player: &PlayerState{X: 2, Y: 2, Direction: "E", Score: 5}}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{
							URL:       "3",
							X:         0,
							Y:         1,
							Direction: "E",
							Score:     3,
						},
						{
							URL:       "2",
							X:         0,
							Y:         2,
							Direction: "E",
							Score:     2,
						},
						{
							URL:       "1",
							X:         2,
							Y:         2,
							Direction: "E",
							Score:     5,
						},
					},
				},
			},
			want: []Player{
				{
					X:         2,
					Y:         2,
					Direction: "E",
				},
				{
					X:         0,
					Y:         1,
					Direction: "E",
				},

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{
				Name:                tt.fields.Name,
				X:                   tt.fields.X,
				Y:                   tt.fields.Y,
				Direction:           tt.fields.Direction,
				WasHit:              tt.fields.WasHit,
				Score:               tt.fields.Score,
				Game:                tt.fields.Game,
				State:               tt.fields.State,
				Strategy:            tt.fields.Strategy,
				Whitelisted:         tt.fields.Whitelisted,
				trappedCount:        tt.fields.trappedCount,
				consecutiveHitCount: tt.fields.consecutiveHitCount,
			}
			ws := WeightedSorter{p: p}
			got := ws.Sort(context.TODO())
			assert.Equal(t, len(tt.want), len(got))
			if tt.want != nil {
				for idx, res := range got {
					assert.NotNil(t, got)
					assert.Equal(t, tt.want[idx].X, res.X)
					assert.Equal(t, tt.want[idx].Y, res.Y)
					assert.Equal(t, tt.want[idx].Direction, res.Direction)
				}
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
