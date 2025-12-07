package advent

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Day6 struct {
	*Options
	input           [][]int
	inputOperations []string
	board           [][]byte
	step            int
	problem         string
	problemSolution int
	viewStr         strings.Builder
	solution1       int
	solution2       int
}

func (d *Day6) Day() int {
	return 6
}

// Init loads in the input from the file and initializes the Day
func (d *Day6) Init(filename string, options *Options) (err error) {
	d.Options = options
	content, err := os.ReadFile(filename)

	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	d.input = make([][]int, len(lines)-1)
	d.board = make([][]byte, 0, len(lines))
	for i, line := range lines {
		d.board = append(d.board, []byte(line)) // part 2, build a byte board
		split := strings.Split(string(line), " ")
		for _, s := range split {
			trimmed := strings.TrimSpace(s)
			if trimmed == "" {
				continue
			}
			if trimmed == "*" || trimmed == "+" {
				// operand
				d.inputOperations = append(d.inputOperations, trimmed)
				continue
			}
			num, err := strconv.Atoi(trimmed)
			if err != nil {
				return fmt.Errorf("error parsing number: %s", trimmed)
			}
			d.input[i] = append(d.input[i], num)
		}
	}

	return nil
}

func (d *Day6) Run(updates chan<- DayUpdate) error {

	d.part1()
	d.part2(updates)

	updates <- DayUpdate{
		View:     d.view(),
		Solution: d.viewSolution(),
		Done:     true,
	}
	return nil
}

func (d *Day6) part1() {
	// go top to bottom
	for i := 0; i < len(d.input[0]); i++ {
		operation := d.inputOperations[i]

		result := d.input[0][i]
		for j := 1; j < len(d.input); j++ {
			switch operation {
			case "*":
				result *= d.input[j][i]
			case "+":
				result += d.input[j][i]
			}
		}

		d.solution1 += result
	}
}

func (d *Day6) part2(updates chan<- DayUpdate) {
	// left to right/right to left doesn't matter
	// so go left to right
	// looking for a pattern like this, with a space at the end
	//  79  338
	//  84  921
	// 562  5154
	// 814  6586
	// *    +

	result := 0
	operator := byte(0)
	var sb strings.Builder
NEXTPROBLEM:
	for x := 0; x < len(d.board[0]); x++ {
		column := make([]byte, len(d.board))
		// build the column, top to bottom
		empty := true
		for y := 0; y < len(d.board)-1; y++ {
			column[y] = d.board[y][x]
			if column[y] != byte(' ') {
				empty = false
			}
		}
		// operator is always the first
		if operator == 0 {
			operator = d.board[len(d.board)-1][x]
		}

		// go to the next problem if we encounter a space or we're at the end
		if empty {
			d.problemSolution = result
			d.solution2 += result
			result = 0
			operator = 0

			d.problem = sb.String()
			sb.Reset()

			if !(d.Quiet) {
				updates <- DayUpdate{
					View:     d.view(),
					Solution: d.viewSolution(),
					Done:     false,
				}
			}

			continue NEXTPROBLEM
		}

		num := d.byteSliceToNumber(column)
		if result != 0 {
			sb.WriteString(fmt.Sprintf(" %s ", string(operator)))
		}
		sb.WriteString(strconv.Itoa(num))
		if result == 0 {
			result = num
		} else {
			switch operator {
			case '*':
				result *= num
			case '+':
				result += num
			}
		}

		// one last update for the final number
		if x == len(d.board[0])-1 {
			d.solution2 += result
		}
		d.step++
	}
}

// byteSliceToNumber converts a byteslice like " 423" to 423
func (d *Day6) byteSliceToNumber(bytes []byte) int {
	num := 0
	tens := 1
	for i := len(bytes) - 1; i >= 0; i-- {
		if bytes[i] == ' ' || bytes[i] == 0 {
			continue
		}
		num = tens*int(bytes[i]-'0') + num
		tens *= 10
	}
	return num
}

func (d *Day6) view() string {
	if d.Quiet {
		return ""
	}

	line := fmt.Sprintf("%s = %s\n",
		data1Style.Render(d.problem),
		correctResultStyle.Render(strconv.Itoa(d.problemSolution)),
	)
	d.viewStr.WriteString(line)
	return d.viewStr.String()
}

func (d *Day6) viewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		solutionStyle.Render(strconv.Itoa(d.solution1)),
		solutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
