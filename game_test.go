package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_UpdateArena(t *testing.T) {
	type args struct {
		a ArenaUpdate
		m Mode
	}
	tests := []struct {
		name string
		args args
		want Game
	}{
		{
			name: "update game",
			args: args{
				m: AggressiveMode,
				a: ArenaUpdate{
					Links: struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
					}{
						Self: struct {
							Href string `json:"href"`
						}{
							Href: "http://testing.run",
						},
					},
					Arena: struct {
						Dimensions []int                  `json:"dims"`
						State      map[string]PlayerState `json:"state"`
					}{
						Dimensions: []int{5, 4},
						State: map[string]PlayerState{
							"http://testing.run": {
								X:         1,
								Y:         1,
								Direction: "W",
								WasHit:    false,
								Score:     4,
							},
							"http://testing-02.run": {
								X:         2,
								Y:         1,
								Direction: "W",
								WasHit:    false,
								Score:     1,
							},
							"http://testing-03.run": {
								X:         2,
								Y:         2,
								Direction: "W",
								WasHit:    false,
								Score:     10,
							},
						},
					},
				},
			},
			want: Game{
				Mode: AggressiveMode,
				Arena: Arena{
					Width:  5,
					Height: 4,
					Grid: [][]Cell{
						{{}, {}, {}, {}, {}},
						{{}, {Player: &PlayerState{URL: "http://testing.run", X: 1, Y: 1, Direction: "W", Score: 4}}, {Player: &PlayerState{URL: "http://testing-02.run", X: 2, Y: 1, Direction: "W", Score: 1}}, {}, {}},
						{{}, {}, {Player: &PlayerState{URL: "http://testing-03.run", X: 2, Y: 2, Direction: "W", Score: 10}}, {}, {}},
						{{}, {}, {}, {}, {}},
					},
				},
				PlayerStateByURL: map[string]PlayerState{
					"http://testing.run": {
						X:         1,
						Y:         1,
						Direction: "W",
						WasHit:    false,
						Score:     4,
					},
					"http://testing-02.run": {
						X:         2,
						Y:         1,
						Direction: "W",
						WasHit:    false,
						Score:     1,
					},
					"http://testing-03.run": {
						X:         2,
						Y:         2,
						Direction: "W",
						WasHit:    false,
						Score:     10,
					},
				},
				LeaderBoard: []PlayerState{
					{
						URL:       "http://testing-03.run",
						X:         2,
						Y:         2,
						Direction: "W",
						WasHit:    false,
						Score:     10,
					}, {
						URL:       "http://testing.run",
						X:         1,
						Y:         1,
						Direction: "W",
						WasHit:    false,
						Score:     4,
					}, {
						URL:       "http://testing-02.run",
						X:         2,
						Y:         1,
						Direction: "W",
						WasHit:    false,
						Score:     1,
					},
				},
			},
		},
		{
			name: "update game without updating the leaderboard",
			args: args{
				m: NormalMode,
				a: ArenaUpdate{
					Links: struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
					}{
						Self: struct {
							Href string `json:"href"`
						}{
							Href: "http://testing.run",
						},
					},
					Arena: struct {
						Dimensions []int                  `json:"dims"`
						State      map[string]PlayerState `json:"state"`
					}{
						Dimensions: []int{5, 4},
						State: map[string]PlayerState{
							"http://testing.run": {
								X:         1,
								Y:         1,
								Direction: "W",
								WasHit:    false,
								Score:     4,
							},
							"http://testing-02.run": {
								X:         2,
								Y:         1,
								Direction: "W",
								WasHit:    false,
								Score:     1,
							},
							"http://testing-03.run": {
								X:         2,
								Y:         2,
								Direction: "W",
								WasHit:    false,
								Score:     10,
							},
						},
					},
				},
			},
			want: Game{
				Mode: NormalMode,
				Arena: Arena{
					Width:  5,
					Height: 4,
					Grid: [][]Cell{
						{{}, {}, {}, {}, {}},
						{{}, {Player: &PlayerState{URL: "http://testing.run", X: 1, Y: 1, Direction: "W", Score: 4}}, {Player: &PlayerState{URL: "http://testing-02.run", X: 2, Y: 1, Direction: "W", Score: 1}}, {}, {}},
						{{}, {}, {Player: &PlayerState{URL: "http://testing-03.run", X: 2, Y: 2, Direction: "W", Score: 10}}, {}, {}},
						{{}, {}, {}, {}, {}},
					},
				},
				PlayerStateByURL: map[string]PlayerState{
					"http://testing.run": {
						X:         1,
						Y:         1,
						Direction: "W",
						WasHit:    false,
						Score:     4,
					},
					"http://testing-02.run": {
						X:         2,
						Y:         1,
						Direction: "W",
						WasHit:    false,
						Score:     1,
					},
					"http://testing-03.run": {
						X:         2,
						Y:         2,
						Direction: "W",
						WasHit:    false,
						Score:     10,
					},
				},
				LeaderBoard: []PlayerState{
					{
						URL:       "http://testing.run",
						X:         1,
						Y:         1,
						Direction: "W",
						WasHit:    false,
						Score:     4,
					},
					{
						URL:       "http://testing-02.run",
						X:         2,
						Y:         1,
						Direction: "W",
						WasHit:    false,
						Score:     1,
					},
					{
						URL:       "http://testing-03.run",
						X:         2,
						Y:         2,
						Direction: "W",
						WasHit:    false,
						Score:     10,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame(WithGameMode(tt.args.m))
			g.UpdateArena(context.TODO(), tt.args.a)
			assert.Equal(t, tt.want, g)
		})
	}
}

