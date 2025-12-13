package advent

import (
	"bufio"
	"fmt"
	"log/slog"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	gridWidth  = 200
	gridHeight = 200
)

type int64Range struct {
	low  int64
	high int64
}

type Day5 struct {
	*Options
	inputRanges  []int64Range // two ingredient id ranges, i.e. 474206951121632-478696506672479
	inputRange   int64Range   // min/max of the range
	inputIDs     []int64      // an id to check, i.e. 223088071752434
	ranges       []int64Range
	mergedRanges map[int]bool
	grid         [][]byte
	solution1    int
	solution2    int64
}

func (d *Day5) Day() int {
	return 5
}

// Init loads in the input from the file and initializes the Day
func (d *Day5) Init(filename string, options *Options) (err error) {
	d.Options = options
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	idMode := false
	d.inputRange.low = math.MaxInt64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// blank line is new input type
		if strings.TrimSpace(line) == "" {
			idMode = true
			continue
		}

		if idMode {
			num, err := strconv.ParseInt(line, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing number on line: %s", line)
			}
			d.inputIDs = append(d.inputIDs, num)
			continue
		}

		// split ids, 92714816788170-94137721164754
		split := strings.Split(line, "-")
		if len(split) != 2 {
			return fmt.Errorf("error splitting range input line")
		}

		low, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing number on line: %s", line)
		}

		high, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing number on line: %s", line)
		}

		d.inputRange.low = min(low, d.inputRange.low)
		d.inputRange.high = max(high, d.inputRange.high)
		d.inputRanges = append(d.inputRanges, int64Range{low, high})
	}

	slog.Debug(fmt.Sprintf("input range: %d..%d (dist: %d)", d.inputRange.low, d.inputRange.high, (d.inputRange.high - d.inputRange.low)))

	// Allocate 400x200 grid (y-major: grid[y][x])
	d.grid = make([][]byte, gridHeight)
	for y := range d.grid {
		d.grid[y] = make([]byte, gridWidth)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	return nil
}

func (d *Day5) Run(updates chan<- DayUpdate) error {

	d.part1()

	d.part2(func() {
		if !d.Quiet {
			d.calcSolution2()
			updates <- DayUpdate{
				View:     d.view(),
				Solution: d.viewSolution(),
				Done:     false,
			}
		}
	})

	updates <- DayUpdate{
		View:     d.view(),
		Solution: d.viewSolution(),
		Done:     true,
	}
	return nil
}

func (d *Day5) part1() {
NEXTID:
	for _, id := range d.inputIDs {
		for _, idRange := range d.inputRanges {
			if id >= idRange.low && id <= idRange.high {
				d.solution1++
				continue NEXTID
			}
		}
	}
}

func (d *Day5) part2(update func()) {
	// assume all ranges are valid
	d.mergedRanges = make(map[int]bool)
	d.ranges = make([]int64Range, len(d.inputRanges))
	copy(d.ranges, d.inputRanges)
	slices.SortFunc(d.ranges, func(a, b int64Range) int {
		return int(a.low - b.low)
	})

	// keep trying until we no longer merge any ranges
	mergedRange := true
	for {
		if !mergedRange {
			break
		}
		mergedRange = false
		ranges := make([]int64Range, len(d.ranges))
		copy(ranges, d.ranges)
		// empty out the ranges store
		d.ranges = d.ranges[:0]
	NEXTRANGE:
		for _, idRange := range ranges {
			// update each looop
			update()

			// find overlap with existing ranges
			for i, r := range d.ranges {
				merged := checkIn64RangeOverlap(idRange, r)
				if merged == (int64Range{}) {
					continue
				}
				d.ranges[i] = merged
				d.mergedRanges[i] = true
				continue NEXTRANGE
			}

			// add this to the ranges to track
			d.ranges = append(d.ranges, idRange)
		}
	}

	// compute distance inclusive of the range
	d.calcSolution2()
}

func (d *Day5) calcSolution2() {
	d.solution2 = 0
	for _, r := range d.ranges {
		d.solution2 += r.high - r.low + 1
	}
}

