package main

import (
	"reflect"
	"testing"
)

func TestPlayerState_MoveForward(t *testing.T) {
	type fields struct {
		X         int
		Y         int
		Direction string
		WasHit    bool
		Score     int
	}
	type args struct {
		g Game
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Decision
	}{
		{
			name:   "move forward",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
			},
			args:   args{
				g: Game{
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
			want:   "F",
		},
		{
			name:   "found edge, turn right",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "W",
			},
			args:   args{
				g: Game{
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
			want:   "R",
		},
		{
			name:   "found enemy in front, turn right",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
			},
			args:   args{
				g: Game{
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
			want:   "R",
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
			if got := p.Walk(tt.args.g); got != tt.want {
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
	}
	type args struct {
		g         Game
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
			},
			args: args{
				direction: East,
				g: Game{
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
			if got := p.GetPlayersInRange(tt.args.g, tt.args.direction, 3); len(got) != tt.want {
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
	}
	type args struct {
		g         Game
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
			},
			args: args{
				g: Game{
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
			},
			args: args{
				g: Game{
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
			},
			args: args{
				g: Game{
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
			},
			args: args{
				g: Game{
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
			}
			if got := p.FindShooterOnDirection(tt.args.g, tt.args.direction); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindShooterOnDirection() = %v, want %v", got, tt.want)
			}
		})
	}
}
