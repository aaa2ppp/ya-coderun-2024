package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

func run(in io.Reader, out io.Writer) error {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var n int
	if _, err := fmt.Fscan(br, &n); err != nil {
		return err
	}

	problems := make([][2]int, n)
	var maxDay int
	for i := 0; i < n; i++ {
		var d, w int
		if _, err := fmt.Fscan(br, &d, &w); err != nil {
			return err
		}
		problems[i] = [2]int{d, w}
		maxDay = max(maxDay, d)
	}

	sort.Slice(problems, func(i, j int) bool {
		return problems[i][1] > problems[j][1]
	})

	days := make([]int, maxDay+1)
	days[0] = -1
	totalStress := 0
	for i := range problems {
		if !solve(problems[i], days) {
			totalStress += problems[i][1]
		}
	}

	fmt.Fprintln(bw, totalStress)
	return nil
}

func solve(p [2]int, days []int) bool {
	for i := p[0]; i > 0; {
		if days[i] == 0 {
			next := i - 1
			days[i] = next
			for next >= 0 && days[next] != 0 {
				next = days[next]
			}
			for i := p[0]; i != next; {
				i, days[i] = days[i], next
			}
			return true
		}
		i = days[i]
	}
	return false
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
