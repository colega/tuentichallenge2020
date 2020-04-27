package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type nodeState int

const (
	UNPROCESSED nodeState = iota
	PROCESSING
	PROCESSED
)

type node struct {
	edge  []int
	state nodeState
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	testCases := scanInt(scanner)

	for testCase := 1; testCase <= testCases; testCase++ {
		matchesCount := scanInt(scanner)
		p := 0
		graph := make(map[int]*node)
		for m := 0; m < matchesCount; m++ {
			mustScan(scanner)
			s := strings.Split(scanner.Text(), " ")
			if len(s) != 3 {
				panic(fmt.Errorf("not three digits: %s", s))
			}
			a, err := strconv.Atoi(s[0])
			assertNoError(err)
			b, err := strconv.Atoi(s[1])
			assertNoError(err)
			res, err := strconv.Atoi(s[2])
			assertNoError(err)
			if res == 1 {
				// b is always better than a
				a, b = b, a
			}
			//fmt.Printf("Testcase %d with %d matches players %d %d => %d\n", testCase, matchesCount, a, b, res)
			p = max(p, a)
			p = max(p, b)

			if graph[a] == nil {
				graph[a] = &node{}
			}

			graph[a].edge = append(graph[a].edge, b)
		}

		m := 0
		for i := 1; i <= p; i++ {
			if graph[i] != nil {
				m = findMax(graph, 1)
				break
			}
		}
		fmt.Printf("Case #%d: %d\n", testCase, m)
	}

}

func findMax(graph map[int]*node, n int) int {
	for _, edge := range graph[n].edge {
		if graph[edge] == nil {
			return edge // no-one is better than this
		}
		return findMax(graph, edge)
	}
	panic("no edges")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func scanInt(scanner *bufio.Scanner) int {
	mustScan(scanner)
	i, err := strconv.Atoi(scanner.Text())
	assertNoError(err)
	return i
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func mustScan(scanner *bufio.Scanner) {
	if !scanner.Scan() {
		panic(fmt.Errorf("can't scan"))
	}
}
