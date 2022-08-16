package main

import (
	"reflect"
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
			if got := p.GetPlayersInRange(tt.args.g, tt.args.direction, 3); len(got) != tt.want {
				t.Errorf("GetPlayersInRange() = %v, want %v", got, tt.want)
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

func TestPlayerState_Escape(t *testing.T) {
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
		// bug muter2
		{
			name:   "opponent is attacking from the top and bottom, player heading west",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,0": {
							X:         1,
							Y:         0,
							Direction: "S",
						},
						"1,2": {
							X:         1,
							Y:         2,
							Direction: "N",
						},
					},
				},
			},
			want:   "F",
		},
		{
			name:   "opponent is attacking from the top and bottom, player heading south",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "S",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,0": {
							X:         1,
							Y:         0,
							Direction: "S",
						},
						"1,2": {
							X:         1,
							Y:         2,
							Direction: "N",
						},
					},
				},
			},
			want:   "R",
		},
		{
			name:   "opponent is attacking from the top and bottom, player heading north",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "N",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,0": {
							X:         1,
							Y:         0,
							Direction: "S",
						},
						"1,2": {
							X:         1,
							Y:         2,
							Direction: "N",
						},
					},
				},
			},
			want:   "R",
		},
		{
			name:   "opponent is attacking from the top and bottom, player heading east",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,0": {
							X:         1,
							Y:         0,
							Direction: "S",
						},
						"1,2": {
							X:         1,
							Y:         2,
							Direction: "N",
						},
					},
				},
			},
			want:   "F",
		},
		{
			name:   "opponent is attacking from the right, player heading east",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "E",
						},
						"1,2": {
							X:         1,
							Y:         2,
							Direction: "N",
						},
					},
				},
			},
			want:   "F",
		},
		{
			name:   "opponent is attacking from the bottom, player heading west",
			fields: fields{
				X:         6,
				Y:         2,
				Direction: "W",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{7,5},
					PlayersByPosition: map[string]PlayerState{
						"6,2": {
							X:         6,
							Y:         2,
							Direction: "W",
						},
						"6,4": {
							X:         6,
							Y:         4,
							Direction: "N",
						},
					},
				},
			},
			want:   "F",
		},

		// -----
		{
			name:   "opponent is attacking from the front",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "E",
						},
						"2,1": {
							X:         2,
							Y:         1,
							Direction: "W",
						},
					},
				},
			},
			want:   "R",
		},
		{
			name:   "opponent is attacking from the left",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "E",
						},
						"1,0": {
							X:         1,
							Y:         0,
							Direction: "S",
						},
					},
				},
			},
			want:   "F",
		},
		{
			name:   "opponent is attacking from the back",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "E",
						},
						"0,1": {
							X:         0,
							Y:         1,
							Direction: "E",
						},
					},
				},
			},
			want:   "R",
		},
		{
			name:   "opponent is attacking from the back, user is on the edge",
			fields: fields{
				X:         0,
				Y:         1,
				Direction: "W",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "W",
						},
					},
				},
			},
			want:   "R",
		},
		{
			name:   "opponent is attacking from the back, user is not on the edge",
			fields: fields{
				X:         2,
				Y:         2,
				Direction: "W",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"3,2": {
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
			name:   "opponent is attacking from the back, but there is other opponent",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"0,1": {
							X:         0,
							Y:         1,
							Direction: "W",
						},
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "W",
						},
						"2,1": {
							X:         2,
							Y:         1,
							Direction: "W",
						},
					},
				},
			},
			want:   "R",
		},
		{
			name:   "opponent is attacking other user",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "E",
						},
						"2,1": {
							X:         2,
							Y:         1,
							Direction: "E",
						},
						"3,1": {
							X:         3,
							Y:         1,
							Direction: "W",
						},
					},
				},
			},
			want:   "R",
		},
		{
			name:   "none attacking",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				WasHit:    true,
				Score:     0,
			},
			args:   args{
				g: Game{
					Dimension:         []int{4,3},
					PlayersByPosition: map[string]PlayerState{
						"1,1": {
							X:         1,
							Y:         1,
							Direction: "E",
						},
					},
				},
			},
			want:   "F",
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
			if got := p.Escape(tt.args.g); got != tt.want {
				t.Errorf("Escape() = %v, want %v", got, tt.want)
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
				g:         Game{
					Dimension:         []int{4,3},
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
				g:         Game{
					Dimension:         []int{4,3},
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
				g:         Game{
					Dimension:         []int{4,4},
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
				g:         Game{
					Dimension:         []int{4,3},
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
			if got := p.FindShooterFromDirection(tt.args.g, tt.args.direction); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindShooterFromDirection() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
					Dimension:         []int{4,3},
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
					Dimension:         []int{4,3},
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
					Dimension:         []int{4,3},
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
			if got := p.MoveForward(tt.args.g); got != tt.want {
				t.Errorf("MoveForward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirection_Opposite(t *testing.T) {
	type fields struct {
		Name   string
		Degree int
	}
	tests := []struct {
		name   string
		fields fields
		want   Direction
	}{
		{
			name:   "south",
			fields: fields{
				Name:  South.Name,
				Degree: South.Degree,
			},
			want:   North,
		},
		{
			name:   "west",
			fields: fields{
				Name:  West.Name,
				Degree: West.Degree,
			},
			want:   East,
		},
		{
			name:   "north",
			fields: fields{
				Name:  North.Name,
				Degree: North.Degree,
			},
			want:   South,
		},
		{
			name:   "east",
			fields: fields{
				Name:  East.Name,
				Degree: East.Degree,
			},
			want:   West,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Direction{
				Name:   tt.fields.Name,
				Degree: tt.fields.Degree,
			}
			if got := d.Opposite(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Opposite() = %v, want %v", got, tt.want)
			}
		})
	}
}