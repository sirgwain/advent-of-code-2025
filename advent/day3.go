package advent

import (
	"fmt"
	"strconv"
)

type Day3 struct {
	solution1 int
	solution2 int
}

func (d *Day3) Day() int {
	return 3
}

func (d *Day3) Run(filename string) error {
	if err := d.Init(filename); err != nil {
		return err
	}

	for {
		done := d.Progress()
		if done {
			break
		}

		fmt.Println(d.View())
	}

	fmt.Println(d.ViewSolution())
	return nil
}

// Init loads in the input from the file and initializes the Day
func (d *Day3) Init(filename string) (err error) {
	return nil
}

func (d *Day3) Progress() bool {
	return true
}

func (d Day3) View() string {
	return ""
}

func (d Day3) ViewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		SolutionStyle.Render(strconv.Itoa(d.solution1)),
		SolutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
