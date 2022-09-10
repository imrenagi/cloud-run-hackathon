package main

import (
	"reflect"
	"testing"
)

func TestDirection_Opposite(t *testing.T) {
	type fields struct {
		Name   string
		Degree int
	}
	tests := []struct {
		name   string
		fields fields
		want   Direction
	}{
		{
			name: "south",
			fields: fields{
				Name:   South.Name,
				Degree: South.Degree,
			},
			want: North,
		},
		{
			name: "west",
			fields: fields{
				Name:   West.Name,
				Degree: West.Degree,
			},
			want: East,
		},
		{
			name: "north",
			fields: fields{
				Name:   North.Name,
				Degree: North.Degree,
			},
			want: South,
		},
		{
			name: "east",
			fields: fields{
				Name:   East.Name,
				Degree: East.Degree,
			},
			want: West,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Direction{
				Name:   tt.fields.Name,
				Degree: tt.fields.Degree,
			}
			if got := d.Backward(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Backward() = %v, wantAnyOf %v", got, tt.want)
			}
		})
	}
}

