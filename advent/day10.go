package advent

import (
	"fmt"
	"math"
	"math/bits"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code-2025/advent/color"
)

var (
	renderedLightOn  = lipgloss.NewStyle().Foreground(color.BrightGreen82).Render(" ● ")
	renderedLightOff = lipgloss.NewStyle().Foreground(color.Gray8).Render(" ○ ")

	renderedButtonOn  = lipgloss.NewStyle().Foreground(color.BrightRed124).Render(" ■ ")
	renderedButtonOff = lipgloss.NewStyle().Foreground(color.Gray8).Render(" □ ")
)

type Day10 struct {
	*Options
	input     []day10Light
	solution1 int
	solution2 int
}

// [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
type day10Light struct {
	light         uint     // bitmask of the light 0b0110
	buttons       []uint   // button bitmasks, i.e. 0b0001, 0b0101, etc...
	buttonIndices [][]int  // (3) (1,3) (2) (2,3) (0,2) (0,1)
	coeffs        [][]int  // {1,1,1,1,1,0}, // (0,1,2,3,4)
	buttonStrs    []string // button bitmasks, i.e. 0b0001, 0b0101, etc...
	joltage       []int    // joltage requirements
}

func (d *Day10) Day() int {
	return 10
}

// Init loads in the input from the file and initializes the Day
func (d *Day10) Init(filename string, options *Options) (err error) {
	d.Options = options
	// format:
	// [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
	// [...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
	// [.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var lineRe = regexp.MustCompile(
		`^\s*\[([.#]+)\]\s*((?:\([0-9,]+\)\s*)+)\{([0-9,]+)\}\s*$`)

	lines := strings.Split(string(content), ("\n"))
	d.input = make([]day10Light, 0, len(lines))
	for _, line := range lines {
		m := lineRe.FindStringSubmatch(line)
		if m == nil {
			return fmt.Errorf("line doesn't match Day10 format: %q", line)
		}

		lights := m[1]       // e.g. ".##."
		buttonsChunk := m[2] // e.g. "(3) (1,3) (2)"
		joltageChunk := m[3] // e.g. "3,5,4,7"

		var l day10Light
		// Parse indicator mask
		var mask uint
		for i, state := range lights {
			if state == '#' {
				mask |= 1 << i
			}
		}
		l.light = mask

		for _, btnStr := range strings.Split(strings.TrimSpace(buttonsChunk), " ") {
			var mask uint

			buttons := []int{}
			for _, lightStr := range strings.Split(btnStr[1:len(btnStr)-1], ",") {
				n, err := strconv.Atoi(lightStr)
				if err != nil {
					return err
				}
				// add 0b0...n...0 as a button
				mask |= 1 << n
				buttons = append(buttons, n)
			}
			l.buttons = append(l.buttons, mask)
			l.buttonIndices = append(l.buttonIndices, buttons)
			l.buttonStrs = append(l.buttonStrs, btnStr)
		}

		for _, joltageStr := range strings.Split(joltageChunk, ",") {
			n, err := strconv.Atoi(joltageStr)
			if err != nil {
				return err
			}
			l.joltage = append(l.joltage, n)
		}

		// Build coeffs: indicator vectors (len == len(goal) == len(joltage))
		numVars := len(l.joltage)
		l.coeffs = make([][]int, 0, len(l.buttonIndices))
		for bi, idxs := range l.buttonIndices {
			indicator := make([]int, numVars)
			for _, idx := range idxs {
				if idx < 0 || idx >= numVars {
					return fmt.Errorf("line %s button %d index %d out of range [0,%d) token=%q",
						line, bi, idx, numVars, l.buttonStrs[bi],
					)
				}
				indicator[idx] = 1
			}
			l.coeffs = append(l.coeffs, indicator)
		}

		d.input = append(d.input, l)
	}

	return nil
}

func (d *Day10) Run(updates chan<- DayUpdate) error {
	d.part1()

	// d.part2()

	// found this solution on reddit: https://www.reddit.com/r/adventofcode/comments/1pk87hl/2025_day_10_part_2_bifurcate_your_way_to_victory/
	// My attempt to implement it was unsuccessful
	for i, l := range d.input {
		sub := solveSingle(l.coeffs, l.joltage)
		fmt.Printf("Line %d/%d: joltage: %v, answer %d\n", i+1, len(d.input), l.joltage, sub)
		d.solution2 += sub
	}

	updates <- DayUpdate{
		View:     d.view(),
		Solution: d.viewSolution(),
		Done:     true,
	}
	return nil
}

