package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayer_FindClosestPlayers(t *testing.T) {
	type fields struct {
		X         int
		Y         int
		Direction string
		Game      Game
	}
	tests := []struct {
		name   string
		fields fields
		want   []Player
	}{
		{
			name: "get closest player",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  7,
						Height: 5,
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {}, {}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {Player: &PlayerState{X: 5, Y: 3, Direction: "E"}}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{X: 0, Y: 0, Direction: "E"},
						{X: 1, Y: 1, Direction: "E"},
						{X: 5, Y: 3, Direction: "E"},
					},
				},
			},
			want: []Player{
				{
					X:         0,
					Y:         0,
					Direction: "E",
				},
				{
					X:         5,
					Y:         3,
					Direction: "E",
				},
			},
		},
		{
			name: "get closest player",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
				Game: Game{
					Arena: Arena{
						Width:  7,
						Height: 5,
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {}, {}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 2, Direction: "E"}}, {}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{X: 0, Y: 2, Direction: "E"},
						{X: 0, Y: 0, Direction: "E"},
						{X: 1, Y: 1, Direction: "W"},
					},
				},
			},
			want: []Player{
				{
					X:         0,
					Y:         0,
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
			name: "get closest player",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "N",
				Game: Game{
					Arena: Arena{
						Width:  7,
						Height: 5,
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {Player: &PlayerState{X: 2, Y: 0, Direction: "E"}}, {}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "N"}}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{X: 2, Y: 0, Direction: "E"},
						{X: 0, Y: 0, Direction: "E"},
						{X: 1, Y: 1, Direction: "W"},
					},
				},
			},
			want: []Player{
				{
					X:         0,
					Y:         0,
					Direction: "E",
				},
				{
					X:         2,
					Y:         0,
					Direction: "E",
				},
			},
		},
		{
			name: "get closest player",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  7,
						Height: 5,
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {Player: &PlayerState{X: 2, Y: 0, Direction: "E"}}, {}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{X: 2, Y: 0, Direction: "E"},
						{X: 0, Y: 0, Direction: "E"},
						{X: 1, Y: 1, Direction: "W"},
					},
				},
			},
			want: []Player{
				{
					X:         2,
					Y:         0,
					Direction: "E",
				},
				{
					X:         0,
					Y:         0,
					Direction: "E",
				},
			},
		},
		{
			name: "no closest player found",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  7,
						Height: 5,
						Grid: [][]Cell{
							{{}, {}, {}, {}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
							{{}, {}, {}, {}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{X: 1, Y: 1, Direction: "E"},
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{
				X:         tt.fields.X,
				Y:         tt.fields.Y,
				Direction: tt.fields.Direction,
				Game:      tt.fields.Game,
			}
			ss := StepSorter{p: p}
			got := ss.Sort(context.TODO())
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
