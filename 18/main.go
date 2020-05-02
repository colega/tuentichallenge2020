/*
Package main contains code that I woudldn't like to debug, really.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var debugEnabled, verboseEnabled bool

const IMPOSSIBLE = -1
const UNKNOWN = -2

func main() {
	var onlyTestCase int
	flag.IntVar(&onlyTestCase, "only", 0, "Run only one test case. 0 runs all")
	flag.BoolVar(&verboseEnabled, "verbose", false, "Verbose logging")
	flag.BoolVar(&verboseEnabled, "debug", false, "Debug logging")
	flag.Parse()
	if verboseEnabled {
		debugEnabled = true
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	testCases := atoi(scanner.Text())

	for testCase := 1; testCase <= testCases; testCase++ {
		scanner.Scan()
		L := atoi(scanner.Text())
		inp := make(input, 0, L)
		for i := 0; i < L; i++ {
			scanner.Scan()
			inp = append(inp, &inputLine{s: scanner.Text()})
		}

		if onlyTestCase > 0 && testCase != onlyTestCase {
			continue
		}

		t0 := time.Now()
		all := inp.String()
		debug("Solving case %d with %d lines with a total of %d characters:\n=========\n%s\n========= end of test case %d", testCase, len(inp), len(all), all, testCase)
		solution := solve(inp)
		if solution >= 0 {
			fmt.Printf("Case #%d: %d\n", testCase, solution)
		} else if solution == IMPOSSIBLE {
			fmt.Printf("Case #%d: IMPOSSIBLE\n", testCase)
		} else if solution == UNKNOWN {
			fmt.Printf("Case #%d: UNKNOWN\n", testCase)
		}
		debug("TestCase %d: %d characters in %s", testCase, len(all), time.Since(t0))
	}
}

func solve(lines input) int {
	debug("Is valid ESC: %t", isValidESC(lines))
	if !isValidESC(lines) {
		return makeValidESC(lines, 0)
	}
	return makeValidLolmaoList([]byte(lines.String()))
}

type input []*inputLine

func (inp input) String() string {
	var s []string
	for _, line := range inp {
		s = append(s, line.String())
	}
	return strings.Join(s, "\n")
}

// isLOLMAOList checks whether the input is a LOLMAO list, assuming it's a valid LOLMAO input
func (inp input) isLOLMAOList() bool {
	return len(inp) > 0 && len(inp[0].s) > 0 && inp[0].s[0] == '['
}

// lolmaoLiteralToLolmaoList assumes that the input is a literal, i,e, it does not contain opening or closing bracket
func (inp input) lolmaoLiteralToLolmaoList() (input, bool) {
	cp := inp.copy()
	if len(cp) == 0 {
		// empty input cant be a list
		return nil, false
	} else if len(cp) == 1 {
		if cp[0].len() < 2 {
			// one liner with one char can't be a list
			// we can make two oneliners a list, and we can make a two-chars a list
			return nil, false
		}
		cp[0] = &inputLine{s: "[" + cp[0].s[1:cp[0].len()-1] + "]"}
		return cp, true
	}

	if len(cp[0].s) > 0 {
		// first line has chars, and we know its not a list start, so we can make it a list start
		cp[0] = &inputLine{s: "[" + cp[0].s[1:]}
	} else {
		// first line is empty, replace the newline
		cp = cp[1:]
		cp[0] = &inputLine{s: "[" + cp[0].s}
	}

	lastLineIndex := len(cp) - 1
	lastLine := cp[lastLineIndex]

	lastLineLen := lastLine.len()
	if lastLineIndex == 0 && lastLineLen == 1 {
		// our last line is what just a `[` we built before, see two empty lines testcase
		return nil, false
	} else if lastLineLen > 0 {
		cp[lastLineIndex] = &inputLine{s: lastLine.s[:lastLineLen-1] + "]"}
	} else if len(cp) == 1 {
		// can't join
		return nil, false
	} else {
		// join the last line replacing the newline by ]
		cp = cp[:lastLineIndex]
		cp[lastLineIndex-1] = &inputLine{s: cp[lastLineIndex-1].s + "]"}
	}
	return cp, true
}

func (inp input) replace(line, pos int, s string) input {
	cp := inp.copy()
	ls := cp[line].s
	cp[line] = &inputLine{s: ls[:pos] + s + ls[pos+1:]}
	return cp
}

func (inp input) copy() input {
	cp := make(input, len(inp))
	copy(cp, inp)
	return cp
}

type inputLine struct {
	s           string
	commasPlus1 int
}

func (l *inputLine) len() int {
	return len(l.s)
}

func (l *inputLine) commas() int {
	if l.commasPlus1 == 0 {
		l.commasPlus1 = strings.Count(l.s, ",") + 1
	}
	return l.commasPlus1 - 1
}

func (l *inputLine) String() string {
	return l.s
}

func (l *inputLine) isValidESC() bool {
	return l.commas() > 0
}

func isValidESC(lines input) bool {
	for _, line := range lines {
		if !line.isValidESC() {
			return false
		}
	}
	return true
}

func makeValidESC(lines input, changes int) int {
	for li, line := range lines {
		if !line.isValidESC() {
			if !lines.isLOLMAOList() {
				if lolmaoList, ok := lines.lolmaoLiteralToLolmaoList(); !ok {
					return IMPOSSIBLE // we can't make it a list and we need a comma
				} else {
					return makeValidESC(lolmaoList, changes+2)
				}
			}

			if li < len(lines)-1 {
				return makeValidESC(join(lines, li), changes+1)
			} else if li > 0 {
				return makeValidESC(join(lines, li-1), changes+1)
			}
			// this is a single fucking line, but at least it's a list.
			// lets try to find a place where to put a comma
			for charPos, c := range line.s {
				if c != '[' && c != ']' && c != ',' {
					return makeValidESC(lines.replace(li, charPos, ","), changes+1)
				}
			}

			// last chance, find a an open-close on a single line
			lastChar := '^'
			for charPos, c := range line.s[1 : len(line.s)-1] {
				if lastChar == '[' && c == ']' {
					return makeValidESC(lines.replace(li, charPos, ",").replace(li, charPos+1, ","), changes+2)
				}
				lastChar = c
			}

			return IMPOSSIBLE
		}
	}
	return changes
}

func makeValidLolmaoList(in []byte) int {
	last := len(in) - 1
	c := 0
	if len(in) < 3 {
		// we know it has a comma since it's an ESC
		return IMPOSSIBLE
	}

	p(in, 0, c, 0, "our start")

	open := 1
	if in[0] != '[' {
		in[0] = '['
		c++
		p(in, 0, c, open, "adding an opening bracket")
	}

	if in[last] != ']' {
		in[last] = ']'
		c++
		p(in, last, c, open, "adding a closing bracket")
	}
	// clean cache
	cache = map[cacheKey]int{}
	best = math.MaxInt32
	alreadyChanged := make([]bool, len(in))
	return makeValidLolmaoListRecursive(in, alreadyChanged, open, c, 0)
}

type cacheKey struct {
	i0, open, changes                      int
	canStartHere, canStartNext, isValidESC bool
}

var cache = map[cacheKey]int{}
var best = math.MaxInt32

func makeValidLolmaoListReplacing(in []byte, alreadyChanged []bool, open, c, i0 int, newChar byte, comment string) (result int) {
	in = copyBytes(in)
	alreadyChanged = copyBools(alreadyChanged)

	in[i0] = newChar

	if newChar == '[' {
		open++
	} else if newChar == ']' {
		open--
	}

	if hasChanged(alreadyChanged, i0) {
		c++
	}
	p(in, i0, c, open, comment)
	if c >= best {
		verbose("Cutting branch ")
		return best
	}
	ckey := cacheKey{i0, open, c, canStartElement(in, i0), canStartElement(in, i0+1), isValidESC(bytesToInput(in))}
	if cached, ok := cache[ckey]; ok {
		if verboseEnabled {
			p(in, i0, c, open, fmt.Sprintf("%s returning cached %d for : %+v", comment, cached, ckey))
		}
		return cached
	}
	defer func() {
		if verboseEnabled {
			if _, ok := cache[ckey]; !ok {
				verbose("Caching state %+v as %d", ckey, result)
			}
		}
		cache[ckey] = result
		if result < best {
			best = result
		}
	}()
	return makeValidLolmaoListRecursive(in, alreadyChanged, open, c, i0)
}

func hasChanged(alreadyChanged []bool, i int) bool {
	if !alreadyChanged[i] {
		alreadyChanged[i] = true
		return true
	}
	return false
}
func makeValidLolmaoListKeeping(in []byte, alreadyChanged []bool, open, c, i0 int) (result int) {
	if in[i0] == ']' {
		open--
	} else if in[i0] == '[' {
		open++
	}
	ckey := cacheKey{i0, open, c, canStartElement(in, i0), canStartElement(in, i0+1), isValidESCBytes(in)}
	if cached, ok := cache[ckey]; ok {
		if verboseEnabled {
			p(in, i0, c, open, fmt.Sprintf("keeping returning cached %d for : %+v", cached, ckey))
		}
		return cached
	}
	p(in, i0, c, open, "keeping")
	defer func() {
		if verboseEnabled {
			if _, ok := cache[ckey]; !ok {
				verbose("Caching state %+v as %d", ckey, result)
			}
		}
		cache[ckey] = result
		if result < best {
			best = result
		}
	}()
	return makeValidLolmaoListRecursive(copyBytes(in), copyBools(alreadyChanged), open, c, i0)
}

func makeValidLolmaoListRecursive(in []byte, alreadyChanged []bool, open, c, i0 int) int {
	last := len(in) - 1
	i := i0 + 1
	if i < len(in)-1 {
		p(in, i, c, open, "deciding")

		score := best
		if canKeep(in, i, open) {
			// the cheapest always goes first, if it's feasible to not to change, it's great
			score = makeValidLolmaoListKeeping(in, alreadyChanged, open, c, i)
		}
		if in[i] != ',' {
			score = makeValidLolmaoListReplacing(in, alreadyChanged, open, c, i, ',', "changing to a comma")
		}
		if in[i] != '[' && canStartElement(in, i) {
			score = makeValidLolmaoListReplacing(in, alreadyChanged, open, c, i, '[', "opening bracket")
		}
		if in[i] != ']' && open > 1 {
			score = makeValidLolmaoListReplacing(in, alreadyChanged, open, c, i, ']', "closing bracket")
		}
		return score
	}

	p(in, last, c, open, "we're at the end")

	// try to change opens followed by a comma, those are double points!
	i = last - 1
	closable := open
	for open > 2 && closable > 2 && i > 0 {
		if in[i] != ']' {
			if in[i] == '[' {
				if in[i+1] == ',' || in[i+1] == ']' {
					open -= 2
					closable -= 2
					if hasChanged(alreadyChanged, i) {
						c++
						if c >= best {
							verbose("Cutting branch ")
							return c
						}
					}
					in[i] = ']'
					p(in, i, c, open, fmt.Sprintf("still %d open, removing an opening", open+2))
				} else {
					closable--
				}
			}
		}
		i--
	}

	i = last - 1
	for open > 1 {
		if in[i] != ']' {
			if in[i] == '[' {
				open--
			}
			if hasChanged(alreadyChanged, i) {
				c++
				if c >= best {
					verbose("Cutting branch ")
					return c
				}
			}
			if open == 1 {
				in[i] = ','
				p(in, i, c, open, "putting a comma since open=1 now")
			} else {
				in[i] = ']'
				open--
				p(in, i, c, open, fmt.Sprintf("still %d open, so putting a closing bracket here", open+2))
			}
		}
		i--
	}

	p(in, i, c, open, "opens matched")

	if c == 4 {
		debug("THE SOLUTION (pending fix) ================================================================")
		p(in, 0, c, open, "this is the solution")
	}

	inp := bytesToInput(in)
	if !isValidESC(inp) {
		fixChanges := makeValidESC(inp, c)
		if fixChanges == IMPOSSIBLE {
			return math.MaxInt32
		}
		return fixChanges
	}

	return c
}

func bytesToInput(bytes []byte) input {
	var inp input
	for _, s := range strings.Split(string(bytes), "\n") {
		inp = append(inp, &inputLine{s: s})
	}
	return inp
}

func copyBytes(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)
	return out
}

func copyBools(in []bool) []bool {
	out := make([]bool, len(in))
	copy(out, in)
	return out
}

func isValidESCBytes(in []byte) bool {
	hasComma := false
	for _, n := range in {
		if n == ',' {
			hasComma = true
		}
		if n == '\n' {
			if !hasComma {
				return false
			}
			hasComma = false
		}
	}
	return hasComma
}

func canKeep(in []byte, i, open int) bool {
	if literal(in[i]) {
		return canStartElement(in, i) || literal(in[i-1])
	} else if in[i] == '[' {
		return canStartElement(in, i)
	} else if in[i] == ']' {
		return open > 1
	} else if in[i] == ',' {
		return true
	}
	panic("what else?")
}

func canStartElement(in []byte, i int) bool {
	return in[i-1] == '[' || in[i-1] == ','
}

func literal(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '\n'
}

func join(lines []*inputLine, i int) []*inputLine {
	joined := make([]*inputLine, len(lines)-1)
	copy(joined[:i], lines[:i])
	joined[i] = &inputLine{s: lines[i].s + "," + lines[i+1].s}
	copy(joined[i+1:], lines[i+2:])
	return joined
}

func debug(msg string, args ...interface{}) {
	if debugEnabled {
		fmt.Fprintf(os.Stderr, msg+"\n", args...)
	}
}

func p(in []byte, i, changes, open int, comment string) {
	if verboseEnabled {
		withoutNewline := strings.Replace(string(in), "\n", "N", -1)
		verbose("CHANGES: %d OPEN: %d COMMENT: %s", changes, open, comment)
		verbose("STATE:   %s", withoutNewline)
		verbose("i=%3d:   %s^", i, strings.Repeat(" ", i))
		verbose("")
	}
}

func verbose(msg string, args ...interface{}) {
	if verboseEnabled {
		fmt.Fprintf(os.Stderr, msg+"\n", args...)
	}
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

func min(nums ...int) int {
	m := math.MaxInt32
	for _, n := range nums {
		if n < m {
			m = n
		}
	}
	return m
}
