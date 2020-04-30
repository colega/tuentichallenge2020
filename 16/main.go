package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/draffensperger/golp"
)

var debugEnabled = true
var verboseEnabled = false

var skip = flag.Int("skip", 0, "test cases to skip")

func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	testCases := scanInts(scanner)[0]

	for testCase := 1; testCase <= testCases; testCase++ {
		nums := scanInts(scanner)
		F, G := nums[0], nums[1]
		var groups []*group
		for g := 0; g < G; g++ {
			employees := scanInts(scanner)[0]
			floors := scanInts(scanner)
			groups = append(groups, &group{
				employees: employees,
				floors:    floors,
			})
		}
		if testCase <= *skip {
			debug("Skipping Case #%d", testCase)
			continue
		}
		solution := solve(F, groups)
		fmt.Printf("Case #%d: %d\n", testCase, solution)
	}
}

type group struct {
	employees int
	floors    []int
	allocated []bool
}

func solve(F int, groups []*group) int {
	G := len(groups)
	allEmployees := 0
	for _, g := range groups {
		allEmployees += g.employees
	}
	debug("Solve %d floors with %d groups and %d employees", F, G, allEmployees)
	for wcs := 0; wcs <= allEmployees; wcs++ {
		lp := golp.NewLP(G+F, G*F)
		for i, g := range groups {
			sort.Ints(g.floors)
			var rows []golp.Entry
			for _, f := range g.floors {
				rows = append(rows, golp.Entry{
					Col: F*i + f,
					Val: float64(1),
				})
			}
			//debug("Adding constraints: %+v", rows)
			_ = lp.AddConstraintSparse(rows, golp.EQ, float64(g.employees))
		}
		for f := 0; f < F; f++ {
			var rows []golp.Entry
			for i := range groups {
				rows = append(rows, golp.Entry{
					Col: i*F + f,
					Val: 1,
				})
			}
			_ = lp.AddConstraintSparse(rows, golp.LE, float64(wcs))
		}
		var obj []float64
		for i := 0; i < F*G; i++ {
			lp.SetInt(i, true)
			obj = append(obj, float64(i))
		}
		lp.SetObjFn(obj)
		t := lp.Solve()
		if t == golp.OPTIMAL {
			//fmt.Println(lp.WriteToString())
			debug("Optimal solution with wcs=%d", wcs)
			return wcs
		} else if t != golp.INFEASIBLE {
			debug("Solution for wcs=%d: %s\n", wcs, t)
		}
		//fmt.Println(lp.WriteToString())
		//fmt.Printf("Solution for wcs=%d: %s\n", wcs, t)
	}
	panic("NOOO")
}

func debug(msg string, args ...interface{}) {
	if debugEnabled {
		fmt.Fprintf(os.Stderr, msg+"\n", args...)
	}
}

func verbose(msg string, args ...interface{}) {
	if verboseEnabled {
		fmt.Fprintf(os.Stderr, msg+"\n", args...)
	}
}

func scanInts(scanner *bufio.Scanner) []int {
	if !scanner.Scan() {
		panic("can't scan")
	}
	var nums []int
	for _, s := range strings.Split(scanner.Text(), " ") {
		nums = append(nums, atoi(s))
	}
	return nums
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	assertNoError(err)
	return i
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}
