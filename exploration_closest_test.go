package main

import (
	"reflect"
	"testing"
)

func TestCheckTargetSurroundingAttackRangeFn(t *testing.T) {
	t.Fail()
	// t.Skip()
	type args struct {
		target Player
	}
	tests := []struct {
		name string
		args args
		want IsUnblockFn
	}{
		{
			name: "target can attack a point",
			args: args{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckTargetSurroundingAttackRangeFn(tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckTargetSurroundingAttackRangeFn() = %v, want %v", got, tt.want)
			}
		})
	}
}
