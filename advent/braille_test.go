package advent_test

import (
	"testing"

	"github.com/sirgwain/advent-of-code-2025/advent"
)

func TestRenderBraille(t *testing.T) {
	tests := []struct {
		name string
		grid [][]byte
		want string
	}{
		{
			name: "empty grid",
			grid: nil,
			want: "",
		},
		{
			name: "single off cell (1x1, all zero)",
			grid: [][]byte{
				{0},
			},
			want: "⠀\n", // U+2800 BLANK BRAILLE
		},
		{
			name: "dot1 only",
			grid: [][]byte{
				{1},
			},
			want: "⠁\n", // U+2801
		},
		{
			name: "dot2 only",
			grid: [][]byte{
				{0},
				{1},
			},
			want: "⠂\n", // U+2802
		},
		{
			name: "dot4 only",
			grid: [][]byte{
				{0, 1},
			},
			want: "⠈\n", // U+2808
		},
		{
			name: "all dots on",
			grid: [][]byte{
				{1, 1},
				{1, 1},
				{1, 1},
				{1, 1},
			},
			want: "⣿\n", // U+28FF FULL BRAILLE
		},
		{
			name: "two cells horizontally",
			grid: [][]byte{
				{1, 0, 0, 1},
			},
			// left = ⠁  dot1
			// right = ⠈ dot4
			want: "⠁⠈\n",
		},
		{
			name: "two cells vertically",
			grid: [][]byte{
				{1},
				{0},
				{0},
				{0},
				{1}, // second tile dot1
			},
			// first row: ⠁
			// second row: ⠁
			want: "⠁\n⠁\n",
		},
		{
			name: "two braille characters from 4x4 grid",
			grid: [][]byte{
				{1, 0, 0, 1},
				{0, 1, 1, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			// left cell:  ⠑ (dot1 + dot4 + dot5)
			// right cell: ⠊ (dot2 + dot4)
			want: "⠑⠊\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := advent.RenderBraille(tt.grid)

			t.Logf("got:\n%s", got)

			if got != tt.want {
				t.Errorf("RenderBraille() = %q, want %q", got, tt.want)
			}
		})
	}
}
