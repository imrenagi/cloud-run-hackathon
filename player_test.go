package main

import (
	"context"
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
			if got := p.Walk(context.TODO()); got != tt.want {
				t.Errorf("Walk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_GetPlayersInRange(t *testing.T) {
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
			if got := p.GetPlayersInRange(context.TODO(), tt.args.direction, 3); len(got) != tt.want {
				t.Errorf("GetPlayersInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_FindShooterOnDirection(t *testing.T) {
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
		want   *Player
	}{
		{
			name: "no shooter",
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
			args: args{
				direction: North,
			},
			want: nil,
		},
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
			want: &Player{
				X:         1,
				Y:         0,
				Direction: "S",
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
			want: &Player{
				X:         3,
				Y:         1,
				Direction: "W",
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
			want: &Player{
				X:         1,
				Y:         2,
				Direction: "N",
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
			want: &Player{
				X:         0,
				Y:         1,
				Direction: "E",
			},
		},
		{
			name: "multiple shooters facing from front",
			fields: fields{
				X:         2,
				Y:         1,
				Direction: "W",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 2, Y: 1, Direction: "W"}}, {}},
							{{}, {}, {}, {}},
						},
					},
				},
			},
			args: args{

				direction: West,
			},
			want: &Player{
				X:         1,
				Y:         1,
				Direction: "E",
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
			got := p.FindShooterOnDirection(context.TODO(), tt.args.direction)
			if tt.want != nil {
				assert.Equal(t, tt.want.Direction, got.Direction)
				assert.Equal(t, tt.want.X, got.X)
				assert.Equal(t, tt.want.Y, got.Y)
			} else {
				assert.Nil(t, got)
			}

		})
	}
}

