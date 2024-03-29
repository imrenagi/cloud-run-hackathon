package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscape_Play(t *testing.T) {

	type fields struct {
		Player Player
	}
	tests := []struct {
		skip      bool
		name      string
		fields    fields
		wantAnyOf []Move
	}{
		{
			name: "opponent is attacking from the top(left) and bottom(right), player heading west",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "opponent is attacking from the top and bottom, player heading south",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{TurnLeft, TurnRight},
		},
		{
			name: "opponent is attacking from the top and bottom, player heading north",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{TurnRight, TurnLeft},
		},
		{
			name: "opponent is attacking from the top and bottom, player heading east",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "opponent is attacking from the right, player heading east",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "opponent is attacking from the bottom, player heading west",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{WalkForward},
		},
		// ---- diserang beberapa sekaligus
		{
			name: "opponent is attacking from the bottom left and back, player heading west, should move forward",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "we are corner, tried to escape when we got attacked from 3 direction",
			fields: fields{
				Player: Player{
					X:         2,
					Y:         0,
					Direction: "E",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {Player: &PlayerState{X: 2, Y: 0, Direction: "E"}}, {}, {Player: &PlayerState{X: 4, Y: 0, Direction: "W"}}, {}, {}},
								{{}, {}, {Player: &PlayerState{X: 2, Y: 1, Direction: "N"}}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "we are corner, tried to escape when we got attacked from 3 direction",
			fields: fields{
				Player: Player{
					X:         2,
					Y:         0,
					Direction: "N",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y: 0, Direction: "E"}}, {}, {Player: &PlayerState{X: 2, Y: 0, Direction: "N"}}, {}, {Player: &PlayerState{X: 4, Y: 0, Direction: "W"}}, {}, {}},
								{{}, {}, {Player: &PlayerState{X: 2, Y: 1, Direction: "N"}}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{TurnRight, TurnLeft},
		},
		{
			name: "opponent is attacking from the bottom left right and back, player heading north, should move forward",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "opponent is attacking from the bottom left right and back, player heading west, should turn right",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{TurnRight},
		},
		{
			name: "opponent is attacking from the bottom left right and back, player heading west, should turn right",
			fields: fields{
				Player: Player{
					X:         4,
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
								{{}, {}, {}, {Player: &PlayerState{X: 3, Y: 3, Direction: "E"}}, {Player: &PlayerState{X: 4, Y: 3, Direction: "W"}}, {}, {Player: &PlayerState{X: 6, Y: 3, Direction: "W"}}},
								{{}, {}, {}, {}, {}, {Player: &PlayerState{X: 5, Y: 4, Direction: "N"}}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{TurnRight, TurnLeft},
		},
		{
			name: "we are cornered, should attack the front player if any",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{Throw},
		},
		{
			name: "we are cornered, should turn to right if trapped and hit 3 times",
			fields: fields{
				Player: Player{
					X:         6,
					Y:         4,
					Direction: "W",
					trappedCount: 3,
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
			wantAnyOf: []Move{TurnRight},
		},
		{
			name: "we are cornered, should turn to left if trapped and hit 3 times",
			fields: fields{
				Player: Player{
					X:         6,
					Y:         4,
					Direction: "N",
					trappedCount: 3,
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {Player: &PlayerState{X: 6, Y: 3, Direction: "S"}}},
								{{}, {}, {}, {}, {}, {Player: &PlayerState{X: 5, Y: 4, Direction: "E"}}, {Player: &PlayerState{X: 6, Y: 4, Direction: "N"}}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{TurnLeft},
		},

		{
			name: "we are cornered, but no front adjacent",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{TurnLeft},
		},
		{
			name: "we are cornered from dPair (has some adjacents)",
			fields: fields{
				Player: Player{
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
								{{}, {}, {}, {}, {}, {}, {Player: &PlayerState{X: 6, Y: 2, Direction: "S"}}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {Player: &PlayerState{X: 4, Y: 4, Direction: "E"}}, {}, {Player: &PlayerState{X: 6, Y: 4, Direction: "E"}}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{TurnLeft},
		},
		{
			name: "we are cornered from dPair (has some adjacents)",
			fields: fields{
				Player: Player{
					X:         0,
					Y:         1,
					Direction: "E",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y:0, Direction: "S"}}, {}, {}, {}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y:1, Direction: "E"}}, {}, {Player: &PlayerState{X: 2, Y:1, Direction: "W"}}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y:3, Direction: "N"}}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "should not move forward when right or left is empty and enemy is attacking from the front",
			fields: fields{
				Player: Player{
					X:         0,
					Y:         1,
					Direction: "E",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{}, {}, {}, {}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y:1, Direction: "E"}}, {}, {}, {Player: &PlayerState{X: 3, Y:1, Direction: "W"}}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{TurnLeft, TurnRight},
		},
		{
			name: "should not move forward when right or left is empty and enemy is attacking from the front and back",
			fields: fields{
				Player: Player{
					X:         1,
					Y:         1,
					Direction: "E",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 2,
							Grid: [][]Cell{
								{{}, {}, {}, {}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y:1, Direction: "E"}}, {Player: &PlayerState{X: 1, Y:1, Direction: "E"}}, {}, {Player: &PlayerState{X: 3, Y:1, Direction: "W"}}, {}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{TurnLeft},
		},

		{
			name: "we are cornered from dPair (has some adjacents)",
			fields: fields{
				Player: Player{
					X:         0,
					Y:         1,
					Direction: "S",
					Game: Game{
						Arena: Arena{
							Width:  7,
							Height: 5,
							Grid: [][]Cell{
								{{Player: &PlayerState{X: 0, Y:0, Direction: "S"}}, {}, {}, {}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y:1, Direction: "S"}}, {}, {Player: &PlayerState{X: 2, Y:1, Direction: "W"}}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y:3, Direction: "N"}}, {}, {}, {}, {}, {}, {}},
								{{}, {}, {}, {}, {}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "we are cornered from dPair (has some adjacents)",
			fields: fields{
				Player: Player{
					X:         1,
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
								{{}, {Player: &PlayerState{X: 1, Y: 3, Direction: "S"}}, {}, {}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 4, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 4, Direction: "E"}}, {}, {}, {Player: &PlayerState{X: 4, Y: 4, Direction: "W"}}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "we are cornered from dPair (has some adjacents)",
			fields: fields{
				Player: Player{
					X:         2,
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
								{{}, {Player: &PlayerState{X: 1, Y: 3, Direction: "S"}}, {}, {}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 4, Direction: "E"}}, {}, {Player: &PlayerState{X: 2, Y: 4, Direction: "E"}}, {}, {Player: &PlayerState{X: 4, Y: 4, Direction: "W"}}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{TurnLeft},
		},
		{
			name: "opponent is attacking from the front",
			fields: fields{
				Player: Player{
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
								{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "E"}}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{TurnLeft},
		},
		{
			name: "opponent is attacking from the front, should immediately turn to avoid attack",
			fields: fields{
				Player: Player{
					X:         0,
					Y:         1,
					Direction: "E",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 3,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {}, {Player: &PlayerState{X: 2, Y: 1, Direction: "W"}}, {}},
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{TurnLeft, TurnRight},
		},
		{
			name: "opponent is attacking from the front, should immediately turn when get attack",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{TurnLeft, TurnRight},
		},
		{
			name: "opponent is attacking from the left",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "opponent is attacking from the back",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "opponent is attacking from the back, but on the edge, should turn to right",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{TurnRight},
		},
		{
			name: "opponent is attacking from the front, but on the edge, should turn to left",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{TurnLeft},
		},
		{
			name: "opponent is attacking from the back, user is on the edge",
			fields: fields{
				Player: Player{
					X:         0,
					Y:         1,
					Direction: "W",
					WasHit:    true,
					Score:     0,
					Game: Game{
						Arena: Arena{Width: 4, Height: 2,
							Grid: [][]Cell{
								{{}, {}, {}, {}},
								{{Player: &PlayerState{X: 0, Y: 1, Direction: "W"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{TurnRight},
		},
		{
			name: "opponent is attacking from the back, user is not on the edge",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{WalkForward},
		},
		{
			name: "opponent is attacking from the back, but there is other opponent in front",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{TurnRight, TurnLeft},
		},
		{
			name: "opponent is attacking other user",
			fields: fields{
				Player: Player{
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
			wantAnyOf: []Move{TurnLeft, TurnRight},
		},
		{
			name: "none attacking",
			fields: fields{
				Player: Player{
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
								{{}, {}, {}, {}},
							},
						},
					},
				},
			},
			wantAnyOf: []Move{"F"},
		},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}
			e := Escape{
				Player: &tt.fields.Player,
			}
			got := e.Play(context.TODO())
			assert.Contains(t, tt.wantAnyOf, got, "Play() = %v, wantAnyOf %v", got, tt.wantAnyOf)
		})
	}
}

func TestBraveEscapeDecorator_Play(t *testing.T) {
	type fields struct {
		Escaper Escaper
	}
	tests := []struct {
		name   string
		fields fields
		want   Move
	}{
		{
			name: "only one opponent is attacking from the front, attack back",
			fields: fields{
				Escaper: &mockEscaper{
					p: &Player{
						X:         1,
						Y:         1,
						Direction: "E",
						WasHit:    true,
						Score:     0,
						Game: Game{
							Arena: Arena{Width: 4, Height: 3,
								Grid: [][]Cell{
									{{}, {}, {}, {}},
									{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E", WasHit: true}}, {}, {Player: &PlayerState{X: 3, Y: 1, Direction: "W"}}},
									{{}, {}, {}, {}},
								},
							},
						},
					},
				},
			},
			want: Throw,
		},
		{
			name: "only one opponent is attacking from the right, turn right",
			fields: fields{
				Escaper: &mockEscaper{
					p: &Player{
						X:         1,
						Y:         1,
						Direction: "E",
						WasHit:    true,
						Score:     0,
						Game: Game{
							Arena: Arena{Width: 4, Height: 3,
								Grid: [][]Cell{
									{{}, {}, {}, {}},
									{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E", WasHit: true}}, {}, {}},
									{{}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {}, {}},
								},
							},
						},
					},
				},
			},
			want: TurnRight,
		},
		{
			name: "only one opponent is attacking from the left, turn left",
			fields: fields{
				Escaper: &mockEscaper{
					p: &Player{
						X:         1,
						Y:         1,
						Direction: "E",
						WasHit:    true,
						Score:     0,
						Game: Game{
							Arena: Arena{Width: 4, Height: 3,
								Grid: [][]Cell{
									{{}, {Player: &PlayerState{X: 1, Y: 0, Direction: "S"}}, {}, {}},
									{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E", WasHit: true}}, {}, {}},
									{{}, {}, {}, {}},
								},
							},
						},
					},
				},
			},
			want: TurnLeft,
		},
		{
			name: "only one opponent is attacking from the back, escape",
			fields: fields{
				Escaper: &mockEscaper{
					p: &Player{
						X:         1,
						Y:         1,
						Direction: "E",
						WasHit:    true,
						Score:     0,
						Game: Game{
							Arena: Arena{Width: 4, Height: 3,
								Grid: [][]Cell{
									{{}, {}, {}, {}},
									{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E", WasHit: true}}, {}, {}},
									{{}, {}, {}, {}},
								},
							},
						},
					},
				},
			},
			want: "MOCK",
		},
		{
			name: "more than one opponent is attacking, escape",
			fields: fields{
				Escaper: &mockEscaper{
					p: &Player{
						X:         1,
						Y:         1,
						Direction: "E",
						WasHit:    true,
						Score:     0,
						Game: Game{
							Arena: Arena{Width: 4, Height: 2,
								Grid: [][]Cell{
									{{}, {}, {}, {}},
									{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E", WasHit: true}}, {}, {Player: &PlayerState{X: 3, Y: 1, Direction: "W"}}},
								},
							},
						},
					},
				},
			},
			want: "MOCK",
		},
		{
			name: "consecutive hit is greater than 3, escape",
			fields: fields{
				Escaper: &mockEscaper{
					p: &Player{
						X:         1,
						Y:         1,
						Direction: "E",
						WasHit:    true,
						Score:     0,
						consecutiveHitCount: 4,
						Game: Game{
							Arena: Arena{Width: 4, Height: 2,
								Grid: [][]Cell{
									{{}, {}, {}, {}},
									{{Player: &PlayerState{X: 0, Y: 1, Direction: "E"}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E", WasHit: true}}, {}, {Player: &PlayerState{X: 3, Y: 1, Direction: "W"}}},
								},
							},
						},
					},
				},
			},
			want: "MOCK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &BraveEscapeDecorator{
				Escaper: tt.fields.Escaper,
			}
			if got := e.Play(context.TODO()); got != tt.want {
				t.Errorf("Play() = %v, wantAnyOf %v", got, tt.want)
			}
		})
	}
}

type mockEscaper struct {
	p *Player
}

func (m mockEscaper) Play(ctx context.Context) Move {
	return "MOCK"
}

func (m mockEscaper) GetPlayer() *Player {
	return m.p
}