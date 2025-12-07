package advent

import (
	"fmt"
	"time"

	"github.com/sirgwain/advent-of-code-2025/tui"
)

type DayUpdate struct {
	View     string
	Solution string
	Done     bool
}

type Day interface {
	Day() int
	Init(filename string, options *Options) error
	Run(updates chan<- DayUpdate) error
}

func Run(d Day, filename string, opts ...Option) error {
	options := NewRun(opts...)
	if err := d.Init(filename, options); err != nil {
		return err
	}

	updates := make(chan DayUpdate, 16)
	errCh := make(chan error, 1)

	// Run the day in a goroutine
	go func() {
		err := d.Run(updates)
		close(updates)

		errCh <- err
	}()

	// Consume updates as they arrive
	for u := range updates {
		fmt.Printf("%s %s\n", u.View, u.Solution)
	}

	// Return the error from Run
	return <-errCh
}

func RunVisual(d Day, filename string, opts ...Option) error {
	p := tui.NewViewportProgram(tui.NewModel(fmt.Sprintf("Day %d", d.Day())))
	options := NewRun(opts...)

	if err := d.Init(filename, options); err != nil {
		return err
	}

	updates := make(chan DayUpdate, 16) // buffered so Day isn’t blocked by UI speed
	errCh := make(chan error, 1)

	// start the day's work
	go func() {
		err := d.Run(updates)
		if err != nil {
			errCh <- err
		}
		close(errCh)
		close(updates)
	}()

	view := ""
	solution := ""

	// consume updates and feed Bubble Tea
	go func() {
		for u := range updates {
			p.Send(tui.UpdateViewport(u.View, len(u.View)))
			p.Send(tui.UpdateSolution(u.Solution))

			if options.Delay != 0 {
				time.Sleep(time.Millisecond * time.Duration(options.Delay))
			}

			view = u.View
			solution = u.Solution

			if u.Done {
			}
		}
	}()

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("could not start program: %v", err)
	}

	fmt.Printf("%s\n%s\n", view, solution)

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	default:
	}

	return nil
}
