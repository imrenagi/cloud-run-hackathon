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
						},
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
						},
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
						Arena: Arena{Width: 4, Height: 3},
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
			},
			want: TurnRight,
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
						Arena: Arena{Width: 4, Height: 3},
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
						Arena: Arena{Width: 4, Height: 3},
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
						},
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
						},
						PlayersByPosition: map[string]PlayerState{
							"5,3": {
								X:         5,
								Y:         3,
								Direction: "W",
							},
							"5,4": {
								X:         5,
								Y:         4,
								Direction: "N",
							},
							"6,3": {
								X:         6,
								Y:         3,
								Direction: "W",
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
						},
						PlayersByPosition: map[string]PlayerState{
							"3,3": {
								X:         3,
								Y:         3,
								Direction: "E",
							},
							"5,3": {
								X:         5,
								Y:         3,
								Direction: "N",
							},
							"5,4": {
								X:         5,
								Y:         4,
								Direction: "N",
							},
							"6,3": {
								X:         6,
								Y:         3,
								Direction: "W",
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
						},
						PlayersByPosition: map[string]PlayerState{
							"3,3": {
								X:         3,
								Y:         3,
								Direction: "E",
							},
							"5,3": {
								X:         5,
								Y:         3,
								Direction: "W",
							},
							"5,4": {
								X:         5,
								Y:         4,
								Direction: "N",
							},
							"6,3": {
								X:         6,
								Y:         3,
								Direction: "W",
							},
						},
					},
				},
			},
			want: MoveForward,
		},
		// -----
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
						Arena: Arena{Width: 4, Height: 3},
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
			},
			want: TurnRight,
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
						Arena: Arena{Width: 4, Height: 3},
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
						Arena: Arena{Width: 4, Height: 3},
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
			},
			want: TurnRight,
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
						Arena: Arena{Width: 4, Height: 3},
						PlayersByPosition: map[string]PlayerState{
							"1,1": {
								X:         1,
								Y:         1,
								Direction: "W",
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
						Arena: Arena{Width: 4, Height: 3},
						PlayersByPosition: map[string]PlayerState{
							"3,2": {
								X:         1,
								Y:         1,
								Direction: "W",
							},
						},
					},
				},
			},
			want: "F",
		},
		{
			name: "opponent is attacking from the back, but there is other opponent",
			fields: fields{
				Player: PlayerState{
					X:         1,
					Y:         1,
					Direction: "W",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3},
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
			},
			want: "R",
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
						Arena: Arena{Width: 4, Height: 3},
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
			},
			want: "R",
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
						Arena: Arena{Width: 4, Height: 3},
						PlayersByPosition: map[string]PlayerState{
							"1,1": {
								X:         1,
								Y:         1,
								Direction: "E",
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
