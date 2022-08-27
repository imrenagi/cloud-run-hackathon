package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckTargetSurroundingAttackRangeFn(t *testing.T) {
	type args struct {
		target Player
	}
	type fnArgs struct {
		p Point
	}
	tests := []struct {
		name   string
		args   args
		fnArgs fnArgs
		want   bool
	}{
		{
			name: "target can attack a point",
			args: args{
				target: Player{
					X:         1,
					Y:         1,
					Direction: "E",
					Game: Game{
						LeaderBoard: []PlayerState{
							{X: 1, Y: 1, Direction: "E"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			fnArgs: fnArgs{p: Point{2, 1}},
			want: false,
		},
		{
			name: "there is other player in right-side of the target attacking the point",
			args: args{
				target: Player{
					X:         1,
					Y:         2,
					WasHit: true,
					Direction: "W",
					Game: Game{
						LeaderBoard: []PlayerState{
							{X: 1, Y: 0, Direction: "S"},
							{X: 1, Y: 2, Direction: "W"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "W"}}, {}, {}},
							},
						},
					},
				},
			},
			fnArgs: fnArgs{p: Point{1, 1}},
			want: false,
		},
		{
			name: "there is other player in left-side of the target attacking the point",
			args: args{
				target: Player{
					X:         1,
					Y:         2,
					WasHit: true,
					Direction: "E",
					Game: Game{
						LeaderBoard: []PlayerState{
							{X: 1, Y: 0, Direction: "S"},
							{X: 1, Y: 2, Direction: "E"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "E"}}, {}, {}},
							},
						},
					},
				},
			},
			fnArgs: fnArgs{p: Point{1, 1}},
			want: false,
		},
		{
			name: "there is other player in front-side of the target attacking the point",
			args: args{
				target: Player{
					X:         1,
					Y:         2,
					WasHit: true,
					Direction: "E",
					Game: Game{
						LeaderBoard: []PlayerState{
							{X: 3, Y: 2, Direction: "W"},
							{X: 1, Y: 2, Direction: "E"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "E"}}, {}, {Player: &PlayerState{X: 3, Y: 2, Direction: "W"}}},
							},
						},
					},
				},
			},
			fnArgs: fnArgs{p: Point{2, 2}},
			want: false,
		},
		{
			name: "there is other player in back-side of the target attacking the point",
			args: args{
				target: Player{
					X:         1,
					Y:         2,
					WasHit: true,
					Direction: "W",
					Game: Game{
						LeaderBoard: []PlayerState{
							{X: 3, Y: 2, Direction: "W"},
							{X: 1, Y: 2, Direction: "W"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "W"}}, {}, {Player: &PlayerState{X: 3, Y: 2, Direction: "W"}}},
							},
						},
					},
				},
			},
			fnArgs: fnArgs{p: Point{2, 2}},
			want: false,
		},
		{
			name: "valid should be unblock",
			args: args{
				target: Player{
					X:         1,
					Y:         2,
					Direction: "W",
					WasHit: true,
					Game: Game{
						LeaderBoard: []PlayerState{
							{X: 1, Y: 0, Direction: "S"},
							{X: 1, Y: 2, Direction: "W"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "W"}}, {}, {}},
							},
						},
					},
				},
			},
			fnArgs: fnArgs{p: Point{2, 1}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := CheckTargetSurroundingAttackRangeFn(tt.args.target)
			assert.Equal(t, tt.want, fn(tt.fnArgs.p))
		})
	}
}

func TestClosestEnemy_Explore(t *testing.T) {
	type args struct {
		p *Player
	}
	tests := []struct {
		name string
		args args
		want Move
	}{
		{
			name: "should head towards the closest player",
			args: args{
				p: &Player{
					X:         1,
					Y:         1,
					Direction: "W",
					Game: Game{
						LeaderBoard: []PlayerState{
							{X: 1, Y: 0, Direction: "S"},
							{X: 1, Y: 1, Direction: "W"},
							{X: 0, Y: 2, Direction: "S"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 2, Direction: "S"}}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnRight,
		},
		{
			name: "when tie, take the one who has smaller x and y",
			args: args{
				p: &Player{
					X:         1,
					Y:         1,
					Direction: "W",
					Game: Game{
						LeaderBoard: []PlayerState{
							{X: 0, Y: 0, Direction: "E"},
							{X: 1, Y: 1, Direction: "W"},
							{X: 2, Y: 0, Direction: "W"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {Player: &PlayerState{X: 2, Y: 0, Direction: "W"}}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: WalkForward,
		},
		{
			name: "when first target is blocked, find the next target",
			args: args{
				p: &Player{
					X:         1,
					Y:         1,
					Direction: "W",
					Game: Game{
						LeaderBoard: []PlayerState{
							{X: 0, Y: 0, Direction: "E"},
							{X: 1, Y: 1, Direction: "W"},
							{X: 2, Y: 2, Direction: "S"},
							{X: 0, Y: 3, Direction: "N"},
						},
						Arena: Arena{
							Width:  5,
							Height: 4,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y: 0, WasHit: true, Direction: "E"}}, {}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}, {}},
								{{}, {}, {Player: &PlayerState{X: 2, Y: 2, Direction: "S"}}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 3, Direction: "N"}}, {}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "when first target is blocked, find the next target",
			args: args{
				p: &Player{
					X:         1,
					Y:         1,
					Direction: "S",
					Game: Game{
						LeaderBoard: []PlayerState{
							{X: 0, Y: 0, Direction: "E"},
							{X: 1, Y: 1, Direction: "S"},
							{X: 2, Y: 2, Direction: "S"},
							{X: 0, Y: 3, Direction: "N"},
						},
						Arena: Arena{
							Width:  5,
							Height: 4,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y: 0, WasHit: true, Direction: "E"}}, {}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "S"}}, {}, {}, {}},
								{{}, {}, {Player: &PlayerState{X: 2, Y: 2, Direction: "S"}}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 3, Direction: "N"}}, {}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: WalkForward,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ClosestEnemy{}
			if got := a.Explore(tt.args.p); got != tt.want {
				t.Errorf("Explore() = %v, want %v", got, tt.want)
			}
		})
	}
}