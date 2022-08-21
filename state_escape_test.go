package main

import "testing"

func TestEscape_Play(t *testing.T) {

	type fields struct {
		Player PlayerState
	}
	tests := []struct {
		name   string
		fields fields
		want   Decision
	}{
		{
			name: "opponent is attacking from the top(left) and bottom(right), player heading west",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "W",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {}, {}},
							},
						},
					},
				},
			},
			want: MoveForward,
		},
		{
			name: "opponent is attacking from the top and bottom, player heading south",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "S",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{
							Width:  4,
							Height: 3,
							Grid: [][]Cell{
								{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "S"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnRight,
		},
		{
			name: "opponent is attacking from the top and bottom, player heading north",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "N",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "N"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {}, {}},
							}},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "opponent is attacking from the top and bottom, player heading east",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {}, {}},
							}},
					},
				},
			},
			want: MoveForward,
		},
		{
			name: "opponent is attacking from the right, player heading east",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {}, {}},
							},
						},
					},
				},
			},
			want: MoveForward,
		},
		{
			name: "opponent is attacking from the bottom, player heading west",
			fields: fields{
				Player: PlayerState{
					X:         6,
					Y:         2,
					Direction: "W",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {Player: &PlayerState{X: 6, Y: 2, Direction: "W"}}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {Player: &PlayerState{X: 6, Y: 4, Direction: "N"}}},
							},
						},
					},
				},
			},
			want: MoveForward,
		},
		// ---- diserang beberapa sekaligus
		{
			name: "opponent is attacking from the bottom left and back, player heading west, should move forward",
			fields: fields{
				Player: PlayerState{
					X:         5,
					Y:         3,
					Direction: "W",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {Player: &PlayerState{X: 5, Y: 3, Direction: "W"}}, {Player: &PlayerState{X: 6, Y: 3, Direction: "W"}}},
								{{}, {}, {}, {}, {}, {Player: &PlayerState{X: 5, Y: 4, Direction: "N"}}, {}},
							},
						},
					},
				},
			},
			want: MoveForward,
		},
		{
			name: "opponent is attacking from the bottom left right and back, player heading north, should move forward",
			fields: fields{
				Player: PlayerState{
					X:         5,
					Y:         3,
					Direction: "N",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {Player: &PlayerState{X: 3, Y: 3, Direction: "E"}}, {}, {Player: &PlayerState{X: 5, Y: 3, Direction: "N"}}, {Player: &PlayerState{X: 6, Y: 3, Direction: "W"}}},
								{{}, {}, {}, {}, {}, {Player: &PlayerState{X: 5, Y: 4, Direction: "N"}}, {}},
							},
						},
					},
				},
			},
			want: MoveForward,
		},
		{
			name: "opponent is attacking from the bottom left right and back, player heading west, should turn right",
			fields: fields{
				Player: PlayerState{
					X:         5,
					Y:         3,
					Direction: "W",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {Player: &PlayerState{X: 3, Y: 3, Direction: "E"}}, {}, {Player: &PlayerState{X: 5, Y: 3, Direction: "W"}}, {Player: &PlayerState{X: 6, Y: 3, Direction: "W"}}},
								{{}, {}, {}, {}, {}, {Player: &PlayerState{X: 5, Y: 4, Direction: "N"}}, {}},
							},
						},
					},
				},
			},
			want: MoveForward,
		},
		{
			name: "we are cornered, should attack the front player if any",
			fields: fields{
				Player: PlayerState{
					X:         6,
					Y:         4,
					Direction: "W",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {Player: &PlayerState{X: 6, Y: 3, Direction: "S"}}},
								{{}, {}, {}, {}, {}, {Player: &PlayerState{X: 5, Y: 4, Direction: "E"}}, {Player: &PlayerState{X: 6, Y: 4, Direction: "W"}}},
							},
						},
					},
				},
			},
			want: Fight,
		},
		{
			name: "we are cornered, but no front adjacent",
			fields: fields{
				Player: PlayerState{
					X:         6,
					Y:         4,
					Direction: "E",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {Player: &PlayerState{X: 6, Y: 3, Direction: "S"}}},
								{{}, {}, {}, {}, {}, {Player: &PlayerState{X: 5, Y: 4, Direction: "E"}}, {Player: &PlayerState{X: 6, Y: 4, Direction: "E"}}},
							},
						},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "opponent is attacking from the front",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 2, Y: 1, Direction: "W"}}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "opponent is attacking from the left",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: MoveForward,
		},
		{
			name: "opponent is attacking from the back",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: MoveForward,
		},
		{
			name: "opponent is attacking from the back, but on the edge, should turn to right",
			fields: fields{
				Player: PlayerState{
					X:         0,
					Y:         0,
					Direction: "E",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 0, Direction: "E"}}, {}, {}},
								{{}, {}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnRight,
		},
		{
			name: "opponent is attacking from the front, but on the edge, should turn to left",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         0,
					Direction: "W",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 0, Direction: "W"}}, {}, {}},
								{{}, {}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "opponent is attacking from the back, user is on the edge",
			fields: fields{
				Player: PlayerState{
					X:         0,
					Y:         1,
					Direction: "W",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 1, Direction: "W"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnRight,
		},
		{
			name: "opponent is attacking from the back, user is not on the edge",
			fields: fields{
				Player: PlayerState{
					X:         2,
					Y:         2,
					Direction: "W",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{}, {}, {}, {}},
								{{}, {}, {Player: &PlayerState{X: 2, Y: 2, Direction: "W"}}, {Player: &PlayerState{X: 3, Y: 2, Direction: "W"}}},
							},
						},
					},
				},
			},
			want: MoveForward,
		},
		{
			name: "opponent is attacking from the back, but there is other opponent in front",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "W",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 1, Direction: "W"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {Player: &PlayerState{X: 2, Y: 1, Direction: "W"}}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnRight,
		},
		{
			name: "opponent is attacking other user",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 2, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 3, Y: 1, Direction: "W"}}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "none attacking",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{}, {Player: &PlayerState{X:1, Y: 1, Direction: "E"}}, {}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			want: "F",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Escape{
				Player: tt.fields.Player,
			}
			if got := e.Play(); got != tt.want {
				t.Errorf("Play() = %v, want %v", got, tt.want)
			}
		})
	}
}
