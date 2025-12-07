package advent

import "testing"

func Test_highestTwoDigits(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    int
		wantErr bool
	}{
		{name: "987654321111111", str: "987654321111111", want: 98, wantErr: false},
		{name: "811111111111119", str: "811111111111119", want: 89, wantErr: false},
		{name: "234234234234278", str: "234234234234278", want: 78, wantErr: false},
		{name: "818181911112111", str: "818181911112111", want: 92, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := highestTwoDigits(tt.str)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("highestTwoDigits() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("highestTwoDigits() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("highestTwoDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_highestNDigits(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		n       int
		want    int
		wantErr bool
	}{
		// 3 to 4 digits for debugger
		{name: "987654321111111 n=3", str: "987654321111111", n: 3, want: 987, wantErr: false},
		{name: "811111111111119 n=3", str: "811111111111119", n: 3, want: 819, wantErr: false},
		{name: "234234234234278 n=4", str: "234234234234278", n: 4, want: 4478, wantErr: false},
		{name: "818181911112111 n=3", str: "818181911112111", n: 3, want: 921, wantErr: false},

		// 12 digits
		{name: "987654321111111 n = 12", str: "987654321111111", n: 12, want: 987654321111, wantErr: false},
		{name: "811111111111119 n = 12", str: "811111111111119", n: 12, want: 811111111119, wantErr: false},
		{name: "234234234234278 n = 12", str: "234234234234278", n: 12, want: 434234234278, wantErr: false},
		{name: "818181911112111 n = 12", str: "818181911112111", n: 12, want: 888911112111, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := highestNDigits(tt.str, tt.n)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("highestNDigits() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("highestNDigits() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("highestNDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDay4Part2(b *testing.B) {
	d := Day3{}
	if err := d.Init("../inputs/day3.txt", &Options{}); err != nil {
		b.Fatalf("failed to init %v", err)
	}
	data := d.input

	b.Run("sirgwain", func(b *testing.B) {

		b.ResetTimer()
		for b.Loop() {
			for _, str := range data {
				highestNDigits(str, 12)
			}
		}
	})

	b.Run("shantz", func(b *testing.B) {

		b.ResetTimer()
		for b.Loop() {
			for _, str := range data {
				shantz_highestNDigits(str, 12)
			}
		}
	})

}
