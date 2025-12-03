package cmd

import (
	"fmt"
	"time"

	"github.com/sirgwain/advent-of-code-2025/advent"
	"github.com/spf13/cobra"
)

type dayRunner interface {
	Run(filename string) error
}

func newRunCmd() *cobra.Command {
	var day int
	var input string
	var visualization bool
	var delay int
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run a day",
		Long:  `run the solution for a day`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var runner dayRunner
			switch day {
			case 1:
				runner = &advent.Day1{}
			case 2:
				runner = &advent.Day2{}
			case 3:
				runner = &advent.Day3{}
				// case 4:
				// 	runner = &advent.Day4{}
				// case 5:
				// 	runner = &advent.Day5{}
				// case 6:
				// 	runner = &advent.Day6{}
				// case 7:
				// 	runner = &advent.Day7{}
				// case 8:
				// 	runner = &advent.Day8{}
				// case 9:
				// 	runner = &advent.Day9{}
				// case 10:
				// 	runner = &advent.Day10{}
				// case 11:
				// 	runner = &advent.Day11{}
				// case 12:
				// 	runner = &advent.Day12{}
			}

			if runner == nil {
				return fmt.Errorf("day %d not found", day)
			}

			// run the visualizer if specified
			if v, ok := runner.(advent.Day); ok && visualization {
				return advent.RunVisual(v, input, advent.WithDelay(delay))
			}

			start := time.Now()
			defer func() { fmt.Printf("\nTime taken %v\n", time.Since(start)) }()

			return runner.Run(input)
		},
	}

	cmd.Flags().IntVarP(&day, "day", "d", 0, "the part to run, a or b")
	cmd.Flags().StringVarP(&input, "input", "i", "", "the input file to load")
	cmd.Flags().IntVar(&delay, "delay", 0, "a delay, in ms to add to the UI")
	cmd.Flags().BoolVarP(&visualization, "visualization", "v", false, "run the visualization for this day, if available")

	cmd.MarkFlagRequired("day")
	cmd.MarkFlagRequired("input")

	return cmd
}

func init() {
	rootCmd.AddCommand(newRunCmd())
}