func TestGame_GetPlayerByRank(t *testing.T) {
	type fields struct {
		Arena            Arena
		PlayerStateByURL map[string]PlayerState
		Players          LeaderBoard
	}
	type args struct {
		rank int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Player
	}{
		{
			name: "should get correct player",
			fields: fields{
				Arena: Arena{
					Width:  5,
					Height: 3,
					Grid: [][]Cell{
						{{}, {}, {}, {}, {}},
						{{}, {}, {}, {}, {}},
						{{}, {Player: &PlayerState{
							URL:       "http://testing2",
							X:         1,
							Y:         2,
							Direction: "E",
							WasHit:    false,
							Score:     4,
						}}, {Player: &PlayerState{
							URL:       "http://testing3",
							X:         2,
							Y:         2,
							Direction: "E",
							WasHit:    false,
							Score:     4,
						}}, {Player: &PlayerState{
							URL:       "http://testing1",
							X:         3,
							Y:         2,
							Direction: "E",
							WasHit:    false,
							Score:     2,
						}}, {}},
					},
				},
				Players: []PlayerState{
					{
						URL:       "http://testing2",
						X:         1,
						Y:         2,
						Direction: "E",
						WasHit:    false,
						Score:     4,
					},
					{
						URL:       "http://testing3",
						X:         2,
						Y:         2,
						Direction: "E",
						WasHit:    false,
						Score:     4,
					},
					{
						URL:       "http://testing1",
						X:         3,
						Y:         2,
						Direction: "E",
						WasHit:    false,
						Score:     2,
					},
				},
			},
			args: args{
				rank: 1,
			},
			want: &Player{
				Name:      "http://testing3",
				X:         2,
				Y:         2,
				Direction: "E",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Game{
				Arena:            tt.fields.Arena,
				PlayerStateByURL: tt.fields.PlayerStateByURL,
				LeaderBoard:      tt.fields.Players,
			}
			got := g.GetPlayerByRank(tt.args.rank)
			if tt.want != nil {
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.X, got.X)
				assert.Equal(t, tt.want.Y, got.Y)
				assert.Equal(t, tt.want.Direction, got.Direction)
			}
		})
	}
}

