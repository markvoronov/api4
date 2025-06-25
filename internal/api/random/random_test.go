package random

import (
	"math/rand"
	"testing"
)

func NewRandomizer() *rand.Rand {
	return rand.New(rand.NewSource(0))
}

func TestNewRandomString1(t *testing.T) {

	type args struct {
		length int
		r      *rand.Rand
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "6 symbol",
			args: args{6, NewRandomizer()},
			want: "cubyhi",
		},
		{
			name: "0 symbol",
			args: args{0, NewRandomizer()},
			want: "",
		},
		{
			name: "1 symbol",
			args: args{1, NewRandomizer()},
			want: "c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRandomString(tt.args.length, tt.args.r); got != tt.want {
				t.Errorf("NewRandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
