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

func (d *DayN) Run(updates chan<- DayUpdate) error {

	updates <- DayUpdate{
		View:     d.View(),
		Solution: d.ViewSolution(),
		Done:     true,
	}
	return nil
}

// Init loads in the input from the file and initializes the Day
func (d *DayN) Init(filename string, opts ...Option) (err error) {
	return nil
}

func (d *DayN) View() string {
	return ""
}

func (d *DayN) ViewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		SolutionStyle.Render(strconv.Itoa(d.solution1)),
		SolutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
