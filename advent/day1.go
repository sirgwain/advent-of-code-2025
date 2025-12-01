package advent

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Day1 struct {
}

func (d *Day1) Run(part int, filename string, opts ...Option) error {
	switch part {
	case 1:
		return d.part1(filename)
	case 2:
		return d.part2(filename)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

func (d *Day1) readInput(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var slice []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		num, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, fmt.Errorf("error parsing number on line: %s", line)
		}

		if line[0] == 'L' {
			slice = append(slice, -num)
		} else {
			slice = append(slice, num)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return slice, nil
}

func (d *Day1) part1(filename string) error {
	// dial stars at 50
	dial := 50
	password := 0

	nums, err := d.readInput(filename)
	if err != nil {
		return err
	}

	for _, num := range nums {
		dial += num
		for dial < 0 {
			dial += 100
		}
		for dial >= 100 {
			dial = dial - 100
		}

		if dial == 0 {
			password++
		}

		fmt.Printf("Dial: %d -> %d\n", num, dial)
	}

	fmt.Printf("day1: %s password = %d\n", filename, password)
	return nil
}

func (d *Day1) part2(filename string) error {
	// dial stars at 50
	dial := 50
	password := 0

	nums, err := d.readInput(filename)
	if err != nil {
		return err
	}

	for _, num := range nums {
		start := dial
		dial += num

		for dial < 0 {
			dial += 100
			if start != 0 {
				password++
			}
			start = dial
		}
		for dial >= 100 {
			dial = dial - 100
			if dial != 0 {
				password++
			}
		}

		if dial == 0 {
			password++
		}

		fmt.Printf("Dial: %d + %d -> %d = %d times\n", start, num, dial, password)
	}

	fmt.Printf("day1: %s password = %d\n", filename, password)
	return nil
}
