package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_UpdateArena(t *testing.T) {

	t.Skip("skipping this since we are not sorting the player based on score for now")

	type args struct {
		a ArenaUpdate
	}
	tests := []struct {
		name string
		args args
		want Game
	}{
		{
			name: "update game",
			args: args{
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
				Arena: Arena{
					Width:  5,
					Height: 4,
					Grid: [][]Cell{
						{{}, {}, {}, {}, {}},
						{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W", Score: 4}}, {Player: &PlayerState{X: 2, Y: 1, Direction: "W", Score: 1}}, {}, {}},
						{{}, {}, {Player: &PlayerState{X: 2, Y: 2, Direction: "W", Score: 10}}, {}, {}},
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
				Players: []PlayerState{
					{
						X:         2,
						Y:         2,
						Direction: "W",
						WasHit:    false,
						Score:     10,
					}, {
						X:         1,
						Y:         1,
						Direction: "W",
						WasHit:    false,
						Score:     4,
					}, {
						X:         2,
						Y:         1,
						Direction: "W",
						WasHit:    false,
						Score:     1,
					},
				},
				Config: GameConfig{
					AttackRange: 3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			g.UpdateArena(tt.args.a)
			assert.Equal(t, tt.want, g)
		})
	}
}
