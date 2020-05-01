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

var debugEnabled = true
var verboseEnabled = false

const IMPOSSIBLE = -1
const UNKNOWN = -2

var onlyTestCase = flag.Int("only", 0, "Run only one test case. 0 runs all")

func main() {
	flag.Parse()

	//http.ListenAndServe(":8080", pprof.)

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

		if *onlyTestCase > 0 && testCase != *onlyTestCase {
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
	alreadyChanged := make([]bool, len(in))
	return makeValidLolmaoListRecursive(in, alreadyChanged, math.MaxInt32, open, c, 0)
}

func makeValidLolmaoListReplacing(in []byte, alreadyChanged []bool, maxChanges, open, c, i0 int, newChar byte, comment string) int {
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
	if c >= maxChanges {
		verbose("Cutting branch ")
		return c
	}
	return makeValidLolmaoListRecursive(in, alreadyChanged, maxChanges, open, c, i0)
}

func hasChanged(alreadyChanged []bool, i int) bool {
	if !alreadyChanged[i] {
		alreadyChanged[i] = true
		return true
	}
	return false
}

func makeValidLolmaoListKeeping(in []byte, alreadyChanged []bool, maxChanges, open, c, i0 int) int {
	if in[i0] == ']' {
		open--
	} else if in[i0] == '[' {
		open++
	}
	return makeValidLolmaoListRecursive(copyBytes(in), copyBools(alreadyChanged), maxChanges, open, c, i0)
}

func makeValidLolmaoListRecursive(in []byte, alreadyChanged []bool, maxChanges, open, c, i0 int) int {
	last := len(in) - 1

	p(in, i0, c, open, "branching")
	for i := i0 + 1; i < len(in)-1; i++ {
		p(in, i, c, open, "loop")
		if in[i] == ',' {
			if canPutABracketHere(in, i) && bracketHereWouldHaveAClosingMatch(in, open, i) {
				withBracket := makeValidLolmaoListReplacing(in, alreadyChanged, maxChanges, open, c, i, '[', "opening bracket would be closed, so we have to change less stuff later")
				withoutBracket := makeValidLolmaoListKeeping(in, alreadyChanged, withBracket, open, c, i)
				return min(withBracket, withoutBracket)
			}
		} else if literal(in[i]) {
			if canPutABracketHere(in, i) && bracketHereWouldHaveAClosingMatch(in, open, i) {
				withBracket := makeValidLolmaoListReplacing(in, alreadyChanged, maxChanges, open, c, i, '[', "opening bracket would be closed, so we have to change less stuff later (literal)")
				withoutBracket := makeValidLolmaoListKeeping(in, alreadyChanged, withBracket, open, c, i)
				return min(withBracket, withoutBracket)
			} else if in[i-1] != ',' && in[i-1] != '[' {
				keepingLiteralIfPossible := maxChanges
				if literal(in[i-1]) {
					keepingLiteralIfPossible = makeValidLolmaoListKeeping(in, alreadyChanged, maxChanges, open, c, i)
				}
				if open > 1 && !closingHereWouldRequireChangingNext(in, i) {
					closingIt := makeValidLolmaoListReplacing(in, alreadyChanged, keepingLiteralIfPossible, open, c, i, ']', "we have enough open lets close now we can")
					replacingByAComma := makeValidLolmaoListReplacing(in, alreadyChanged, closingIt, open, c, i, ',', "closing here would require changing next literal, so lets avoid changing it")
					return min(closingIt, replacingByAComma, keepingLiteralIfPossible)
				} else {
					replacingByAComma := makeValidLolmaoListReplacing(in, alreadyChanged, keepingLiteralIfPossible, open, c, i, ',', "can't have an opening here and closing would require changing next or we don't have enough open")
					return min(keepingLiteralIfPossible, replacingByAComma)
				}
			} else if canPutABracketHere(in, i) {
				// like first case but again
				withBracket := makeValidLolmaoListReplacing(in, alreadyChanged, maxChanges, open, c, i, '[', "opening bracket would be closed, so we have to change less stuff later (literal)")
				withoutBracket := makeValidLolmaoListKeeping(in, alreadyChanged, withBracket, open, c, i)
				return min(withBracket, withoutBracket)
			} else {
				for literal(in[i+1]) {
					i++
				}
			}
		} else if in[i] == '[' {
			if canPutABracketHere(in, i) {
				// we can leave this opening
				open++
			} else {
				if open > 1 && !closingHereWouldRequireChangingNext(in, i) {
					closingIt := makeValidLolmaoListReplacing(in, alreadyChanged, maxChanges, open, c, i, ']', "we have enough open lets close now we can")
					replacingByAComma := makeValidLolmaoListReplacing(in, alreadyChanged, closingIt, open, c, i, ',', "closing here would require changing next literal, so lets avoid changing it")
					return min(closingIt, replacingByAComma)
				} else {
					return makeValidLolmaoListReplacing(in, alreadyChanged, maxChanges, open, c, i, ',', "can't have an opening here and closing would require changing next or we don't have enough open")
				}
			}
		} else if in[i] == ']' {
			if open > 1 {
				if closingHereWouldRequireChangingNext(in, i) {
					replacingByAComma := makeValidLolmaoListReplacing(in, alreadyChanged, maxChanges, open, c, i, ',', "closing here would require changing next literal, so lets avoid changing it")
					keepingClosingBracket := makeValidLolmaoListKeeping(in, alreadyChanged, replacingByAComma, open, c, i)
					return min(replacingByAComma, keepingClosingBracket)
				} else {
					open--
				}
				// just close igt
			} else {
				// we have to change it, maybe to something useful?
				if bracketHereWouldHaveAClosingMatch(in, open, i) && canPutABracketHere(in, i) {
					return makeValidLolmaoListReplacing(in, alreadyChanged, maxChanges, open, c, i, '[', "opening bracket would be closed, not sure if its a good idea")
				} else {
					return makeValidLolmaoListReplacing(in, alreadyChanged, maxChanges, open, c, i, ',', "a default comma is okay here")
				}
			}
		} else {
			panic(fmt.Errorf("unexpected char %c", in[i]))
		}
	}

	p(in, last, c, open, "loop done")

	// try to change opens followed by a comma, those are double points!
	i := last - 1
	closable := open
	for open > 2 && closable > 2 && i > 0 {
		if in[i] != ']' {
			if in[i] == '[' {
				if in[i+1] == ',' || in[i+1] == ']' {
					open -= 2
					closable -= 2
					if hasChanged(alreadyChanged, i) {
						c++
						if c >= maxChanges {
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
				if c >= maxChanges {
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

	inp := bytesToInput(in)
	if !isValidESC(inp) {
		fixChanges := makeValidESC(inp, c)
		if fixChanges == IMPOSSIBLE {
			return maxChanges
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

func placeForACheaperClosingBracket(in []byte, i int) (int, bool) {
	i++
	tmp := i
	// first try to find one where we can also remove an opening bracket
	for ; i < len(in)-1; i++ {
		if in[i] == '[' && (in[i+1] == ',' || in[i+1] == ']') {
			return i, true
		}
	}

	i = tmp
	// no luck, then anything is okay
	for ; i < len(in)-1; i++ {
		if in[i+1] == ',' || in[i+1] == ']' {
			return i, true
		}
	}
	return 0, false
}

func placeForACheaperOpeningBracketToReplaceByComma(in []byte, open, i int) (int, bool) {
	i++
	open++
	for ; i < len(in)-1; i++ {
		if in[i] == '[' {
			if (in[i-1] == ',' || in[i-1] == '[') && open > 2 {
				return i, true
			}
			open++
		} else if in[i] == ']' {
			open--
			if open == 1 {
				// next opening brackets dont affect this one
				return 0, false
			}
		}
	}
	return 0, false
}

func closingHereWouldRequireChangingNext(in []byte, i int) bool {
	return in[i+1] != ',' && in[i+1] != ']'
}
func canPutABracketHere(in []byte, i int) bool {
	return in[i-1] == '[' || in[i-1] == ','
}

func bracketHereWouldHaveAClosingMatch(in []byte, open, i int) bool {
	return !bracketHereWouldBeUnclosed(in, open, i)
}

// would it be unclosed checks whether the bracket would be unclosed if we place one here
func bracketHereWouldBeUnclosed(in []byte, open, i int) bool {
	open++
	i++
	for ; i < len(in)-1; i++ {
		if in[i] == '[' {
			open++
		} else if in[i] == ']' /*&& (in[i+1] == ',' || in[i+1] == ']')*/ {
			open--
		}
		if open == 1 {
			// this is not the last closing brackets and we went to zero
			// so so it's definitely not unclosed
			return false
		}
	}
	return open > 1
}

func bracketHereWouldBeUnclosedWithOneExtra(in []byte, open, i int) bool {
	open++
	i++
	for ; i < len(in)-1; i++ {
		if in[i] == '[' {
			open++
		} else if in[i] == ']' && (in[i+1] == ',' || in[i+1] == ']') /* checking if it will stay there */ {
			open--
			if open == 1 {
				// this is not the last closing brackets and we went to zero
				// so so it's definitely not unclosed
				return false
			}
		}
	}
	return open > 2
}

func excessOpen(in []byte, open, i int) bool {
	for ; i < len(in)-1; i++ {
		if in[i] == '[' {
			open++
		} else if in[i] == ']' {
			open--
			if open == 0 {
				// this is not the last closing brackets and we went to zero, not good
				return false
			}
		}
	}
	if i == len(in) {
		return true
	}
	return false
}

func count(in []byte, c byte) int {
	count := 0
	for _, b := range in {
		if c == b {
			count++
		}
	}
	return count
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
