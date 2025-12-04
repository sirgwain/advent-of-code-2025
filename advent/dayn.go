package advent

type DayN struct {
	solution1 int
	solution2 int
}

func (d *DayN) Day() int {
	return 0
}

func (d *DayN) Run(updates chan<- DayUpdate) error {

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
func (d *DayN) Init(filename string, opts ...Option) (err error) {
	return nil
}

// Progress progresses one "step" and returns true if finished
func (d *DayN) Progress() (done bool) {
	return d.Done()
}

func (d *DayN) Done() bool {
	return true
}

func (d *DayN) View() string {
	return ""
}

func (d *DayN) ViewSolution() string {
	return ""
}
