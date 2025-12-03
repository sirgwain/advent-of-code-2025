package advent

import "testing"

func Test_isTwoRepeatingNumbers(t *testing.T) {
	tests := []struct {
		name string
		num  int
		want bool
	}{
		{name: "11", num: 11, want: true},
		{name: "123123", num: 123123, want: true},
		{name: "123124", num: 123124, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTwoRepeatingNumbers(tt.num)
			if got != tt.want {
				t.Errorf("hasRepeatingNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasRepeatingNumbers(t *testing.T) {
	tests := []struct {
		name string
		num  int
		want bool
	}{
		{name: "11", num: 11, want: true},
		{name: "111", num: 111, want: true},
		{name: "1112", num: 1112, want: false},
		{name: "123123", num: 123123, want: true},
		{name: "123123123", num: 123123123, want: true},
		{name: "565656", num: 565656, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasRepeatingNumbers(tt.num)
			if got != tt.want {
				t.Errorf("hasRepeatingNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
