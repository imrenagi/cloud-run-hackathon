package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArena(t *testing.T) {
	type args struct {
		w int
		h int
	}
	tests := []struct {
		name string
		args args
		want Arena
	}{
		{
			name: "5 cols, 3 row",
			args: args{
				w: 5,
				h: 3,
			},
			want: Arena{
				Width:  5,
				Height: 3,
				Grid: [][]Cell{
					{{}, {}, {}, {}, {}},
					{{}, {}, {}, {}, {}},
					{{}, {}, {}, {}, {}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArena(tt.args.w, tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArena() = %v, wantAnyOf %v", got, tt.want)
			}
		})
	}
}

func TestArena_GetAdjacent(t *testing.T) {

	type args struct {
		p       Point
		options []AdjacentOption
	}
	tests := []struct {
		name  string
		arena Arena
		args  args
		want  []Point
	}{
		{
			name:  "surround by 8 points",
			arena: NewArena(5, 3),
			args: args{
				p: Point{1, 1},
				options: []AdjacentOption{WithDiagonalAdjacents()},
			},
			want: []Point{
				{0, 0}, {0, 1}, {0, 2},
				{1, 0}, {1, 2},
				{2, 0}, {2, 1}, {2, 2},
			},
		},
		{
			name:  "on the top left corner, surround by 3 points",
			arena: NewArena(5, 3),
			args:  args{
				p: Point{0, 0},
				options: []AdjacentOption{WithDiagonalAdjacents()},
			},
			want: []Point{
				{0, 1}, {1, 0}, {1, 1},
			},
		},
		{
			name:  "on the middle top edge",
			arena: NewArena(5, 3),
			args:  args{
				p: Point{1, 0},
				options: []AdjacentOption{WithDiagonalAdjacents()},
			},
			want: []Point{
				{0, 0}, {0, 1},
				{1, 1},
				{2, 0}, {2, 1},
			},
		},
		{
			name:  "surround by 4 points, without diagonal",
			arena: NewArena(5, 3),
			args: args{
				p: Point{1, 1},
			},
			want: []Point{
				{0, 1},
				{1, 0}, {1, 2},
				{2, 1},
			},
		},
		{
			name:  "on the top left corner, surround by 2 points",
			arena: NewArena(5, 3),
			args:  args{
				p: Point{0, 0},
			},
			want: []Point{
				{0, 1}, {1, 0},
			},
		},
		{
			name:  "on the middle top edge, return no diagonal",
			arena: NewArena(5, 3),
			args:  args{
				p: Point{1, 0},
			},
			want: []Point{
				{0, 0},
				{1, 1},
				{2, 0},
			},
		},
		{
			name:  "on the middle top edge, return no diagonal and only empty cell",
			arena: Arena{
				Width:  5,
				Height: 3,
				Grid:   [][]Cell{
					{{}, {}, {}, {}, {}},
					{{}, {Player: &PlayerState{}}, {}, {}, {}},
					{{}, {}, {}, {}, {}},
				},
			},
			args:  args{
				p: Point{1, 0},
				options: []AdjacentOption{WithEmptyAdjacent()},
			},
			want: []Point{
				{0, 0},
				{2, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := tt.arena.GetAdjacent(context.TODO(), tt.args.p, tt.args.options...)
			assert.Len(t, got, len(tt.want))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAdjacent() = %v, wantAnyOf %v", got, tt.want)
			}

		})
	}
}

func TestArena_Traverse(t *testing.T) {
	type args struct {
		start Point
	}
	tests := []struct {
		name  string
		arena Arena
		args  args
		want  []Point
	}{
		{
			name:  "get correct bfs",
			arena: NewArena(4, 3),
			args:  args{start: Point{1, 1}},
			want: []Point{
				{1, 1}, {0, 0}, {0, 1}, {0, 2},
				{1, 0}, {1, 2},
				{2, 0}, {2, 1}, {2, 2},
				{3, 0}, {3, 1}, {3, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arena.Traverse(context.TODO(), tt.args.start); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Traverse() = %v, wantAnyOf %v", got, tt.want)
			}
		})
	}
}
