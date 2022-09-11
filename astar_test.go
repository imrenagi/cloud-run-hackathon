package main

import (
	"context"
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
		name    string
		fields  fields
		args    args
		want    Path
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
						{{Player: nil}, {Player: &PlayerState{X: 1, Y: 1, Direction: "N"}}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 1},
				dest: Point{1, 0},
			},
			want: []Point{{1, 1}, {1, 0}},
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
						{{Player: nil}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {Player: nil}},
					},
				},
				src:  Point{1, 2},
				dest: Point{1, 0},
			},
			want: []Point{{1, 2}, {1, 1}, {1, 0}},
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
						{{Player: nil}, {Player: &PlayerState{X: 1, Y: 2, Direction: "E"}}, {Player: nil}},
					},
				},
				src:  Point{1, 2},
				dest: Point{1, 0},
			},
			want: []Point{{1, 2}, {1, 1}, {1, 0}},
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
						{{Player: nil}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {Player: nil}},
					},
				},
				src:  Point{1, 2},
				dest: Point{2, 0},
			},
			want: []Point{{1, 2}, {1, 1}, {1, 0}, {2, 0}},
		},
		{
			name: "test search 2 step to north and 1 step to east, but there is obstacle",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: &PlayerState{}}, {Player: nil}},
						{{Player: nil}, {Player: &PlayerState{X: 1, Y: 2, Direction: "N"}}, {Player: nil}},
					},
				},
				src:  Point{1, 2},
				dest: Point{2, 0},
			},
			want: []Point{{1, 2}, {2, 2}, {2, 1}, {2, 0}},
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
						{{Player: nil}, {Player: &PlayerState{X: 1, Y: 1, Direction: "N"}}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 1},
				dest: Point{2, 1},
			},
			want: []Point{{1, 1}, {2, 1}},
		},
		{
			name: "test search 2 step to north and 1 step to west",
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
						{{Player: nil}, {Player: nil}, {Player: &PlayerState{X: 2, Y: 2, Direction: "N"}}},
					},
				},
				src:  Point{2, 2},
				dest: Point{0, 0},
			},
			want: []Point{{2, 2}, {2, 1}, {2, 0}, {1, 0}, {0, 0}},
		},
		{
			name: "test search heading to west and choose the shortest path",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 1},
				dest: Point{0, 0},
			},
			want: []Point{{1, 1}, {0, 1}, {0, 0}},
		},
		{
			name: "test search heading to north and choose the shortest path",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: &PlayerState{X: 1, Y: 1, Direction: "N"}}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 1},
				dest: Point{0, 0},
			},
			want: []Point{{1, 1}, {1, 0}, {0, 0}},
		},
		{
			name: "test search heading to east and choose the shortest path",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 1},
				dest: Point{0, 0},
			},
			want: []Point{{1, 1}, {1, 0}, {0, 0}},
		},
		{
			name: "test search heading to south and choose the shortest path",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: nil}, {Player: &PlayerState{X: 1, Y: 1, Direction: "S"}}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 1},
				dest: Point{0, 0},
			},
			want: []Point{{1, 1}, {0, 1}, {0, 0}},
		},
		{
			name: "test search heading to west and choose the shortest path if there is obstacle",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: nil}, {Player: nil}},
						{{Player: &PlayerState{}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 1},
				dest: Point{0, 0},
			},
			want: []Point{{1, 1}, {1, 0}, {0, 0}},
		},
		{
			name: "no path at all because blocked by obstacle",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  3,
					Height: 3,
					Grid: [][]Cell{
						{{Player: nil}, {Player: &PlayerState{}}, {Player: nil}},
						{{Player: &PlayerState{}}, {Player: &PlayerState{X: 1, Y: 1, Direction: "W"}}, {Player: nil}},
						{{Player: nil}, {Player: nil}, {Player: nil}},
					},
				},
				src:  Point{1, 1},
				dest: Point{0, 0},
			},
			want: nil,
			wantErr: ErrPathNotFound,
		},
		{
			name: "path to diagonal",
			fields: fields{
				distanceCalculator: &ManhattanDistance{},
			},
			args: args{
				a: Arena{
					Width:  4,
					Height: 3,
					Grid: [][]Cell{
						{{}, {}, {}, {}},
						{{}, {Player: &PlayerState{X: 1, Y: 1, Direction: "E"}}, {}, {}},
						{{}, {}, {}, {}},
					},
				},
				src:  Point{1, 1},
				dest: Point{2, 0},
			},
			want: []Point{{1,1}, {2,1}, {2,0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := NewAStar(tt.args.a)
			path, err := as.SearchPath(context.TODO(), tt.args.src, tt.args.dest)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, path)
		})
	}
}

func TestManhattanDistance_Distance(t *testing.T) {
	type args struct {
		p1 Point
		p2 Point
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "",
			args: args{
				p1: Point{X: 0, Y: 0},
				p2: Point{X: 1, Y: 0},
			},
			want: 1.0,
		},
		{
			name: "",
			args: args{
				p1: Point{X: 0, Y: 0},
				p2: Point{X: 1, Y: 1},
			},
			want: 2.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ManhattanDistance{}
			got := m.Distance(tt.args.p1, tt.args.p2)
			assert.InDelta(t, tt.want, got, 0.001, "got %v expect %v", got, tt.want)
		})
	}
}

func Test_ppair_RequiredRotation(t *testing.T) {
	type fields struct {
		F         float64
		X         int
		Y         int
		Direction Direction
	}
	type args struct {
		pt Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "heading west, target is in south",
			fields: fields{
				X:         0,
				Y:         1,
				Direction: West,
			},
			args: args{
				pt: Point{0, 2},
			},
			want:    1,
		},
		{
			name: "heading west, target is in east",
			fields: fields{
				X:         0,
				Y:         1,
				Direction: West,
			},
			args: args{
				pt: Point{1, 1},
			},
			want:    2,
		},
		{
			name: "heading west, target is in north",
			fields: fields{
				X:         0,
				Y:         1,
				Direction: West,
			},
			args: args{
				pt: Point{0, 0},
			},
			want:    1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ppair{
				F:         tt.fields.F,
				X:         tt.fields.X,
				Y:         tt.fields.Y,
				Direction: tt.fields.Direction,
			}
			if got := p.requiredRotation(tt.args.pt); got != tt.want {
				t.Errorf("requiredRotation() = %v, wantAnyOf %v", got, tt.want)
			}
		})
	}
}