func TestPlayer_MoveToAdjacent(t *testing.T) {
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
			got := p.RequiredMoves(context.TODO(), tt.args.forPath, tt.args.opts...)
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
			got := p.FindClosestPlayers(context.TODO())
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

func TestPlayer_CanHitPoint(t *testing.T) {
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
		opts []HitOption
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		skip   bool
		want   bool
	}{
		{
			name: "can attack",
			args: args{
				pt: Point{3, 1},
			},
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 2,
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {Player: &PlayerState{X: 3, Y: 1, Direction: "W"}}},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "can attack even when there is no player",
			args: args{pt: Point{3, 1}},
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 2,
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "cant attack because there is other player in attack range",
			args: args{pt: Point{3, 1}},
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 2,
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 2, Y: 1, Direction: "N"}}, {Player: &PlayerState{X: 3, Y: 1, Direction: "W"}}},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "can attack even if there is other player in attack range because we enable the option to ignore players",
			args: args{
				opts: []HitOption{WithIgnorePlayer()},
				pt: Point{3, 1},
			},
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 2,
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 2, Y: 1, Direction: "N"}}, {}},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "cant attack because there is other player in attack range, even when target has no player",
			args: args{pt: Point{3, 1}},
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 2,
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 2, Y: 1, Direction: "N"}}, {}},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "cant attack because outside range",
			args: args{pt: Point{3, 1}},
			fields: fields{
				X:         1,
				Y:         1,
				Direction: "N",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 2,
						Grid: [][]Cell{
							{{}, {}, {}, {}},
							{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "N"}}, {}, {Player: &PlayerState{X: 3, Y: 1, Direction: "W"}}},
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.skip {
				t.Skip()
			}

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
			if got := p.CanHitPoint(context.TODO(), tt.args.pt, tt.args.opts...); got != tt.want {
				t.Errorf("CanHitPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_GetRank(t *testing.T) {

	type fields struct {
		Name string
		X    int
		Y    int
		Game Game
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "should get my rank",
			fields: fields{
				Name: "http://testing1",
				Game: Game{
					LeaderBoard: []PlayerState{
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
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{
				Name: tt.fields.Name,
				X:    tt.fields.X,
				Y:    tt.fields.Y,
				Game: tt.fields.Game,
			}
			if got := p.GetRank(); got != tt.want {
				t.Errorf("GetRank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_GetPlayerOnNextPodium(t *testing.T) {
	type fields struct {
		Name      string
		X         int
		Y         int
		Direction string
		Score     int
		Game      Game
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Player
	}{
		{
			name: "if alone return nil",
			fields: fields{
				X: 0, Y: 0, Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  2,
						Height: 1,
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{
							URL:       "",
							X:         0,
							Y:         0,
							Direction: "E",
							Score:     0,
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: nil,
		},
		{
			name: "if 3rd place, return second place",
			fields: fields{
				Name: "3",
				X:    0, Y: 0, Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 2, Direction: "E"}}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{
							URL:       "1",
							X:         0,
							Y:         1,
							Direction: "E",
							Score:     5,
						},
						{
							URL:       "2",
							X:         0,
							Y:         2,
							Direction: "E",
							Score:     4,
						},
						{
							URL:       "3",
							X:         0,
							Y:         0,
							Direction: "E",
							Score:     3,
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: &Player{
				X:         0,
				Y:         2,
				Direction: "E",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{
				Name:      tt.fields.Name,
				X:         tt.fields.X,
				Y:         tt.fields.Y,
				Direction: tt.fields.Direction,
				Score:     tt.fields.Score,
				Game:      tt.fields.Game,
			}

			got := p.GetPlayerOnNextPodium(tt.args.ctx)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want.X, got.X)
				assert.Equal(t, tt.want.Y, got.Y)
				assert.Equal(t, tt.want.Direction, got.Direction)
			}

		})
	}
}

func TestPlayer_GetHighestRank(t *testing.T) {
	type fields struct {
		Name         string
		X            int
		Y            int
		Direction    string
		WasHit       bool
		Score        int
		Game         Game
		State        State
		Strategy     Strategy
		trappedCount int
		Whitelisted  map[string]string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Player
	}{
		{
			name: "return highest rank",
			fields: fields{
				Name: "3",
				X:    0, Y: 0, Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 2, Direction: "E"}}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{
							URL:       "1",
							X:         0,
							Y:         1,
							Direction: "E",
							Score:     5,
						},
						{
							URL:       "2",
							X:         0,
							Y:         2,
							Direction: "E",
							Score:     4,
						},
						{
							URL:       "3",
							X:         0,
							Y:         0,
							Direction: "E",
							Score:     3,
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: &Player{
				X:         0,
				Y:         1,
				Direction: "E",
			},
		},
		{
			name: "return highest rank, skip whitelisted player",
			fields: fields{
				Whitelisted: map[string]string{
					"1": "1",
				},
				Name: "3",
				X:    0, Y: 0, Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
						Grid: [][]Cell{
							{{Player: &PlayerState{URL: "3", X: 0, Y: 0, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{URL: "1", X: 0, Y: 1, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{URL: "2", X: 0, Y: 2, Direction: "E"}}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{
							URL:       "1",
							X:         0,
							Y:         1,
							Direction: "E",
							Score:     5,
						},
						{
							URL:       "2",
							X:         0,
							Y:         2,
							Direction: "E",
							Score:     4,
						},
						{
							URL:       "3",
							X:         0,
							Y:         0,
							Direction: "E",
							Score:     3,
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: &Player{
				X:         0,
				Y:         2,
				Direction: "E",
			},
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
				State:        tt.fields.State,
				Strategy:     tt.fields.Strategy,
				trappedCount: tt.fields.trappedCount,
				Whitelisted:  tt.fields.Whitelisted,
			}

			got := p.GetHighestRank(tt.args.ctx)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want.X, got.X)
				assert.Equal(t, tt.want.Y, got.Y)
				assert.Equal(t, tt.want.Direction, got.Direction)
			}
		})
	}
}

func TestPlayer_GetLowestRank(t *testing.T) {
	type fields struct {
		Name         string
		X            int
		Y            int
		Direction    string
		WasHit       bool
		Score        int
		Game         Game
		State        State
		Strategy     Strategy
		trappedCount int
		Whitelisted  map[string]string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Player
	}{
		{
			name: "return lowest rank",
			fields: fields{
				Name: "3",
				X:    0, Y: 0, Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
						Grid: [][]Cell{
							{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{X: 0, Y: 2, Direction: "E"}}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{
							URL:       "1",
							X:         0,
							Y:         1,
							Direction: "E",
							Score:     5,
						},
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
					},
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: &Player{
				X:         0,
				Y:         2,
				Direction: "E",
			},
		},
		{
			name: "return lowest rank, skip whitelisted player",
			fields: fields{
				Whitelisted: map[string]string{
					"4": "4",
				},
				Name: "3",
				X:    0, Y: 0, Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
						Grid: [][]Cell{
							{{Player: &PlayerState{URL: "3", X: 0, Y: 0, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{URL: "1", X: 0, Y: 1, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{URL: "2", X: 0, Y: 2, Direction: "E"}}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{
							URL:       "1",
							X:         0,
							Y:         1,
							Direction: "E",
							Score:     5,
						},
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
							URL:       "4",
							X:         2,
							Y:         2,
							Direction: "E",
							Score:     1,
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: &Player{
				X:         0,
				Y:         2,
				Direction: "E",
			},
		},
		{
			name: "return lowest rank, skip if it is me player",
			fields: fields{
				Name: "3",
				X:    0, Y: 0, Direction: "E",
				Game: Game{
					Arena: Arena{
						Width:  4,
						Height: 3,
						Grid: [][]Cell{
							{{Player: &PlayerState{URL: "3", X: 0, Y: 0, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{URL: "1", X: 0, Y: 1, Direction: "E"}}, {}, {}, {}},
							{{Player: &PlayerState{URL: "2", X: 0, Y: 2, Direction: "E"}}, {}, {}, {}},
						},
					},
					LeaderBoard: []PlayerState{
						{
							URL:       "1",
							X:         0,
							Y:         1,
							Direction: "E",
							Score:     5,
						},
						{
							URL:       "2",
							X:         0,
							Y:         2,
							Direction: "E",
							Score:     3,
						},
						{
							URL:       "3",
							X:         0,
							Y:         0,
							Direction: "E",
							Score:     2,
						},

					},
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: &Player{
				X:         0,
				Y:         2,
				Direction: "E",
			},
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
				State:        tt.fields.State,
				Strategy:     tt.fields.Strategy,
				trappedCount: tt.fields.trappedCount,
				Whitelisted:  tt.fields.Whitelisted,
			}

			got := p.GetLowestRank(tt.args.ctx)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want.X, got.X)
				assert.Equal(t, tt.want.Y, got.Y)
				assert.Equal(t, tt.want.Direction, got.Direction)
			}
		})
	}
}