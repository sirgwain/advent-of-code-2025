package advent

import (
	"fmt"
	"strconv"
)

type DayN struct {
	*Options
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
func (d *DayN) Init(filename string, options *Options) (err error) {
	d.Options = options
	return nil
}

func (d *DayN) View() string {
	if d.Quiet {
		return ""
	}
	return ""
}

func (d *DayN) ViewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		solutionStyle.Render(strconv.Itoa(d.solution1)),
		solutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
