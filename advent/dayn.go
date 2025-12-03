package advent

import (
	"fmt"
	"strconv"
)

type DayN struct {
	solution1 int
	solution2 int
}

func (d *DayN) Day() int {
	return 0
}

func (d *DayN) Run(filename string) error {
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
func (d *DayN) Init(filename string) (err error) {
	return nil
}

func (d *DayN) Progress() bool {
	return true
}

func (d DayN) View() string {
	return ""
}

func (d DayN) ViewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		SolutionStyle.Render(strconv.Itoa(d.solution1)),
		SolutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
