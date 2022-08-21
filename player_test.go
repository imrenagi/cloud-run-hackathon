package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerState_Walk(t *testing.T) {
	type fields struct {
		X         int
		Y         int
		Direction string
		WasHit    bool
		Score     int
		Game      Game
	}

	tests := []struct {
		name   string
		fields fields
		want   Decision
	}{
		{
			name: "move forward",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
					},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "W",
						},
					},
				},
			},
			want: "F",
		},
		{
			name: "found edge, turn right",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "W",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
					},
					PlayersByPosition: map[string]PlayerState{
						"0,0": {
							X:         0,
							Y:         0,
							Direction: "W",
						},
					},
				},
			},
			want: "R",
		},
		{
			name: "found enemy in front, turn right",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
					},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "W",
						},
						"0,1": {
							X:         0,
							Y:         1,
							Direction: "W",
						},
					},
				},
			},
			want: "R",
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
				Game:      tt.fields.Game,
			}
			if got := p.Walk(); got != tt.want {
				t.Errorf("Walk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayerState_GetPlayersInFront(t *testing.T) {
	type fields struct {
		X         int
		Y         int
		Direction string
		WasHit    bool
		Score     int
		Game      Game
	}
	type args struct {
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
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 4,
					},
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
			args: args{
				direction: East,
			},
			want: 3,
		},
		{
			name: "get all player in front east direction",
			fields: fields{
				X:         2,
				Y:         0,
				Direction: "E",
				WasHit:    false,
				Score:     0,
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 4,
					},
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
			args: args{
				direction: East,
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
				Game:      tt.fields.Game,
			}
			if got := p.GetPlayersInRange(tt.args.direction, 3); len(got) != tt.want {
				t.Errorf("GetPlayersInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayerState_FindShooterFromDirection(t *testing.T) {
	type fields struct {
		X         int
		Y         int
		Direction string
		WasHit    bool
		Score     int
		Game      Game
	}
	type args struct {
		direction Direction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []PlayerState
	}{
		{
			name: "shooter from right",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
					},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "W",
						},
						"1,0": {
							X:         1,
							Y:         0,
							Direction: "S",
						},
					},
				},
			},
			args: args{
				direction: North,
			},
			want: []PlayerState{
				{
					X:         1,
					Y:         0,
					Direction: "S",
				},
			},
		},
		{
			name: "shooter from back",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
					},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "W",
						},
						"3,1": {
							X:         3,
							Y:         1,
							Direction: "W",
						},
					},
				},
			},
			args: args{
				direction: East,
			},
			want: []PlayerState{
				{
					X:         3,
					Y:         1,
					Direction: "W",
				},
			},
		},
		{
			name: "shooter from left",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 4,
					},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "W",
						},
						"1,2": {
							X:         1,
							Y:         2,
							Direction: "N",
						},
					},
				},
			},
			args: args{
				direction: South,
			},
			want: []PlayerState{
				{
					X:         1,
					Y:         2,
					Direction: "N",
				},
			},
		},
		{
			name: "shooter from front",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
					},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "W",
						},
						"0,1": {
							X:         0,
							Y:         1,
							Direction: "E",
						},
					},
				},
			},
			args: args{

				direction: West,
			},
			want: []PlayerState{
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
			p := PlayerState{
				X:         tt.fields.X,
				Y:         tt.fields.Y,
				Direction: tt.fields.Direction,
				WasHit:    tt.fields.WasHit,
				Score:     tt.fields.Score,
				Game:      tt.fields.Game,
			}
			got := p.FindShooterOnDirection(tt.args.direction)
			assert.Equal(t, len(tt.want), len(got))
			for i, p := range tt.want {
				assert.Equal(t, p.X, tt.want[i].X)
				assert.Equal(t, p.Y, tt.want[i].Y)
			}
		})
	}
}


func TestPlayerState_GetShortestRotation(t *testing.T) {
	type fields struct {
		X         int
		Y         int
		Direction string
		WasHit    bool
		Score     int
		Game      Game
	}
	type args struct {
		toPt Point
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Decision
		wantErr error
	}{
		{
			name: "heading west, target is in south",
			fields: fields{
				X:         0,
				Y:         1,
				Direction: "W",
			},
			args: args{
				toPt: Point{0, 2},
			},
			want:    []Decision{TurnLeft},
		},
		{
			name: "heading west, target is in east",
			fields: fields{
				X:         0,
				Y:         1,
				Direction: "W",
			},
			args: args{
				toPt: Point{1, 1},
			},
			want:    []Decision{TurnLeft, TurnLeft},
		},
		{
			name: "heading west, target is in north",
			fields: fields{
				X:         0,
				Y:         1,
				Direction: "W",
			},
			args: args{
				toPt: Point{0, 0},
			},
			want:    []Decision{TurnRight},
		},
		{
			name: "heading west, target is in north",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "W",
			},
			args: args{
				toPt: Point{2, 0},
			},
			want:    nil,
			wantErr: ErrDestNotFound,
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
				Game:      tt.fields.Game,
			}
			got, err := p.GetShortestRotation(tt.args.toPt)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVector_Angle(t *testing.T) {

	t.Skip()
	t.Log("still not complete")

	type fields struct {
		X float64
		Y float64
	}
	type args struct {
		v2 Vector
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name:   "90 degree",
			fields: fields{X: 1, Y: 0},
			args:   args{
				v2: Vector{X: 0, Y: 1},
			},
			want:   90,
		},
		{
			name:   "180 degree",
			fields: fields{X: 1, Y: 0},
			args:   args{
				v2: Vector{X: -1, Y: 0},
			},
			want:   180,
		},
		{
			name:   "-90 degree",
			fields: fields{X: 1, Y: 0},
			args:   args{
				v2: Vector{X: 0, Y: -1},
			},
			want:   -90,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if got := v.Angle(tt.args.v2); got != tt.want {
				t.Errorf("Angle() = %v, want %v", got, tt.want)
			}
		})
	}
}