package advent

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Day3 struct {
	*Options
	input     []string
	highest2  []int
	highest12 []int
	solution1 int
	solution2 int
}

type day3Workload struct {
	num       int    // the job number
	str       string // the string to evaluate
	highest2  int
	highest12 int
}

func (d *Day3) Day() int {
	return 3
}

func (d *Day3) Run(updates chan<- DayUpdate) error {

	numWorkers := len(d.input)

	jobs := make(chan *day3Workload, len(d.input))    // Channel to queue jobs
	results := make(chan *day3Workload, len(d.input)) // Channel to collect results

	// Worker function
	worker := func(jobs <-chan *day3Workload, results chan<- *day3Workload) error {
		for job := range jobs {
			highest2, err := highestTwoDigits(job.str)
			if err != nil {
				return err
			}

			highest12, err := highestNDigits(job.str, 12)
			if err != nil {
				return err
			}
			job.highest2 = highest2
			job.highest12 = highest12
			results <- job
		}
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Start the workers
	for i := 0; i < numWorkers; i++ {
		go worker(jobs, results)
	}

	// send all the jobs to the workers
	for i := 0; i < len(d.input); i++ {
		jobs <- &day3Workload{
			num: i,
			str: d.input[i],
		}
	}
	close(jobs)

	// Collect results (this blocks until results is closed)
	for i := 0; i < len(d.input); i++ {
		// wait for a result
		result := <-results
		// update with this result
		d.highest2[result.num] = result.highest2
		d.highest12[result.num] = result.highest12
		d.solution1 += result.highest2
		d.solution2 += result.highest12

		if d.Quiet {
			continue
		}
		updates <- DayUpdate{
			View:     d.View(),
			Solution: d.ViewSolution(),
			Done:     false,
		}
	}

	updates <- DayUpdate{
		View:     d.View(),
		Solution: d.ViewSolution(),
		Done:     true,
	}

	return nil
}

// Init loads in the input from the file and initializes the Day
func (d *Day3) Init(filename string, opts ...Option) (err error) {
	d.Options = NewRun(opts...)
	content, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	d.input = strings.Split(string(content), "\n")
	d.highest2 = make([]int, len(d.input))
	d.highest12 = make([]int, len(d.input))
	return nil
}

func (d *Day3) View() string {
	if d.Quiet {
		return "" // reduce output on quiet
	}

	var sb strings.Builder
	for i, input := range d.input {
		if d.highest12[i] == 0 {
			// skip unfinished data
			continue
		}
		sb.WriteString(fmt.Sprintf("S%d %s highest 2: %s, highest 12: %s\n",
			i,
			Data1Style.Render(input),
			correctResultStyle.Render(strconv.Itoa(d.highest2[i])),
			correctResultStyle.Render(strconv.Itoa(d.highest12[i])),
		))
	}
	return sb.String()
}

func (d *Day3) ViewSolution() string {
	return fmt.Sprintf("solution1: %s solution2: %s",
		SolutionStyle.Render(strconv.Itoa(d.solution1)),
		SolutionStyle.Render(strconv.Itoa(d.solution2)),
	)
}

func highestTwoDigits(str string) (int, error) {
	var high1, high2 int

	for i := range str {
		num, err := strconv.Atoi(string(str[i]))
		if err != nil {
			return 0, err
		}

		if num > high1 && i < len(str)-1 {
			high1 = num
			high2 = 0
			continue
		}
		if num > high2 {
			high2 = num
		}
	}

	return high1*10 + high2, nil
}

func highestNDigits(str string, n int) (int, error) {

	high := make([]byte, n)

	j := 0
	// for each digit find the highest number from right to left
	for h := range n {
		// go right to left through the str to find the largest number
		highest := str[len(str)-n+h]
		highestIndex := len(str) - n + h
		for i := len(str) - n + h; i >= j; i-- {
			if str[i] >= highest {
				highest = str[i]
				highestIndex = i
			}
		}
		j = highestIndex + 1
		high[h] = highest
	}

	return strconv.Atoi(string(high))
}

// friend's algorithm for benchmark comparison
func shantz_highestNDigits(str string, n int) (int, error) {

	high := make([]byte, n)

	j := 0
	// for each digit find the highest number from right to left
	for h := range n {
		// go right to left through the str to find the largest number
		highest := str[len(str)-n+h]
		highestIndex := len(str) - n + h
		for i := len(str) - n + h; i >= j; i-- {
			if str[i] >= highest {
				highest = str[i]
				highestIndex = i
			}
		}
		j = highestIndex + 1
		high[h] = highest
	}

	return strconv.Atoi(string(high))
}
