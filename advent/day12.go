package advent

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Day12 struct {
	*Options
	tetriminos []Tetrimino
	boards     []day12Board
	solution1  int
	solution2  int
}

// from https://github.com/Broderick-Westrope/tetrigo/blob/main/pkg/tetris/tetrimino.go
// A Tetrimino is a geometric Tetris piece formed by connecting four square blocks (Minos) along their edges.
// Each Tetrimino has a unique shape, position, rotation state, and rotation behavior defined by
// the Super Rotation System (SRS).
type Tetrimino struct {
	// Value is the character identifier for the Tetrimino (I, O, T, S, Z, J, L).
	// This is used internally and may differ from the display representation.
	Value byte

	// Cells represents the shape of the Tetrimino as a 2D grid.
	// True indicates a Mino is present in the cell, false indicates empty space.
	Cells [][]bool

	// Position tracks the coordinate of the Tetrimino's top-left corner in the game matrix.
	// This serves as the reference point for all movement and rotation operations.
	Position Point

	fillCount int
}

// rotateClockwise rotates the Tetrimino clockwise.
// This rotation is done by reversing the order of the rows
// we have square tetriminos so no need to transpose
func (t *Tetrimino) rotateClockwise() {
	// Reverse the order of the rows
	for i, j := 0, len(t.Cells)-1; i < j; i, j = i+1, j-1 {
		t.Cells[i], t.Cells[j] = t.Cells[j], t.Cells[i]
	}
}

// DeepCopy creates a deep copy of the Tetrimino.
// This is useful when you need to modify a Tetrimino without affecting the original.
func (t *Tetrimino) DeepCopy() *Tetrimino {
	cells := make([][]bool, len(t.Cells))
	for i := range t.Cells {
		cells[i] = make([]bool, len(t.Cells[i]))
		copy(cells[i], t.Cells[i])
	}
	return &Tetrimino{
		Value:    t.Value,
		Cells:    cells,
		Position: t.Position,
	}

}

// count how many squares are filled
func (t *Tetrimino) FillCount() int {
	if t.fillCount == 0 {
		for r := range 3 {
			for c := range 3 {
				if t.Cells[r][c] {
					t.fillCount++
				}
			}
		}
	}
	return t.fillCount
}

// Matrix represents the board of cells on which the game is played.
type Matrix [][]byte

// NewMatrix creates a new Matrix with the given height and width.
// It returns an error if the height is less than 20 to allow for a buffer zone of 20 lines.
func NewMatrix(height, width int) Matrix {
	matrix := make(Matrix, height)
	for i := range matrix {
		matrix[i] = make([]byte, width)
	}
	return matrix
}

func (m *Matrix) DeepCopy() *Matrix {
	duplicate := make(Matrix, len(*m))
	for i := range *m {
		duplicate[i] = make([]byte, len((*m)[i]))
		copy(duplicate[i], (*m)[i])
	}
	return &duplicate
}

type day12Board struct {
	width        int
	height       int
	m            Matrix
	requirements []int
}

func (d *Day12) Day() int {
	return 12
}

// Init loads in the input from the file and initializes the Day
func (d *Day12) Init(filename string, options *Options) (err error) {
	d.Options = options
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// 0:
	// ###
	// ##.
	// ##.

	// 1:
	// ###
	// ##.
	// .##

	// 2:
	// .##
	// ###
	// ##.

	// 3:
	// ##.
	// ###
	// ##.

	// 4:
	// ###
	// #..
	// ###

	// 5:
	// ###
	// .#.
	// ###

	// 4x4: 0 0 0 0 2 0
	// 12x5: 1 0 1 0 2 2
	// 12x5: 1 0 1 0 3 2

	lines := strings.Split(string(content), "\n")
	var tmNum byte
	for i, line := range lines {
		// check for board
		if strings.Contains(line, "x") {
			// 4x4: 0 0 0 0 2 0
			split := strings.Split(line, ": ")

			// find dimensions, i.e. 4x4
			dim := strings.Split(split[0], "x")
			w, err := strconv.Atoi(dim[0])
			if err != nil {
				return fmt.Errorf("couldn't determine dimension: %s %v", line, err)
			}
			h, err := strconv.Atoi(dim[1])
			if err != nil {
				return fmt.Errorf("couldn't determine dimension: %s %v", line, err)
			}

			// find requirements, just a list of numbers
			reqStrs := strings.Split(strings.TrimSpace(split[1]), " ")
			reqs := make([]int, len(reqStrs))
			for j, s := range reqStrs {
				n, err := strconv.Atoi(s)
				if err != nil {
					return fmt.Errorf("couldn't requirement: %s %v", line, err)
				}
				reqs[j] = n
			}
			d.boards = append(d.boards, day12Board{
				width:        w,
				height:       h,
				m:            NewMatrix(h, w),
				requirements: reqs,
			})

			continue
		}
		// check for tetrimino
		if strings.Contains(line, ":") {
			tmNum++
			// tetriminos are all 3x3
			// 4:
			// ###
			// #..
			// ###
			cells := make([][]bool, 3)
			for row := range 3 {
				cells[row] = make([]bool, 3)
				for col := range 3 {
					cells[row][col] = lines[row+1][col] == '#'
				}
			}
			tm := Tetrimino{
				Value: tmNum,
				Cells: cells,
			}
			d.tetriminos = append(d.tetriminos, tm)
			i += 3
		}
	}

	return nil
}

func (d *Day12) Run(updates chan<- DayUpdate) error {

	// lol, glad this worked
	for i, b := range d.boards {
		area := b.width * b.height
		minosRequired := 0
		for i, req := range b.requirements {
			minosRequired += d.tetriminos[i].FillCount() * req
		}
		fmt.Printf("board %d area %d, required minos: %d\n", i, area, minosRequired)

		if minosRequired > area {
			fmt.Printf("board %d invalid\n", i)
		} else {
			d.solution1++
		}
	}

	updates <- DayUpdate{
		View:     d.view(),
		Solution: d.viewSolution(),
		Done:     true,
	}
	return nil
}

func (d *Day12) view() string {
	if d.Quiet {
		return ""
	}
	return ""
}

func (d *Day12) viewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		solutionStyle.Render(strconv.Itoa(d.solution1)),
		solutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