func (d *Day10) part1() {
	for _, l := range d.input {

		minPresses, buttons := d.minPressesToToggle(l.light, l.buttons)

		if !d.Quiet {
			fmt.Printf("%s %s %d presses\n", d.viewLight(l.light, len(l.joltage)), d.viewButtons(buttons, l.buttonIndices), minPresses)
			d.solution1 += minPresses
		}
	}

}

// part2 attempts to recursively reduce each joltage to 0 by pressing buttons to make the joltage even, then dividing by 2, then making even, etc
// it doesn't work on all inputs, and it panics if there is an odd pattern that can't be made even
// to fix, it needs to try every odd -> even pattern in a recursive loop, instead of just picking the first one it finds
func (d *Day10) part2() {
	for _, l := range d.input {

		joltage := make([]int, len(l.joltage))
		copy(joltage, l.joltage)
		presses := d.reduceVoltageToZero(&l)
		if !d.Quiet {
			fmt.Printf("%v %d presses\n",
				joltage,
				presses,
			)
		}

		d.solution2 += presses
	}
}

// reduceVoltageToZero gradually presses buttons/divides voltage to try and get to zero
func (d *Day10) reduceVoltageToZero(l *day10Light) int {
	if l.isZeroJoltage() {
		// no more presses require
		return 0
	}

	joltage := make([]int, len(l.joltage))
	copy(joltage, l.joltage)

	oddJoltageBits := l.oddJoltageBits()
	if oddJoltageBits != 0 {
		// turn it odd
		minPresses, buttons := d.minPressesToToggle(l.oddJoltageBits(), l.buttons)
		l.reduceJoltage(buttons)

		if !d.Quiet {
			fmt.Printf("%v -> %v %d presses to make even %s\n",
				joltage,
				l.joltage,
				minPresses,
				d.viewButtons(buttons, l.buttonIndices),
			)
		}

		// keep track of these presses
		return minPresses + d.reduceVoltageToZero(l)
	} else {
		// divide voltage by 2
		l.halveJoltage()
		if !d.Quiet {
			fmt.Printf("%v -> %v halved joltage\n",
				joltage,
				l.joltage,
			)
		}

		return 2 * d.reduceVoltageToZero(l)
	}
}

// minPressesToToggle finds the minimum presses to toggle lights into a configuration
func (d *Day10) minPressesToToggle(desired uint, buttons []uint) (presses int, bestButtons uint) {

	presses = math.MaxInt
	bestButtons = 0
	// fmt.Printf("target: %s\n\n", d.viewLight(l.light))
	// go through each combination of button presses
	for perm := uint(1); perm < (1 << len(buttons)); perm++ {
		// fmt.Printf("perm: %2d, press: ", perm)
		var state uint
		for j := 0; j < len(buttons); j++ {
			// check if this button is on in this permutation
			if perm&(1<<j) != 0 {
				// xor to toggle the lights for this button
				// 0b0001
				// 0b0101 xor
				// =====
				// 0b0100 - toggles on light 2 and off light 0
				state ^= buttons[j]
				// fmt.Printf("b%d %5s %s  ", j+1, l.buttonStrs[j], d.viewButton(l.buttons[j]))
			}
		}

		// fmt.Printf("\nstate: %s\n\n", d.viewLight(state))

		// check if we found a match
		if state == desired {
			count := bits.OnesCount(uint(perm))
			if count < presses {
				presses = count
				bestButtons = perm
			}
			// fmt.Printf("%s %d presses\n", d.viewLight(l.light), best)
		}
	}

	return presses, bestButtons
}

func (l *day10Light) oddJoltageBits() uint {
	var mask uint
	for i, j := range l.joltage {
		if j%2 == 1 {
			// odd number
			mask |= 1 << i
		}
	}
	return mask
}

func (l *day10Light) isZeroJoltage() bool {
	for _, j := range l.joltage {
		if j > 0 {
			return false
		}
	}
	return true
}

// joltageAfterPressing shows the joltage result after pressing buttons
// in this case we are going towards 0
func (l *day10Light) reduceJoltage(buttons uint) {
	for b := range bits.Len(buttons) {
		if buttons&(1<<b) != 0 {
			// button b was pressed
			for _, light := range l.buttonIndices[b] {
				l.joltage[light]--
			}
		}
	}
}

func (l *day10Light) halveJoltage() {
	for i, j := range l.joltage {
		l.joltage[i] = j / 2
	}
}

