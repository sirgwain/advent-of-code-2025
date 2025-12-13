// DayN is an empty day to be copy/pasted for a new day
// It contains an Init() function to load input, setup the day, and Run to run the puzzle
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

// Init loads in the input from the file and initializes the Day
func (d *DayN) Init(filename string, options *Options) (err error) {
	d.Options = options
	return nil
}

func (d *DayN) Run(updates chan<- DayUpdate) error {

	updates <- DayUpdate{
		View:     d.view(),
		Solution: d.viewSolution(),
		Done:     true,
	}
	return nil
}

func (d *DayN) view() string {
	if d.Quiet {
		return ""
	}
	return ""
}

func (d *DayN) viewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		solutionStyle.Render(strconv.Itoa(d.solution1)),
		solutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
