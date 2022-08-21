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
			want: WalkForward,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Attack{
				Player: tt.fields.Player,
			}
			if got := a.Play(); got != tt.want {
				t.Errorf("Play() = %v, want %v", got, tt.want)
			}
		})
	}
}
