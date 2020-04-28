package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

	for testCase := int64(1); testCase <= testCases; testCase++ {
		n := scanInt(scanner)
		if n < 20 || (n >= 30 && n < 40) || n == 59 {
			fmt.Printf("Case #%d: %s\n", testCase, "IMPOSSIBLE")
		} else if n < 40 {
			fmt.Printf("Case #%d: 1\n", testCase)
		} else {
			c := n / 20
			fmt.Printf("Case #%d: %d\n", testCase, c)
		}

	}

}
func scanInt(scanner *bufio.Scanner) int64 {
	mustScan(scanner)
	n, err := strconv.ParseInt(scanner.Text(), 10, 64)
	assertNoError(err)
	return n
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