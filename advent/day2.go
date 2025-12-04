package advent

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Day2 struct {
	input     [][2]int
	step      int
	solution1 int64
	solution2 int64
}

func (d *Day2) Day() int {
	return 2
}

func (d *Day2) Run(updates chan<- DayUpdate) error {

	for {
		done := d.Progress()
		if done {
			break
		}

		updates <- DayUpdate{
			View:     d.View(),
			Solution: d.ViewSolution(),
			Done:     d.Done(),
		}
	}

	return nil
}

// Init loads in the input from the file and initializes the Day
func (d *Day2) Init(filename string, opts ...Option) (err error) {

	content, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	ranges := strings.Split(string(content), ",")
	input := make([][2]int, len(ranges))

	for i, idRange := range ranges {
		ids := strings.Split(idRange, "-")
		input[i][0], err = strconv.Atoi(ids[0])
		if err != nil {
			return fmt.Errorf("not a number %s %v", ids[0], err)
		}
		input[i][1], err = strconv.Atoi(ids[1])
		if err != nil {
			return fmt.Errorf("not a number %s %v", ids[1], err)
		}
	}

	d.input = input

	return nil
}

// Progress progresses one "step" and returns true if finished
func (d *Day2) Progress() (done bool) {

	idRange := d.input[d.step]

	// loop through each id looking for duplicates
	// 11-22
	// *11*, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, *22*
	for id := idRange[0]; id <= idRange[1]; id++ {
		// part 1
		if isTwoRepeatingNumbers(id) {
			d.solution1 += int64(id)
		}
		// part 2
		if hasRepeatingNumbers(id) {
			d.solution2 += int64(id)
		}
	}

	d.step++

	return d.Done()
}

func (d *Day2) Done() bool {
	return d.step == len(d.input)
}

func (d Day2) View() string {
	step := min(d.step, len(d.input)-1)
	return fmt.Sprintf("S%d ID Range %s - %s",
		step,
		correctResultStyle.Render(strconv.Itoa(d.input[step][0])),
		correctResultStyle.Render(strconv.Itoa(d.input[step][1])),
	)

}

func (d Day2) ViewSolution() string {
	return fmt.Sprintf("solution1: %s, solution2: %s",
		SolutionStyle.Render(strconv.FormatInt(d.solution1, 10)),
		SolutionStyle.Render(strconv.FormatInt(d.solution2, 10)),
	)
}

// isTwoRepeatingNumbers returns true if this number contains two numbers repeating
func isTwoRepeatingNumbers(num int) bool {
	digits := len(strconv.Itoa(num))
	if digits < 2 || digits%2 > 0 {
		// odd number of digits, no duplicates
		return false
	}

	// for a number like 123123 we want 123, 123
	tens := int(math.Pow10(digits / 2))
	left := num / tens
	right := num % tens

	return left == right
}

// hasRepeatingNumbers returns true if this number contains a repeating pattern
func hasRepeatingNumbers(num int) bool {
	str := strconv.Itoa(num)
	digits := len(str)
	if digits < 2 {
		// can't repeat
		return false
	}

	// for a number like 123123, try
	// 123 123 (6 / 2)
	// 12 31 23 (6 / 3)
	// 1 2 3 1 2 3 (6 / 6)
NEXTSPLIT:
	for split := 2; split <= digits; split++ {
		if digits%split > 0 {
			// uneven split, skip it
			continue NEXTSPLIT
		}
		// two splits for 123 123
		// three splits for 12 31 23
		size := digits / split
		for i := size; i <= digits-size; i += size {
			if str[0:size] != str[i:i+size] {
				continue NEXTSPLIT
			}
		}
		return true
	}

	return false
}
