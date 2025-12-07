package advent

import (
	"fmt"
	"strconv"
	"strings"
)

type Day4 struct {
	*Options
	board        [][]rune
	validSquares map[Position]bool
	solution1    int
	solution2    int
}

// prerender some styled characters
var (
	renderedPaperTowel = BoxStyle.Render("@")
)

func (d *Day4) Day() int {
	return 0
}

func (d *Day4) Run(updates chan<- DayUpdate) error {

	iteration := 0
	for {
		for y := 0; y < len(d.board); y++ {
			for x := 0; x < len(d.board[y]); x++ {
				if GetBoardValue(x, y, d.board) != '@' {
					continue
				}
				// check if there are less than 4 papertowels adjacent to us
				pos := Position{x, y}
				adjacentObstacles := countAdjacent(d.board, pos, '@')
				if adjacentObstacles < 4 {
					if iteration == 0 {
						// solution1 is only a single iteration
						d.solution1++
					}
					d.solution2++
					d.validSquares[pos] = true
				}
			}
			if !d.Quiet {
				updates <- DayUpdate{
					View:     d.View(),
					Solution: d.ViewSolution(),
				}
			}
		}
		if len(d.validSquares) == 0 {
			// all done
			break
		}
		// remove valid squares
		for p := range d.validSquares {
			d.board[p.Y][p.X] = '.'
		}
		d.validSquares = map[Position]bool{}
		iteration++
	}

	updates <- DayUpdate{
		View:     d.View(),
		Solution: d.ViewSolution(),
		Done:     true,
	}

	return nil
}

// countAdjacent check for the existence of a rune adjacent to a position
func countAdjacent(board [][]rune, pos Position, r rune) int {
	count := 0
	for _, d := range AdjacentDirections {
		square := pos.addDirection(d)
		if GetBoardValue(square.X, square.Y, board) == r {
			count++
		}
	}

	return count
}

// Init loads in the input from the file and initializes the Day
func (d *Day4) Init(filename string, opts ...Option) (err error) {
	d.Options = NewRun(opts...)
	input, err := ReadInputAsRunes(filename)
	if err != nil {
		return err
	}
	d.board = input
	d.validSquares = make(map[Position]bool)
	return nil
}

func (d *Day4) View() string {

	var sb strings.Builder
	for y, line := range d.board {
		for x, r := range line {
			pos := Position{x, y}
			if d.validSquares[pos] {
				sb.WriteString(BoxHighlightStyle.Render(string(r)))
				continue
			}
			switch r {
			case '@':
				sb.WriteString(renderedPaperTowel)
			default:
				sb.WriteRune(r)
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (d *Day4) ViewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		SolutionStyle.Render(strconv.Itoa(d.solution1)),
		SolutionStyle.Render(strconv.Itoa(d.solution2)),
	)

}
