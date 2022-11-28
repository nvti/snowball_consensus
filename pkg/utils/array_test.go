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
		name      string
		args      args
		wantCount int
		wantValue K
		wantErr   bool
	}{
		{
			name: "Normal",
			args: args{
				arr: []K{1, 2, 2, 3, 1, 2, 3, 4},
			},
			wantCount: 3,
			wantValue: 2,
			wantErr:   false,
		},
		{
			name: "2 item has same counter",
			args: args{
				arr: []K{1, 2, 2, 3, 1, 2, 1, 4},
			},
			wantCount: 3,
			wantValue: 2,
			wantErr:   false,
		},
		{
			name: "Empty array",
			args: args{
				arr: []K{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := MostFrequent(tt.args.arr)
			if (err != nil) != tt.wantErr {
				t.Errorf("MostFrequent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got != tt.wantCount {
				t.Errorf("MostFrequent() got = %v, want %v", got, tt.wantCount)
				return
			}
			if !reflect.DeepEqual(got1, tt.wantValue) {
				t.Errorf("MostFrequent() got1 = %v, want %v", got1, tt.wantValue)
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
