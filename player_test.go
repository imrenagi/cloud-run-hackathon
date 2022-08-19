package main

import (
	"reflect"
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

func TestPlayerState_GoTo(t *testing.T) {

	t.Fail()
	t.Log("not implemented")

	type fields struct {
		X         int
		Y         int
		Direction string
		WasHit    bool
		Score     int
		Game      Game
	}
	type args struct {
		pt        Point
		direction Direction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Decision
	}{
		// TODO: Add test cases.
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
			if got := p.GoTo(tt.args.pt, tt.args.direction); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoTo() = %v, want %v", got, tt.want)
			}
		})
	}
}