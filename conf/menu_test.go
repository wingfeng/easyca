package conf

import (
	"reflect"
	"testing"
)

func TestGetMenu(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test0", args{
			f: "menu.yaml",
		},
			5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMenu(tt.args.f); !reflect.DeepEqual(len(got.Items), tt.want) {
				t.Errorf("GetMenu() = %v, want %v", len(got.Items), tt.want)
			}
		})
	}
}
