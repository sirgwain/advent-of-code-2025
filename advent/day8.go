package advent

import (
	"bufio"
	"cmp"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Day8 struct {
	*Options
	input     [][3]int
	closestN  int
	solution1 int
	solution2 int
}

type node struct {
	point   [3]int
	index   int
	circuit int
}

type pair struct {
	n1   *node
	n2   *node
	dist int64
}

type circuit struct {
	circuit int
	nodes   []*node
}

func (n *node) String() string {
	return fmt.Sprintf("%d %v", n.index, n.point)
}

func distSquared(a, b [3]int) int64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	dz := a[2] - b[2]
	return int64(dx*dx + dy*dy + dz*dz)
}

func (d *Day8) Day() int {
	return 8
}

func (d *Day8) Run(updates chan<- DayUpdate) error {

	nodes := make([]*node, len(d.input))
	for i, p := range d.input {
		nodes[i] = &node{
			point: p,
			index: i,
		}
	}

	// compute all distances from each junction box to every other junction box, as pairs
	pairs := []pair{}
	for i, n1 := range nodes {
		for j := i + 1; j < len(nodes); j++ {
			n2 := nodes[j]
			pairs = append(pairs, pair{
				n1:   n1,
				n2:   n2,
				dist: distSquared(n1.point, n2.point),
			})
		}
	}

	// sort pairs by distance
	slices.SortFunc(pairs, func(p1, p2 pair) int { return cmp.Compare(p1.dist, p2.dist) })

	circuits := map[int]*circuit{}

	for count := range len(pairs) {
		if count == d.closestN {
			d.recordPart1(circuits)
		}
		pair := pairs[count]
		fmt.Printf("%d: next closest node %v => %v %d\n",
			count,
			data1Style.Render(pair.n1.String()),
			data2Style.Render(pair.n2.String()),
			pair.dist,
		)

		n1 := pair.n1
		n2 := pair.n2

		if n1.circuit != 0 && n1.circuit == n2.circuit {
			fmt.Printf("part of same circuit %s", correctResultStyle.Render(strconv.Itoa(n1.circuit)))
		} else if n1.circuit != 0 && n2.circuit != 0 {
			// merge n2 circuit into n1
			circuits[n1.circuit].nodes = append(circuits[n1.circuit].nodes, (circuits[n2.circuit].nodes)...)
			n2Circuit := n2.circuit
			for _, n := range circuits[n2.circuit].nodes {
				n.circuit = n1.circuit
			}
			delete(circuits, n2Circuit)
			fmt.Printf("already part of existing circuit %s", correctResultStyle.Render(strconv.Itoa(n1.circuit)))
		} else if n1.circuit != 0 {
			// n1 already in a circuit, join n2
			n2.circuit = n1.circuit
			circuits[n1.circuit].nodes = append(circuits[n1.circuit].nodes, n2)
			fmt.Printf("adding to existing circuit %s", correctResultStyle.Render(strconv.Itoa(n1.circuit)))
		} else if n2.circuit != 0 {
			// n2 already in a circuit, join n1
			n1.circuit = n2.circuit
			circuits[n2.circuit].nodes = append(circuits[n2.circuit].nodes, n1)
			fmt.Printf("adding to existing circuit %s", correctResultStyle.Render(strconv.Itoa(n2.circuit)))
		} else {
			// form new circuit
			n1.circuit = count + 1
			n2.circuit = n1.circuit
			circuits[n1.circuit] = &circuit{
				circuit: n1.circuit,
				nodes:   []*node{n1, n2},
			}
			fmt.Printf("creating new circuit %s", correctResultStyle.Render(strconv.Itoa(n1.circuit)))
		}
		fmt.Printf(": %s boxes, (%s)\n\n",
			solutionStyle.Render(strconv.Itoa(len(circuits[n1.circuit].nodes))),
			data2Style.Render(fmt.Sprintf("%s", circuits[n1.circuit].nodes)))

		if len(circuits) == 1 && len(circuits[n1.circuit].nodes) == len(d.input) {
			// found the last pair
			fmt.Printf("found final pair %s %s",
				data1Style.Render(n1.String()),
				data2Style.Render(n2.String()),
			)
			d.solution2 = n1.point[0] * n2.point[0]
			break
		}

	}

	updates <- DayUpdate{
		View:     d.view(),
		Solution: d.viewSolution(),
		Done:     true,
	}
	return nil
}

// after N steps, record part1's score
func (d *Day8) recordPart1(circuits map[int]*circuit) {
	d.solution1 = 1

	sortedCircuits := slices.Collect(maps.Values(circuits))
	slices.SortFunc(sortedCircuits, func(c1, c2 *circuit) int { return cmp.Compare(len(c2.nodes), len(c1.nodes)) })
	for i := range 3 {
		circuit := sortedCircuits[i]
		if len(circuit.nodes) == 0 {
			continue
		}
		fmt.Printf("c: %s -> %s boxes, %s\n",
			correctResultStyle.Render(strconv.Itoa(circuit.circuit)),
			solutionStyle.Render(strconv.Itoa(len(circuit.nodes))),
			data2Style.Render(fmt.Sprintf("%s", circuit.nodes)))
		d.solution1 *= len(circuit.nodes)
	}
}

// Init loads in the input from the file and initializes the Day
func (d *Day8) Init(filename string, options *Options) (err error) {
	d.Options = options

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

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
		z, err := strconv.Atoi(nums[2])
		if err != nil {
			return fmt.Errorf("error parsing number on line: %s", line)
		}

		d.input = append(d.input, [3]int{x, y, z})
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// different data requires different "n closest"
	d.closestN = 1000
	if len(d.input) == 20 {
		d.closestN = 10
	}

	return err
}

func (d *Day8) view() string {
	if d.Quiet {
		return ""
	}
	return ""
}

func (d *Day8) viewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		solutionStyle.Render(strconv.Itoa(d.solution1)),
		solutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
