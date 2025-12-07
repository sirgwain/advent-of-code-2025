package cmd

import (
	"fmt"
	"time"

	"github.com/sirgwain/advent-of-code-2025/advent"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	var day int
	var input string
	var visualization bool
	var quiet bool
	var delay int
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run a day",
		Long:  `run the solution for a day`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var d advent.Day
			switch day {
			case 0:
				d = &advent.DayN{}
			case 1:
				d = &advent.Day1{}
			case 2:
				d = &advent.Day2{}
			case 3:
				d = &advent.Day3{}
			case 4:
				d = &advent.Day4{}
			case 5:
				d = &advent.Day5{}
			case 6:
				d = &advent.Day6{}
			case 7:
				d = &advent.Day7{}
				// case 8:
				// 	d = &advent.Day8{}
				// case 9:
				// 	d = &advent.Day9{}
				// case 10:
				// 	d = &advent.Day10{}
				// case 11:
				// 	d = &advent.Day11{}
				// case 12:
				// 	d = &advent.Day12{}
			}

			if d == nil {
				return fmt.Errorf("day %d not found", day)
			}

			// run the visualizer if specified
			if visualization {
				return advent.RunVisual(d, input, advent.WithDelay(delay))
			}

			start := time.Now()
			defer func() { fmt.Printf("\nTime taken %v\n", time.Since(start)) }()

			return advent.Run(d, input, advent.WithQuiet(quiet))
		},
	}

	cmd.Flags().IntVarP(&day, "day", "d", 0, "the part to run, a or b")
	cmd.Flags().StringVarP(&input, "input", "i", "", "the input file to load")
	cmd.Flags().IntVar(&delay, "delay", 0, "a delay, in ms to add to the UI")
	cmd.Flags().BoolVarP(&visualization, "visualization", "v", false, "run the visualization for this day, if available")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "run in quiet mode")

	cmd.MarkFlagRequired("day")
	cmd.MarkFlagRequired("input")
	// no quiet mode when visualizing
	cmd.MarkFlagsMutuallyExclusive("quiet", "visualization")

	return cmd
}

func init() {
	rootCmd.AddCommand(newRunCmd())
}
