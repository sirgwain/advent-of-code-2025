package advent

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Day1 struct {
	input     []int
	step      int
	dial      int
	num       int // the current input num being processed
	solution1 int
	solution2 int
}

func (d *Day1) Day() int {
	return 1
}

func (d *Day1) Run(updates chan<- DayUpdate) error {

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
func (d *Day1) Init(filename string, opts ...Option) (err error) {
	// dial starts at 50
	d.dial = 50

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		num, err := strconv.Atoi(line[1:])
		if err != nil {
			return fmt.Errorf("error parsing number on line: %s", line)
		}

		if line[0] == 'L' {
			d.input = append(d.input, -num)
		} else {
			d.input = append(d.input, num)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return err
}

func (d *Day1) Progress() bool {
	d.num = d.input[d.step]
	start := d.dial
	d.dial += d.num

	for d.dial < 0 {
		d.dial += 100
		if start != 0 {
			d.solution2++
		}
		start = d.dial
	}
	for d.dial >= 100 {
		d.dial = d.dial - 100
		if d.dial != 0 {
			d.solution2++
		}
	}

	if d.dial == 0 {
		d.solution1++
		d.solution2++
	}
	d.step++

	return d.Done()
}

func (d *Day1) Done() bool {
	return d.step == len(d.input)
}

func (d Day1) View() string {
	command := ""
	if d.num < 0 {
		command = fmt.Sprintf("L%00d", -d.num)
	} else {
		command = fmt.Sprintf("R%00d", d.num)
	}

	return fmt.Sprintf("S%d - dial: %002d, command: %s",
		d.step,
		d.dial,
		command,
	)
}

func (d Day1) ViewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		SolutionStyle.Render(strconv.Itoa(d.solution1)),
		SolutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
