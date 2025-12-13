package advent

import (
	"fmt"
	"strconv"
	"strings"
)

// prerender some styled characters
var (
	renderedStart = data1Style.Render(" S ")
	renderedBeam  = correctResultStyle.Render(" | ")
	renderedSplit = wallStyle.Render(" ^ ")
)

type Day7 struct {
	*Options
	board              [][]rune
	splits             map[Point]bool
	solutionsFromSplit map[Point]int64

	solution1 int
	solution2 int64
}

func (d *Day7) Day() int {
	return 7
}

// Init loads in the input from the file and initializes the Day
func (d *Day7) Init(filename string, options *Options) (err error) {
	d.Options = options
	d.board, err = ReadInputAsRunes(filename)
	if err != nil {
		return err
	}

	d.splits = make(map[Point]bool)
	d.solutionsFromSplit = make(map[Point]int64)
	return nil
}

func (d *Day7) Run(updates chan<- DayUpdate) error {

	d.fireBeams(updates)

	updates <- DayUpdate{
		View:     d.view(),
		Solution: d.viewSolution(),
		Done:     true,
	}
	return nil
}

func (d *Day7) fireBeams(updates chan<- DayUpdate) {
	// starting from the S, fire a beam
	x, y := FindValue(d.board, 'S')
	d.solution2 = d.fireBeam(x, y+1, updates)
}

func (d *Day7) fireBeam(x, y int, updates chan<- DayUpdate) int64 {
	if _, ok := d.solutionsFromSplit[Point{x, y}]; ok {
		// already tried this route
		return d.solutionsFromSplit[Point{x, y}]
	}

	if y >= len(d.board) {
		// finished the board, record it and move on
		if !d.Quiet {
			updates <- DayUpdate{
				View:     d.view(),
				Solution: d.viewSolution(),
				Done:     false,
			}
		}

		return 1
	}

	val := GetBoardValue(x, y, d.board)
	if val == '^' {
		// found split
		d.splits[Point{x, y}] = true
		d.solution1 = len(d.splits)
		d.solutionsFromSplit[Point{x, y}] = 0
		l := d.fireBeam(x-1, y, updates)
		d.solutionsFromSplit[Point{x, y}] = l
		r := d.fireBeam(x+1, y, updates)
		d.solutionsFromSplit[Point{x, y}] += r
		return l + r
	} else {
		d.board[y][x] = '|'
		s := d.fireBeam(x, y+1, updates)
		return s
	}
}

func (d *Day7) view() string {
	if d.Quiet {
		return ""
	}
	var sb strings.Builder
	for y, line := range d.board {
		for x, r := range line {
			switch r {
			case 'S':
				sb.WriteString(renderedStart)
			case '^':
				solutionsfromPosition := d.solutionsFromSplit[Point{x, y}]
				if d.splits[Point{x, y}] {
					sb.WriteString(visitedStyle.Render(fmt.Sprintf("%3d", solutionsfromPosition)))
				} else {
					sb.WriteString(renderedSplit)
				}
			case '|':
				sb.WriteString(renderedBeam)
			default:
				sb.WriteString(" . ")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (d *Day7) viewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		solutionStyle.Render(strconv.Itoa(d.solution1)),
		solutionStyle.Render(strconv.FormatInt(d.solution2, 10)),
	)
}