func (d *Day10) viewLight(l uint, numLights int) string {
	var sb strings.Builder
	for i := range numLights {
		if l&(1<<i) != 0 {
			sb.WriteString(renderedLightOn)
		} else {
			sb.WriteString(renderedLightOff)
		}
	}
	return sb.String()
}
func (d *Day10) viewButton(b uint) string {
	var sb strings.Builder
	for i := range 10 {
		if b&(1<<i) != 0 {
			sb.WriteString(renderedButtonOn)
		} else {
			sb.WriteString(renderedButtonOff)
		}
	}
	return sb.String()
}

func (d *Day10) viewButtons(buttons uint, buttonIndices [][]int) string {
	var sb strings.Builder

	buttonsStr := fmt.Sprintf("%4b", buttons)
	_ = buttonsStr
	for i, b := range buttonIndices {
		if buttons&(1<<i) != 0 {
			iStr := fmt.Sprintf("%6b", 1<<i)
			_ = iStr
			sb.WriteString(fmt.Sprintf("B%d %v ", i, b))
		}
	}
	return sb.String()
}

func (d *Day10) view() string {
	if d.Quiet {
		return ""
	}
	return ""
}

func (d *Day10) viewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		solutionStyle.Render(strconv.Itoa(d.solution1)),
		solutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}

type vecPattern struct {
	v    []int
	cost int
}

// patterns returns press patterns you can make by pressing the buttons in various combinations
// along with the cost (num button presses) to make it
// coeffs should be like this (from buttons)
//
//	[][]int{
//	  {1,1,1,1,1,0}, // (0,1,2,3,4)
//	  {1,0,0,1,1,0}, // (0,3,4)
//	  {1,1,1,0,1,1}, // (0,1,2,4,5)
//	  {0,1,1,0,0,0}, // (1,2)
//	}
//
// returns vecPatterns like this
// {v: [1,1,1,0,1,1], cost: 1}
// {v: [3,2,2,2,3,1], cost: 3}
// {v: [0,1,1,0,0,0], cost: 1}
// from reddit, translated from python: https://www.reddit.com/r/adventofcode/comments/1pk87hl/2025_day_10_part_2_bifurcate_your_way_to_victory/
func patterns(coeffs [][]int) []vecPattern {
	numButtons := len(coeffs)
	if numButtons == 0 {
		return nil
	}
	numVars := len(coeffs[0])

	// key -> minCost
	minCost := make(map[string]int, 1<<min(numButtons, 20))
	// key -> representative vector (so we can return []vecPattern without re-decoding)
	vecByKey := make(map[string][]int, 1<<min(numButtons, 20))

	for mask := 0; mask < (1 << numButtons); mask++ {
		p := make([]int, numVars)
		cost := 0
		for i := 0; i < numButtons; i++ {
			if (mask & (1 << i)) == 0 {
				continue
			}
			cost++
			ci := coeffs[i]
			for j := 0; j < numVars; j++ {
				p[j] += ci[j]
			}
		}

		key := fmt.Sprintf("%v", p)
		if prev, ok := minCost[key]; !ok || cost < prev {
			minCost[key] = cost
			// keep the corresponding vector for this key
			// (safe because p is newly allocated each loop)
			vecByKey[key] = p
		}
	}

	out := make([]vecPattern, 0, len(minCost))
	for key, cost := range minCost {
		out = append(out, vecPattern{v: vecByKey[key], cost: cost})
	}
	return out
}

// from reddit, translated from python: https://www.reddit.com/r/adventofcode/comments/1pk87hl/2025_day_10_part_2_bifurcate_your_way_to_victory/
func solveSingle(coeffs [][]int, goal []int) int {
	pats := patterns(coeffs)

	const INF = 1_000_000
	memo := make(map[string]int, 1<<16)

	var rec func(g []int) int
	rec = func(g []int) int {
		allZero := true
		for _, x := range g {
			if x != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			return 0
		}

		k := fmt.Sprintf("%v", g)
		if v, ok := memo[k]; ok {
			return v
		}

		best := INF
		for _, pat := range pats {
			ok := true
			for i := range g {
				pi := pat.v[i]
				gi := g[i]
				if pi > gi || (pi&1) != (gi&1) {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}

			newGoal := make([]int, len(g))
			for i := range g {
				newGoal[i] = (g[i] - pat.v[i]) / 2
			}

			cand := pat.cost + 2*rec(newGoal)
			if cand < best {
				best = cand
			}
		}

		memo[k] = best
		return best
	}

	return rec(goal)
}
