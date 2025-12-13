package advent

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code-2025/advent/color"
)

var (
	renderedDac = lipgloss.NewStyle().Foreground(color.BrightGreen82).Render("dac")
	renderedFft = lipgloss.NewStyle().Foreground(color.BrightLilac177).Render("fft")
)

type Day11 struct {
	*Options
	input         []*day11Node
	linksCache    map[string]int
	svrToDacLinks int
	svrToFftLinks int
	dacToFftLinks int
	fftToDacLinks int
	dacToOutLinks int
	fftToOutLinks int
	solution1     int
	solution2     int
}

type day11Node struct {
	key   string
	links []*day11Node
}

func (n *day11Node) String() string {
	return n.key
}

func (d *Day11) Day() int {
	return 11
}

// Init loads in the input from the file and initializes the Day
func (d *Day11) Init(filename string, options *Options) (err error) {
	d.Options = options

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// aaa: you hhh
	// you: bbb ccc
	// bbb: ddd eee
	// ccc: ddd eee fff

	lines := strings.Split(string(content), "\n")

	nodesByKey := make(map[string]*day11Node, len(lines))
	d.input = make([]*day11Node, 0, len(lines))
	for _, line := range lines {
		split := strings.Split(line, ": ")
		key := strings.TrimSpace(split[0])
		var n *day11Node
		if existing, ok := nodesByKey[key]; ok {
			n = existing
		} else {
			n = &day11Node{key: key}
			nodesByKey[key] = n
		}

		outKeys := strings.Split(strings.TrimSpace(split[1]), " ")
		n.links = make([]*day11Node, 0, len(outKeys))
		for _, key := range outKeys {
			var out *day11Node
			if existing, ok := nodesByKey[key]; ok {
				out = existing
			} else {
				out = &day11Node{key: key}
				nodesByKey[key] = out
			}
			n.links = append(n.links, out)
		}

		d.input = append(d.input, n)
	}
	return nil
}

func (d *Day11) Run(updates chan<- DayUpdate) error {
	// we cache various root -> sub, sub -> out style link counts
	d.linksCache = make(map[string]int, len(d.input)*6)

	if err := d.part1(); err != nil {
		return err
	}

	if err := d.part2(updates); err != nil {
		return err
	}

	updates <- DayUpdate{
		View:     d.view(),
		Solution: d.viewSolution(),
		Done:     true,
	}
	return nil
}

func (d *Day11) getNode(key string) *day11Node {
	idx := slices.IndexFunc(d.input, func(n *day11Node) bool { return n.key == key })
	if idx == -1 {
		return nil
	}
	return d.input[idx]

}

func (d *Day11) traverse(in *day11Node, out string, path []string) int {
	cacheKey := in.key + "_" + out
	if c, ok := d.linksCache[cacheKey]; ok {
		return c
	}

	if in.key == out {
		// fmt.Printf("%s\n", d.viewPath(path))
		return 1
	}

	links := 0
	for _, link := range in.links {
		links += d.traverse(link, out, append(path, link.key))
	}
	d.linksCache[cacheKey] = links
	return links
}

func (d *Day11) part1() error {
	you := d.getNode("you")
	if you == nil {
		return fmt.Errorf("no you found in data")
	}

	d.solution1 = d.traverse(you, "out", []string{you.key})

	return nil
}

func (d *Day11) part2(updates chan<- DayUpdate) error {
	svr := d.getNode("svr")
	if svr == nil {
		return fmt.Errorf("no svr found in data")
	}

	dac := d.getNode("dac")
	if dac == nil {
		return fmt.Errorf("no dac found in data")
	}

	fft := d.getNode("fft")
	if fft == nil {
		return fmt.Errorf("no fft found in data")
	}

	fmt.Printf("\nfinding svr -> out\n")
	d.traverse(svr, "out", []string{svr.key})

	fmt.Printf("\nfinding svr -> dac\n")
	d.svrToDacLinks = d.traverse(svr, "dac", []string{svr.key})
	fmt.Printf("svr -> dac: %s\n", data2Style.Render(strconv.Itoa(d.svrToDacLinks)))

	fmt.Printf("\nfinding dac -> fft\n")
	d.dacToFftLinks = d.traverse(dac, "fft", []string{dac.key})
	fmt.Printf("dac -> fft: %s\n", data2Style.Render(strconv.Itoa(d.dacToFftLinks)))

	fmt.Printf("\nfinding fft -> out\n")
	d.fftToOutLinks = d.traverse(fft, "out", []string{fft.key})
	fmt.Printf("fft -> out: %s\n", data2Style.Render(strconv.Itoa(d.fftToOutLinks)))

	fmt.Printf("\nfinding svr -> fft\n")
	d.svrToFftLinks = d.traverse(svr, "fft", []string{svr.key})
	fmt.Printf("svr -> fft: %s\n", data2Style.Render(strconv.Itoa(d.svrToFftLinks)))

	fmt.Printf("\nfinding fft -> dac\n")
	d.fftToDacLinks = d.traverse(fft, "dac", []string{fft.key})
	fmt.Printf("fft -> dac: %s\n", data2Style.Render(strconv.Itoa(d.fftToDacLinks)))

	fmt.Printf("\nfinding dac -> out\n")
	d.dacToOutLinks = d.traverse(dac, "out", []string{dac.key})
	fmt.Printf("dac -> out: %s\n", data2Style.Render(strconv.Itoa(d.dacToOutLinks)))

	d.solution2 += d.svrToDacLinks * d.dacToFftLinks * d.fftToOutLinks
	d.solution2 += d.svrToFftLinks * d.fftToDacLinks * d.dacToOutLinks

	updates <- DayUpdate{
		View:     d.view(),
		Solution: d.viewSolution(),
		Done:     false,
	}

	return nil
}

func (d *Day11) viewPath(path []string) string {
	var sb strings.Builder
	for i, p := range path {
		if i > 0 {
			sb.WriteRune(' ')
		}

		switch p {
		case "fft":
			sb.WriteString(renderedFft)
		case "dac":
			sb.WriteString(renderedDac)
		default:
			sb.WriteString(p)
		}
	}
	return sb.String()
}

func (d *Day11) viewNode(n *day11Node) string {
	var sb strings.Builder
	for i, out := range n.links {
		if i > 0 {
			sb.WriteRune(' ')
		}
		switch out.key {
		case "fft", "dac":
			sb.WriteString(correctResultStyle.Render(out.key))
		default:
			sb.WriteString(data2Style.Render(out.key))
		}
	}
	return fmt.Sprintf("%s: [%s]",
		data1Style.Render(n.key),
		sb.String(),
	)
}

func (d *Day11) view() string {
	if d.Quiet {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("svr -> dac: %s\n", data2Style.Render(strconv.Itoa(d.svrToDacLinks))))
	sb.WriteString(fmt.Sprintf("dac -> fft: %s\n", data2Style.Render(strconv.Itoa(d.dacToFftLinks))))
	sb.WriteString(fmt.Sprintf("fft -> out: %s\n", data2Style.Render(strconv.Itoa(d.fftToOutLinks))))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("svr -> fft: %s\n", data2Style.Render(strconv.Itoa(d.svrToFftLinks))))
	sb.WriteString(fmt.Sprintf("fft -> dac: %s\n", data2Style.Render(strconv.Itoa(d.fftToDacLinks))))
	sb.WriteString(fmt.Sprintf("dac -> out: %s\n", data2Style.Render(strconv.Itoa(d.dacToOutLinks))))

	return sb.String()
}

func (d *Day11) viewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		solutionStyle.Render(strconv.Itoa(d.solution1)),
		solutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}
