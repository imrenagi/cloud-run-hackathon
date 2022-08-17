package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlay(t *testing.T) {
	type args struct {
		input ArenaUpdate
	}
	tests := []struct {
		name         string
		args         args
		wantResponse []string
	}{
		{
			name: "should attack on the front",
			wantResponse: []string{"T"},
			args: args{input: ArenaUpdate{
				Links: struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				}{
					Self: struct{
						Href string `json:"href"`
					}{Href: "player1"},
				},
				Arena: struct {
					Dimensions []int                  `json:"dims"`
					State      map[string]PlayerState `json:"state"`
				}{
					Dimensions: []int{4,3},
					State: map[string]PlayerState{
						"player1": {
							X:         1,
							Y:         1,
							Direction: "W",
							WasHit:    false,
							Score:     0,
						},
						"player2": {
							X:         0,
							Y:         1,
							Direction: "S",
							WasHit:    false,
							Score:     0,
						},
					},
				},

			}},
		},
		{
			name: "should not attack ",
			wantResponse: []string{"L", "R", "F"},
			args: args{input: ArenaUpdate{
				Links: struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				}{
					Self: struct{
						Href string `json:"href"`
					}{Href: "player1"},
				},
				Arena: struct {
					Dimensions []int                  `json:"dims"`
					State      map[string]PlayerState `json:"state"`
				}{
					Dimensions: []int{4,3},
					State: map[string]PlayerState{
						"player1": {
							X:         1,
							Y:         1,
							Direction: "W",
							WasHit:    false,
							Score:     0,
						},
						"player2": {
							X:         1,
							Y:         0,
							Direction: "E",
							WasHit:    false,
							Score:     0,
						},
					},
				},

			}},
		},
		{
			name: "should turn left if there is player on left",
			wantResponse: []string{"L"},
			args: args{input: ArenaUpdate{
				Links: struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				}{
					Self: struct{
						Href string `json:"href"`
					}{Href: "player1"},
				},
				Arena: struct {
					Dimensions []int                  `json:"dims"`
					State      map[string]PlayerState `json:"state"`
				}{
					Dimensions: []int{4,3},
					State: map[string]PlayerState{
						"player1": {
							X:         1,
							Y:         1,
							Direction: "N",
							WasHit:    false,
							Score:     0,
						},
						"player2": {
							X:         0,
							Y:         1,
							Direction: "E",
							WasHit:    false,
							Score:     0,
						},
					},
				},

			}},
		},
		{
			name: "should turn right if there is player on right",
			wantResponse: []string{"R"},
			args: args{input: ArenaUpdate{
				Links: struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				}{
					Self: struct{
						Href string `json:"href"`
					}{Href: "player1"},
				},
				Arena: struct {
					Dimensions []int                  `json:"dims"`
					State      map[string]PlayerState `json:"state"`
				}{
					Dimensions: []int{4,3},
					State: map[string]PlayerState{
						"player1": {
							X:         1,
							Y:         1,
							Direction: "N",
							WasHit:    false,
							Score:     0,
						},
						"player2": {
							X:         2,
							Y:         1,
							Direction: "E",
							WasHit:    false,
							Score:     0,
						},
					},
				},

			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse := Play(tt.args.input);
			assert.Contains(t, tt.wantResponse, string(gotResponse))
		})
	}
}
