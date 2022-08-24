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
		want   Move
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
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
							{{}, {}, {}, {}},
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
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "W"}}, {}, {}, {}},
							{{}, {}, {}, {}},
							{{}, {}, {}, {}},
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
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 1, Direction: "W"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
							{{}, {}, {}, {}},
						},
					},
				},
			},
			want: "R",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{
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
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0}}, {Player: &PlayerState{X: 1, Y: 0}}, {Player: &PlayerState{X: 2, Y: 0}}, {Player: &PlayerState{X: 3, Y: 0}}},
							{{}, {}, {}, {}},
							{{}, {}, {}, {}},
							{{}, {}, {}, {}},
						},
					},
					// PlayersByPosition: map[string]Player{
					// 	"0,0": {
					// 		X: 0,
					// 		Y: 0,
					// 	},
					// 	"1,0": {
					// 		X: 1,
					// 		Y: 0,
					// 	},
					// 	"2,0": {
					// 		X: 2,
					// 		Y: 0,
					// 	},
					// 	"3,0": {
					// 		X: 3,
					// 		Y: 0,
					// 	},
					// },
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
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 0, Direction: "E"}}, {Player: &PlayerState{X: 2, Y: 0, Direction: "E"}}, {Player: &PlayerState{X: 3, Y: 0, Direction: "E"}}},
							{{}, {}, {}, {}},
							{{}, {}, {}, {}},
							{{}, {}, {}, {}},
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
			p := Player{
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
		want   []Player
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
						Grid: [][]Cell{
							{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
							{{}, {}, {}, {}},
						},
					},
				},
			},
			args: args{
				direction: North,
			},
			want: []Player{
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
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {Player: &PlayerState{X: 3, Y: 1, Direction: "W"}}},
							{{}, {}, {}, {}},
						},
					},
				},
			},
			args: args{
				direction: East,
			},
			want: []Player{
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
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {}, {}},
							{{}, {}, {}, {}},
						},
					},
				},
			},
			args: args{
				direction: South,
			},
			want: []Player{
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
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
							{{}, {}, {}, {}},
						},
					},
				},
			},
			args: args{

				direction: West,
			},
			want: []Player{
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
		want    []Move
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
			want: []Move{TurnLeft, WalkForward},
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
			want: []Move{TurnLeft, TurnLeft, WalkForward},
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
			want: []Move{TurnRight, WalkForward},
		},
		{
			name: "heading west, target is in west",
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "W",
			},
			args: args{
				toPt: Point{0, 1},
			},
			want: []Move{WalkForward},
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
			p := Player{
				X:         tt.fields.X,
				Y:         tt.fields.Y,
				Direction: tt.fields.Direction,
				WasHit:    tt.fields.WasHit,
				Score:     tt.fields.Score,
				Game:      tt.fields.Game,
			}
			got, err := p.MoveToAdjacent(tt.args.toPt)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPlayer_NextMove(t *testing.T) {
	type fields struct {
		X         int
		Y         int
		Direction string
		Game      Game
	}
	type args struct {
		forPath Path
		opts    []MoveOption
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Move
	}{
		{
			name: "not moving",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "E",
				Game: Game{
					Arena: NewArena(5, 3),
				},
			},
			args: args{
				forPath: Path{{0, 0}},
			},
			want: nil,
		},
		{
			name: "move forward",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "E",
				Game: Game{
					Arena: NewArena(5, 3),
				},
			},
			args: args{
				forPath: Path{{0, 0}, {1, 0}},
			},
			want: []Move{WalkForward},
		},
		{
			name: "move forward and turn",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "E",
				Game: Game{
					Arena: NewArena(5, 3),
				},
			},
			args: args{
				forPath: Path{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
			},
			want: []Move{WalkForward, WalkForward, TurnRight, WalkForward},
		},
		{
			name: "only return next move",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "E",
				Game: Game{
					Arena: NewArena(5, 3),
				},
			},
			args: args{
				forPath: Path{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
				opts:    []MoveOption{WithOnlyNextMove()},
			},
			want: []Move{WalkForward},
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
			got := p.RequiredMoves(tt.args.forPath, tt.args.opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequiredMoves() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_Apply(t *testing.T) {
	type fields struct {
		X         int
		Y         int
		Direction string
		Game      Game
	}
	type args struct {
		m Move
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "turn left",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "N",
			},
			args: args{
				m: TurnLeft,
			},
			want: fields{
				X:         0,
				Y:         0,
				Direction: "W",
			},
		},
		{
			name: "turn right",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "N",
			},
			args: args{
				m: TurnRight,
			},
			want: fields{
				X:         0,
				Y:         0,
				Direction: "E",
			},
		},
		{
			name: "move forward",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "E",
				Game: Game{
					Arena: NewArena(5, 3),
				},
			},
			args: args{
				m: WalkForward,
			},
			want: fields{
				X:         1,
				Y:         0,
				Direction: "E",
			},
		},
		{
			name: "move forward",
			fields: fields{
				X:         0,
				Y:         0,
				Direction: "N",
				Game: Game{
					Arena: NewArena(5, 3),
				},
			},
			args: args{
				m: WalkForward,
			},
			want: fields{
				X:         0,
				Y:         0,
				Direction: "N",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				X:         tt.fields.X,
				Y:         tt.fields.Y,
				Direction: tt.fields.Direction,
				Game:      tt.fields.Game,
			}
			p.Apply(tt.args.m)
			assert.Equal(t, tt.want.X, p.X)
			assert.Equal(t, tt.want.Y, p.Y)
			assert.Equal(t, tt.want.Direction, p.Direction)
		})
	}
}

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
		want   *Player
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
					Players: []PlayerState{
						{X: 0, Y: 0, Direction: "E"},
						{X: 1, Y: 1, Direction: "E"},
						{X: 5, Y: 3, Direction: "E"},
					},
				},
			},
			want: &Player{
				X:         0,
				Y:         0,
				Direction: "E",
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
					Players: []PlayerState{
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
			got := p.FindClosestPlayers()
			if tt.want != nil {
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.X, got.X)
				assert.Equal(t, tt.want.Y, got.Y)
				assert.Equal(t, tt.want.Direction, got.Direction)
			} else {
				assert.Nil(t, got)
			}

		})
	}
}

func TestPlayer_CanAttack(t *testing.T) {
	type fields struct {
		Name         string
		X            int
		Y            int
		Direction    string
		WasHit       bool
		Score        int
		Game         Game
		Strategy     Strategy
		trappedCount int
	}
	type args struct {
		pt Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "can attack",
			args: args{pt: Point{2, 1}},
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  7,
						Height: 5,
					},
				},
			},
			want: true,
		},
		{
			name: "cant attack",
			args: args{pt: Point{1, 0}},
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  7,
						Height: 5,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{
				Name:         tt.fields.Name,
				X:            tt.fields.X,
				Y:            tt.fields.Y,
				Direction:    tt.fields.Direction,
				WasHit:       tt.fields.WasHit,
				Score:        tt.fields.Score,
				Game:         tt.fields.Game,
				Strategy:     tt.fields.Strategy,
				trappedCount: tt.fields.trappedCount,
			}
			if got := p.CanAttack(tt.args.pt); got != tt.want {
				t.Errorf("CanAttack() = %v, want %v", got, tt.want)
			}
		})
	}
}
