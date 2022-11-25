package utils

import (
	"reflect"
	"testing"
)

func TestMostFrequent(t *testing.T) {
	type K int
	type args struct {
		arr []K
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 K
	}{
		{
			name: "Normal",
			args: args{
				arr: []K{1, 2, 2, 3, 1, 2, 3, 4},
			},
			want:  3,
			want1: 2,
		},
		{
			name: "2 item has same counter",
			args: args{
				arr: []K{1, 2, 2, 3, 1, 2, 1, 4},
			},
			want:  3,
			want1: 2,
		},
		{
			name: "Empty array",
			args: args{
				arr: []K{},
			},
			want:  0,
			want1: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := MostFrequent(tt.args.arr)
			if got != tt.want {
				t.Errorf("MostFrequent() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("MostFrequent() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetRandomSubArray(t *testing.T) {
	type K int
	type args struct {
		arr  []K
		size int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Normal",
			args: args{
				arr:  []K{1, 2, 3, 4, 5, 6},
				size: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]K, len(tt.args.arr))
			copy(arr, tt.args.arr)

			got := GetRandomSubArray(tt.args.arr, tt.args.size)
			if len(got) != tt.args.size {
				t.Errorf("GetRandomSubArray() len = %v, want %v", len(got), tt.args.size)
				return
			}

			if !reflect.DeepEqual(arr, tt.args.arr) {
				t.Errorf("GetRandomSubArray() source array was changed = %v, want %v", arr, tt.args.arr)
				return
			}

		})
	}
}
