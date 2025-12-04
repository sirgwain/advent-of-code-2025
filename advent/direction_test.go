package advent

import "testing"

func Test_direction_minTurns(t *testing.T) {
	tests := []struct {
		name  string
		d     Direction
		other Direction
		want  int
	}{
		{"up->right", DirectionUp, DirectionRight, 1},
		{"up->left", DirectionUp, DirectionLeft, 1},
		{"up->down", DirectionUp, DirectionDown, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.MinTurns(tt.other); got != tt.want {
				t.Errorf("direction.minTurns() = %v, want %v", got, tt.want)
			}
		})
	}
}
