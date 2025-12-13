package advent

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code-2025/advent/color"
)

var (
	renderedRedSquare            = lipgloss.NewStyle().Foreground(color.BloodRed52).Render("#")
	renderedHighlightedRedSquare = lipgloss.NewStyle().Foreground(color.BloodRed52).Background(color.VioletsAreBlue105).Render("#")
	renderedGreenSquare          = lipgloss.NewStyle().Foreground(color.BrightGreen82).Render("X")
)

type Day9 struct {
	*Options
	input          [][2]int
	poly           []Point
	board          [][]byte
	min            Point
	max            Point
	p1             Point
	p2             Point
	validRectangle *bool
	solution1      int
	solution2      int
}

func (d *Day9) Day() int {
	return 9
}

// Init loads in the input from the file and initializes the Day
func (d *Day9) Init(filename string, options *Options) (err error) {
	d.Options = options

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	d.min = Point{math.MaxInt, math.MaxInt}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Split(line, ",")

		x, err := strconv.Atoi(nums[0])
		if err != nil {
			return fmt.Errorf("error parsing number on line: %s", line)
		}
		y, err := strconv.Atoi(nums[1])
		if err != nil {
			return fmt.Errorf("error parsing number on line: %s", line)
		}

		d.input = append(d.input, [2]int{x, y})
		d.min = Point{min(x, d.min.X), min(y, d.min.Y)}
		d.max = Point{max(x, d.max.X), max(y, d.max.Y)}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return err
}

func (d *Day9) Run(updates chan<- DayUpdate) error {

	d.board = MakeBoard[byte](d.max.X+1, d.max.Y+1)
	d.poly = make([]Point, len(d.input))
	for i, p := range d.input {
		point := Point{p[0], p[1]}
		d.poly[i] = point
		d.board[point.Y][point.X] = 0x01

		var p1 Point
		if i > 0 {
			// draw from previous point to the current point
			p1 = d.poly[i-1]
		} else {
			// for the first point in the polygon, draw a line from the last point
			p1 = Point{d.input[len(d.input)-1][0], d.input[len(d.input)-1][1]}
		}
		// draw line from the previous point
		x1 := min(p1.X, point.X)
		x2 := max(p1.X, point.X)
		for x := x1 + 1; x < x2; x++ {
			d.board[point.Y][x] |= 0x02
		}
		y1 := min(p1.Y, point.Y)
		y2 := max(p1.Y, point.Y)
		for y := y1 + 1; y < y2; y++ {
			d.board[y][point.X] |= 0x02
		}
	}

	for i, point := range d.input {
		d.p1 = Point{point[0], point[1]}
		for j := i + 1; j < len(d.input); j++ {
			// find the max area between any two points
			d.p2 = Point{d.input[j][0], d.input[j][1]}
			area := area(d.p1, d.p2)
			d.solution1 = max(d.solution1, area)
			d.validRectangle = nil
			if area > d.solution2 && d.validateRectangle(d.p1, d.p2) {
				d.solution2 = area
			}

			if !d.Quiet {
				updates <- DayUpdate{
					View:     d.view(),
					Solution: d.viewSolution(),
					Done:     false,
				}
			}

		}

	}

	updates <- DayUpdate{
		View:     d.view(),
		Solution: d.viewSolution(),
		Done:     true,
	}
	return nil
}

// for these positions, validate that we have a fully filled area
// ..............
// .......#...#.. (7,1) (11,1)
// ..............
// ..#....#...... (2,3) (7,3)
// ..............
// ..#......#.... (2,5) (9,5)
// ..............
// .........#.#.. (9,7) (11,7)
// ..............
func (d *Day9) validateRectangle(p1 Point, p2 Point) bool {
	// find which coord is up/left
	y1 := min(p1.Y, p2.Y)
	x1 := min(p1.X, p2.X)

	// find which coord is down/right
	y2 := max(p1.Y, p2.Y)
	x2 := max(p1.X, p2.X)

	result := true
	for y := y1; y < y2; y++ {
		if !rowInside(d.poly, y, x1, x2) {
			result := false
			d.validRectangle = &result
			return result
		}
	}

	d.validRectangle = &result
	return result
}

func (d *Day9) view() string {
	if d.Quiet {
		return ""
	}

	var sb strings.Builder
	for y := 0; y < len(d.board); y++ {
		for x := 0; x < len(d.board[y]); x++ {
			switch d.board[y][x] {
			case 0:
				sb.WriteRune('.')
			case 1, 3:
				if d.p1 == (Point{x, y}) || d.p2 == (Point{x, y}) {
					sb.WriteString(renderedHighlightedRedSquare)
				} else {
					sb.WriteString(renderedRedSquare)
				}
			case 2:
				sb.WriteString(renderedGreenSquare)
			}
		}
		sb.WriteRune('\n')
	}

	validity := ""
	if d.validRectangle != nil && !*d.validRectangle {
		validity = incorrectResultStyle.Render("invalid")
	} else if d.validRectangle != nil && *d.validRectangle {
		validity = correctResultStyle.Render("valid")
	}

	return fmt.Sprintf("\n%s\np1: %s, p2: %s, area: %d, %s",
		sb.String(),
		data1Style.Render(d.p1.String()),
		data2Style.Render(d.p2.String()),
		area(d.p1, d.p2),
		validity,
	)
}

func (d *Day9) viewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		solutionStyle.Render(strconv.Itoa(d.solution1)),
		solutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
