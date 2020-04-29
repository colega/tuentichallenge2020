package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

// golang is great but overflow sucks, python is better here
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	testCases := scanInt(scanner)

	for testCase := 1; testCase <= testCases; testCase++ {
		p := scanInt(scanner)
		h, c := solve(p)
		if h == 0 {
			fmt.Printf("Case #%d: IMPOSSIBLE\n", testCase)
		} else {
			fmt.Printf("Case #%d: %d %d\n", testCase, h, c)
		}
	}
}

func solve(p int) (int, int) {
	h := sort.Search(math.MaxInt64, func(h int) bool {
		v := f(1, 1, h)
		return v < h || v > p || v < 0
	}) - 1
	if h < 3 {
		return 0, 0
	}
	//fmt.Printf("found h=%d with f(1, 1, h) = %d\n", h, f(1, 1, h))

	m := sort.Search(math.MaxInt64, func(m int) bool {
		v := f(m, m, h)
		return v < m || v > p || v < 0
	}) - 1
	//fmt.Printf("found m=%d with f(m, m, h) = %d\n", m, f(m, m, h))

	n := m
	if f(m, n+1, h) <= p {
		n++
		//fmt.Printf("using n=m+1=%d with f = %d\n", n, f(m, n, h))
	}

	return h, f(m, n, h)
}

func f(m, n, h int) int {
	k := h - 2
	//fmt.Printf("trying f(M=%d, N=%d, H=%d)\n", m, n, h)
	var s int

	//for i := int(0); i < k; i++ {
	//s += 2*(m+4*i)*(n+4*i) - (m+4*i+2)*(n+4*i+2)
	//}
	sumi := k * (k - 1) / 2
	sumi2 := k * (k - 1) * (2*(k-1) + 1) / 6

	s += sumi * 4 * m
	s += sumi * 4 * n
	s += sumi2 * 16
	s += k * m * n
	s -= k * 2 * m
	s -= k * 2 * n
	s -= 16 * sumi
	s -= 4 * k
	s += (m + 4*k) * (n + 4*k) * 2

	//fmt.Printf("f(M=%d, N=%d, H=%d)=%d\n", m, n, h, s)
	return s
}

func scanInt(scanner *bufio.Scanner) int {
	mustScan(scanner)
	i, err := strconv.ParseInt(scanner.Text(), 10, 64)
	assertNoError(err)
	return int(i)
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
