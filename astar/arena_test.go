package astar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAStar_Search(t *testing.T) {
	type fields struct {
		distanceCalculator DistanceCalculator
	}
	type args struct {
		a    Arena
		src  Point
		dest Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Point
		wantErr error
	}{
		{
			name: "test search 1 step to north",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 1},
				dest: Point{1, 0},
			},
			want: []Point{{1,1}, {1, 0}},
		},
		{
			name: "test search 2 step to north",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 2},
				dest: Point{1, 0},
			},
			want: []Point{{1,2}, {1, 1}, {1, 0}},
		},
		{
			name: "test search 2 step to north and 1 step to east",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 2},
				dest: Point{2, 0},
			},
			want: []Point{{1,2}, {1, 1}, {1, 0}, {2, 0}},
		},
		{
			name: "test search 1 step to east",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 1},
				dest: Point{2, 1},
			},
			want: []Point{{1,1}, {2, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := NewAStar(tt.args.a)
			path, err := as.SearchPath(tt.args.src, tt.args.dest)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, path)
		})
	}
}
