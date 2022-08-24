package main

import (
	"testing"
)

func TestAttack_Play(t *testing.T) {

	type fields struct {
		Player Player
	}
	tests := []struct {
		name   string
		fields fields
		want   Move
	}{
		{
			name: "should attack player in front",
			fields: fields{
				Player: Player{
					X:         1,
					Y:         1,
					Direction: "W",
					WasHit:    false,
					Score:     0,
					Game: Game{
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 1, Direction: "S"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: Throw,
		},
		{
			name: "should turn left if there is enemy on the left",
			fields: fields{
				Player: Player{
					X:         1,
					Y:         1,
					Direction: "N",
					WasHit:    false,
					Score:     0,
					Game: Game{
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "N"}}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "should turn right if there is enemy on the right",
			fields: fields{
				Player: Player{
					X:         1,
					Y:         1,
					Direction: "N",
					WasHit:    false,
					Score:     0,
					Game: Game{
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "N"}}, {Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnRight,
		},
		{
			name: "should just move forward when none is observed",
			fields: fields{
				Player: Player{
					X:         1,
					Y:         1,
					Direction: "N",
					WasHit:    false,
					Score:     0,
					Game: Game{
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "N"}}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: Throw,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Attack{
				Player: &tt.fields.Player,
			}
			if got := a.Play(); got != tt.want {
				t.Errorf("Play() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttack_Play_Explore(t *testing.T) {

	type fields struct {
		Player Player
	}
	tests := []struct {
		skip   bool
		name   string
		fields fields
		want   Move
	}{
		{
			name: "should just move forward when trying to reach the target",
			fields: fields{
				Player: Player{
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    false,
					Score:     0,
					Game: Game{
						Players: []PlayerState{
							{X: 3, Y: 0, Direction: "S"},
							{X: 1, Y: 1, Direction: "E"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {Player: &PlayerState{X: 3, Y: 0, Direction: "S"}}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: WalkForward,
		},
		{
			name: "should just move forward when trying to reach the target",
			fields: fields{
				Player: Player{
					X:         2,
					Y:         1,
					Direction: "E",
					WasHit:    false,
					Score:     0,
					Game: Game{
						Players: []PlayerState{
							{X: 3, Y: 0, Direction: "S"},
							{X: 2, Y: 1, Direction: "E"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {Player: &PlayerState{X: 3, Y: 0, Direction: "S"}}},
								{{}, {}, {Player: &PlayerState{X: 2, Y: 1, Direction: "E"}}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "should avoid surrounding attack range",
			fields: fields{
				Player: Player{
					X:         0,
					Y:         1,
					Direction: "E",
					Game: Game{
						Players: []PlayerState{
							{X: 1, Y: 0, Direction: "S"},
							{X: 0, Y: 1, Direction: "E"},
							{X: 1, Y: 2, Direction: "N"},
						},
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "should avoid surrounding attack range",
			skip: true,
			fields: fields{
				Player: Player{
					X:         1,
					Y:         1,
					Direction: "W",
					Game: Game{
						Players: []PlayerState{
							{X: 0, Y: 0, Direction: "E"},
							{X: 1, Y: 1, Direction: "W"},
							{X: 0, Y: 4, Direction: "N"},
						},
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {}, {}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 4, Direction: "N"}}, {}, {}, {}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "if probably other player will attack us, instead of taking turn, run away from its range instead",
			skip: true,
			fields: fields{
				Player: Player{
					X:         1,
					Y:         1,
					Direction: "W",
					Game: Game{
						Players: []PlayerState{
							{X: 0, Y: 0, Direction: "E"},
							{X: 1, Y: 1, Direction: "W"},
						},
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y: 0, Direction: "S"}}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 3, Direction: "E"}}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
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

			if tt.skip {
				t.Skip()
			}

			a := Attack{
				Player: &tt.fields.Player,
			}
			if got := a.Play(); got != tt.want {
				t.Errorf("Play() = %v, want %v", got, tt.want)
			}
		})
	}
}
