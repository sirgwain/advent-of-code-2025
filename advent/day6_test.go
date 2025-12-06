package advent

import "testing"

func TestDay6_byteSliceToNumber(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		bytes []byte
		want  int
	}{
		// {name: "0", bytes: []byte{'0'}, want: 0},
		{name: "23", bytes: []byte("23"), want: 23},
		{name: " 434", bytes: []byte(" 434"), want: 434},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Day6

			if got := d.byteSliceToNumber(tt.bytes); got != tt.want {
				t.Errorf("byteSliceToNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDay6_part2(t *testing.T) {
	tests := []struct {
		name  string
		board [][]byte
		want  int
	}{
		{
			name: "4+431+623",
			board: [][]byte{
				{'6', '4', ' '},
				{'2', '3', ' '},
				{'3', '1', '4'},
				{'+', ' ', ' '},
			},
			want: 1058,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Day6
			d.board = tt.board
			d.part2(nil)

			if got := d.solution2; got != tt.want {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}

		})
	}
}
