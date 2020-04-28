package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	testCases := scanInt(scanner)

	var err error
	for testCase := 1; testCase <= testCases; testCase++ {
		mustScan(scanner)
		input := strings.Split(scanner.Text(), " ")
		nums := make([]int, len(input))
		for i, s := range input {
			nums[i], err = strconv.Atoi(s)
			assertNoError(err)
		}
		f := solve(nums)
		fmt.Printf("Case #%d: %d\n", testCase, f)
	}
}

func solve(inp []int) int {
	x := inp[0]
	m := make(map[int]struct{}, len(inp))
	// yes, we're adding x to m because we don't want to use it
	for i := 0; i < len(inp); i++ {
		m[inp[i]] = struct{}{}
	}
	nums := make([]int, 0, x-(len(m)))
	for i := 1; i < x; i++ {
		if _, ok := m[i]; !ok {
			nums = append(nums, i)
		}
	}

	return count(nums, x)
}

// count solves the coin change problem https://www.geeksforgeeks.org/coin-change-dp-7/
func count(nums []int, x int) int {
	if x == 0 {
		return 1
	}
	if x < 0 {
		return 0
	}
	if len(nums) == 0 && x > 0 {
		return 0
	}

	return count(nums[:len(nums)-1], x) + count(nums, x-nums[len(nums)-1])
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
