package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeaderBoard_GetPlayerByRank(t *testing.T) {
	type args struct {
		rank int
	}
	tests := []struct {
		name string
		l    LeaderBoard
		args args
		want *PlayerState
	}{
		{
			name: "correctly return number",
			l: []PlayerState{
				{
					URL:       "http://testing1",
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    false,
					Score:     5,
				},
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
					X:         1,
					Y:         3,
					Direction: "E",
					WasHit:    false,
					Score:     3,
				},
			},
			args: args{
				rank: 1,
			},
			want: &PlayerState{
				URL:       "http://testing2",
				X:         1,
				Y:         2,
				Direction: "E",
				WasHit:    false,
				Score:     4,
			},
		},
		{
			name: "not found return nil",
			l: []PlayerState{},
			args: args{
				rank: 1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetPlayerByRank(tt.args.rank); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPlayerByRank() = %v, wantAnyOf %v", got, tt.want)
			}
		})
	}
}

func TestLeaderBoard_GetRank(t *testing.T) {
	type args struct {
		p Player
	}
	tests := []struct {
		name string
		l    LeaderBoard
		args args
		want int
	}{
		{
			name: "",
			l: []PlayerState{
				{
					URL:       "http://testing1",
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    false,
					Score:     5,
				},
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
					X:         1,
					Y:         3,
					Direction: "E",
					WasHit:    false,
					Score:     3,
				},
			},
			args: args{
				p: Player{
					Name:         "http://testing3",
					X:            1,
					Y:            3,
					Direction:    "E",
				},
			},
			want: 2,
		},
		{
			name: "",
			l: []PlayerState{
				{
					URL:       "http://testing1",
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    false,
					Score:     5,
				},
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
					X:         1,
					Y:         3,
					Direction: "E",
					WasHit:    false,
					Score:     3,
				},
			},
			args: args{
				p: Player{
					Name:         "http://testing4",
					X:            1,
					Y:            4,
					Direction:    "E",
				},
			},
			want: -1,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetRank(tt.args.p); got != tt.want {
				t.Errorf("GetRank() = %v, wantAnyOf %v", got, tt.want)
			}
		})
	}
}

func TestLeaderBoard_Sort(t *testing.T) {
	tests := []struct {
		name string
		l    LeaderBoard
		want LeaderBoard
	}{
		{
			name: "sort from highest score",
			l: []PlayerState{
				{
					URL:       "http://testing1",
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    false,
					Score:     2,
				},
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
					X:         1,
					Y:         3,
					Direction: "E",
					WasHit:    false,
					Score:     5,
				},
			},
			want: []PlayerState{
				{
					URL:       "http://testing3",
					X:         1,
					Y:         3,
					Direction: "E",
					WasHit:    false,
					Score:     5,
				},
				{
					URL:       "http://testing2",
					X:         1,
					Y:         2,
					Direction: "E",
					WasHit:    false,
					Score:     4,
				},
				{
					URL:       "http://testing1",
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    false,
					Score:     2,
				},
			},
		},
		{
			name: "should tie based on the URL",
			l: []PlayerState{
				{
					URL:       "http://testing1",
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    false,
					Score:     2,
				},
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
					X:         1,
					Y:         3,
					Direction: "E",
					WasHit:    false,
					Score:     4,
				},
			},
			want: []PlayerState{
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
					X:         1,
					Y:         3,
					Direction: "E",
					WasHit:    false,
					Score:     4,
				},
				{
					URL:       "http://testing1",
					X:         1,
					Y:         1,
					Direction: "E",
					WasHit:    false,
					Score:     2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Sort()
			assert.Equal(t, tt.want, tt.l)
		})
	}
}