// checkIn64RangeOverlap checks if two int64 ranges overlap
// returns empty int64Range if no overlap
func checkIn64RangeOverlap(r1, r2 int64Range) int64Range {

	// 16 - 20 followed by 12 - 18 looks like this
	//	           16 17 18 19 20
	// 12 13 14 15 16 17 18
	// new range is 12 - 20

	if r1.low >= r2.low && r1.high <= r2.high {
		// r1 falls within r2
		return r2
	}

	if r2.low >= r1.low && r2.high <= r1.high {
		// r2 falls within r1
		return r1
	}

	if r1.high+1 < r2.low || r1.low-1 > r2.high {
		// if r2 is 10..15 and r1 is n..8 or 17..n, no overlap
		return int64Range{}
	}

	// if we got here, we have an overlap
	merged := r2
	if r1.low > r2.low {
		// if r is 10..15 and we are 12..n, our high becomes their new high
		merged.high = r1.high
	}
	if r1.high < r2.high {
		// if r is 10..15 and we are n..12, our low becomes their new low
		merged.low = r1.low
	}
	// return the merged result
	return merged
}

func (d *Day5) view() string {
	if d.Quiet {
		// return ""
	}

	var sb strings.Builder
	for i, r := range d.ranges {
		var highStr string
		if d.mergedRanges[i] {
			highStr = visitedStyle.Render(strconv.FormatInt(r.high, 10))
		} else {
			highStr = data2Style.Render(strconv.FormatInt(r.high, 10))
		}
		sb.WriteString(fmt.Sprintf("%s..%s → valid ids: %s\n",
			data1Style.Render(strconv.FormatInt(r.low, 10)),
			highStr,
			correctResultStyle.Render(strconv.FormatInt(r.high-r.low+1, 10)),
		))
	}
	return sb.String()
}

func (d *Day5) viewGrid() string {
	if d.Quiet {
		return ""
	}

	d.buildGrid()

	return RenderBrailleWithColor(d.grid, DensityColor)
}

func (d *Day5) viewSolution() string {
	return fmt.Sprintf("solution1: %s, solution2: %s",
		solutionStyle.Render(strconv.Itoa(d.solution1)),
		solutionStyle.Render(strconv.FormatInt(d.solution2, 10)),
	)
}

func (d *Day5) buildGrid() {

	// clear the grid
	for y := 0; y < len(d.grid); y++ {
		for x := 0; x < len(d.grid[y]); x++ {
			d.grid[y][x] = 0
		}
	}

	span := d.inputRange.high - d.inputRange.low

	numCells := int64(gridWidth * gridHeight)
	cellsPerID := float64(numCells) / float64(span)

	// For each valid range, map it into [0, numCells) and mark cells on.
	for _, r := range d.ranges {
		// Clamp to the global input range, just in case
		low := max(r.low, d.inputRange.low)
		high := min(r.high, d.inputRange.high)

		// If you treat ranges as [low, high), then high <= low is empty.
		if high <= low {
			continue
		}

		startOff := low - d.inputRange.low // offset from global low
		endOff := high - d.inputRange.low

		// Map offsets to linear cell indices
		startIdx := int(float64(startOff) * cellsPerID)
		endIdx := int(float64(endOff) * cellsPerID)

		// Clamp to [0, numCells-1]
		if startIdx < 0 {
			startIdx = 0
		}
		if startIdx >= int(numCells) {
			// Entire range is beyond our grid
			continue
		}
		if endIdx < 0 {
			// Entire range is before our grid
			continue
		}
		if endIdx >= int(numCells) {
			endIdx = int(numCells) - 1
		}
		if endIdx < startIdx {
			endIdx = startIdx
		}

		// Mark all covered cells as "on"
		for idx := startIdx; idx <= endIdx; idx++ {
			y := idx / gridWidth
			x := idx % gridWidth
			if y >= 0 && y < gridHeight {
				d.grid[y][x] += 1
			}
		}
	}

}

// part2_sollniss is a solution found on reddit from @sollniss. Wow, that's much faster than mine. :)
// putting here for benchmarking
func (d *Day5) part2_sollniss(ranges []int64Range) int {
	// https://github.com/sollniss/aoc2025/blob/14c88f9798582e0c187504d75f9d4ffeb137abc3/day5/main.go#L153-L164

	slices.SortFunc(ranges, func(a, b int64Range) int {
		if a.low < b.low {
			return -1
		} else {
			return 1
		}
	})

	res := 0
	var curr int64
	for _, r := range ranges {
		if curr > r.high {
			continue
		}
		from := max(curr, r.low)
		if from <= r.high {
			res += int(r.high-from) + 1
		}

		curr = r.high + 1
	}

	return res
}
