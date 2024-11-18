package player

import (
	"testing"
)

func TestPlayerService(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PlayerService()
		})
	}
}

func TestPlayer(t *testing.T) {
	type args struct {
		power    uint
		wakeCh   chan uint
		name     string
		opponent string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Player(tt.args.power, tt.args.wakeCh, tt.args.name, tt.args.opponent); got != tt.want {
				t.Errorf("Player() = %v, want %v", got, tt.want)
			}
		})
	}
}
