package advent

import "testing"

func Test_checkIn64RangeOverlap(t *testing.T) {
	tests := []struct {
		name string
		r1   int64Range
		r2   int64Range
		want int64Range
	}{
		{name: "no overlap low", r1: int64Range{3, 5}, r2: int64Range{10, 14}, want: int64Range{}},
		{name: "no overlap high", r1: int64Range{10, 14}, r2: int64Range{3, 5}, want: int64Range{}},
		{name: "low overlap", r1: int64Range{3, 12}, r2: int64Range{10, 14}, want: int64Range{3, 14}},
		{name: "high overlap", r1: int64Range{12, 15}, r2: int64Range{10, 14}, want: int64Range{10, 15}},
		{name: "r1 in r2", r1: int64Range{12, 13}, r2: int64Range{10, 14}, want: int64Range{10, 14}},
		{name: "r2 in r1", r1: int64Range{10, 14}, r2: int64Range{12, 13}, want: int64Range{10, 14}},
		{name: "r1 high one away from r2 low", r1: int64Range{10, 14}, r2: int64Range{15, 16}, want: int64Range{10, 16}},
		{name: "r1 low one away from r2 high", r1: int64Range{10, 14}, r2: int64Range{5, 9}, want: int64Range{5, 14}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkIn64RangeOverlap(tt.r1, tt.r2); got != tt.want {
				t.Errorf("checkIn64RangeOverlap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDay5Part2(b *testing.B) {
	d := Day5{}
	if err := d.Init("../inputs/day5.txt", &Options{}); err != nil {
		b.Fatalf("failed to load input %v", err)
	}

	b.Run("sirgwian", func(b *testing.B) {
		var update = func() {}
		for b.Loop() {
			d.part2(update)
			d.calcSolution2()
		}
	})

	b.Run("sollniss", func(b *testing.B) {
		ranges := make([]int64Range, len(d.inputRanges))
		copy(ranges, d.inputRanges)
		for b.Loop() {
			d.part2_sollniss(ranges)
		}
	})

}
