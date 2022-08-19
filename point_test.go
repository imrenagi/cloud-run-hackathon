package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint_MoveWithDirection(t *testing.T) {
	type fields struct {
		X int
		Y int
	}
	type args struct {
		distance  int
		direction Direction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Point
	}{
		{
			name:   "",
			fields: fields{0, 0},
			args: args{
				direction: South,
				distance:  1,
			},
			want: Point{0, 1},
		},
		{
			name:   "",
			fields: fields{0, 0},
			args: args{
				direction: East,
				distance:  1,
			},
			want: Point{1, 0},
		},
		{
			name: "",
			fields: fields{1, 1},
			args: args{

				direction: North,
				distance:  3,
			},
			want: Point{1, -2},
		},
		{
			name: "",
			fields: fields{1, 1},
			args: args{
				direction: West,
				distance:  3,
			},
			want: Point{-2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Point{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			got := p.TranslateToDirection(tt.args.distance, tt.args.direction)
			// assert.Equal(t, tt.want, got)
			assert.InDelta(t, tt.want.X, got.X, 0.01, "got %v want %v", got.X, tt.want.X)
			assert.InDelta(t, tt.want.Y, got.Y, 0.01, "got %v want %v", got.Y, tt.want.Y)
		})
	}
}
