package advent

import (
	"fmt"
	"time"

	"github.com/sirgwain/advent-of-code-2025/tui"
)

type Day interface {
	Day() int
	Init(filename string) error
	Progress() bool
	View() string
	ViewSolution() string
}

func RunVisual(d Day, filename string, opts ...Option) error {
	// create a bubbletea program
	p := tui.NewViewportProgram(tui.NewModel(fmt.Sprintf("Day %d", d.Day())))

	options := NewRun(opts...)

	// load in the data, init the day
	if err := d.Init(filename); err != nil {
		return err
	}

	go func() {
		for {
			done := d.Progress()

			view := d.View()
			p.Send(tui.UpdateViewport(view, len(view)))
			p.Send(tui.UpdateSolution(d.ViewSolution()))

			if done {
				break
			}

			if options.Delay != 0 {
				time.Sleep(time.Millisecond * time.Duration(options.Delay))
			}

		}
	}()

	// execute the bubbletea program. This will block until the user pressed q or esc
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("could not start program: %v", err)
	}

	// output the board and final result before exiting the program
	fmt.Println(d.View())
	fmt.Println(d.ViewSolution())

	return nil
}
