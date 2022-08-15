package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerState_GetPlayersInFront(t *testing.T) {
	type fields struct {
		X         int
		Y         int
		Direction string
		WasHit    bool
		Score     int
	}
	type args struct {
		g Game
		direction Direction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "get all player in front east direction",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "E",
				WasHit:    false,
				Score:     0,
			},
			args: args{
				direction: East,
				g: Game{
					Dimension: []int{4, 4},
					PlayersByPosition: map[string]PlayerState{
						"0,0": {
							X: 0,
							Y: 0,
						},
						"1,0": {
							X: 1,
							Y: 0,
						},
						"2,0": {
							X: 2,
							Y: 0,
						},
						"3,0": {
							X: 3,
							Y: 0,
						},
					},
				},
			},
			want: 3,
		},
		{
			name: "get all player in front east direction",
			fields: fields{
				X: 2,
				Y: 0,
				Direction: "E",
				WasHit:    false,
				Score:     0,
			},
			args: args{
				direction: East,
				g: Game{
					Dimension: []int{4, 4},
					PlayersByPosition: map[string]PlayerState{
						"0,0": {
							X: 0,
							Y: 0,
						},
						"1,0": {
							X: 1,
							Y: 0,
						},
						"2,0": {
							X: 2,
							Y: 0,
						},
						"3,0": {
							X: 3,
							Y: 0,
						},
					},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PlayerState{
				X:         tt.fields.X,
				Y:         tt.fields.Y,
				Direction: tt.fields.Direction,
				WasHit:    tt.fields.WasHit,
				Score:     tt.fields.Score,
			}
			if got := p.GetPlayersInDirection(tt.args.g, tt.args.direction); got != tt.want {
				t.Errorf("GetPlayersInDirection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPointAfterMove(t *testing.T) {
	type args struct {
		p         Point
		direction Direction
		distance  int
	}
	tests := []struct {
		name string
		args args
		want Point
	}{
		{
			name: "",
			args: args{
				p:         Point{0, 0},
				direction: South,
				distance:  1,
			},
			want: Point{0, 1},
		},
		{
			name: "",
			args: args{
				p:         Point{0, 0},
				direction: East,
				distance:  1,
			},
			want: Point{1, 0},
		},
		{
			name: "",
			args: args{
				p:         Point{1, 1},
				direction: North,
				distance:  3,
			},
			want: Point{1, -2},
		},
		{
			name: "",
			args: args{
				p:         Point{1, 1},
				direction: West,
				distance:  3,
			},
			want: Point{-2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := NewPointAfterMove(tt.args.p, tt.args.direction, tt.args.distance);
			// assert.Equal(t, tt.want, got)
			assert.InDelta(t, tt.want.X, got.X, 0.01)
			assert.InDelta(t, tt.want.Y, got.Y, 0.01)
		})
	}
}
