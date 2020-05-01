package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var debugEnabled = true
var verboseEnabled = false

const IMPOSSIBLE = -1
const UNKNOWN = -2

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	testCases := atoi(scanner.Text())

	for testCase := 1; testCase <= testCases; testCase++ {
		scanner.Scan()
		L := atoi(scanner.Text())
		lines := make([]*inputLine, 0, L)
		for i := 0; i < L; i++ {
			scanner.Scan()
			lines = append(lines, &inputLine{s: scanner.Text()})
		}

		debug("Solving case %d", testCase)
		for i, line := range lines {
			debug("LINE %2d: %s", i, line)
		}
		solution := solve(lines)
		if solution >= 0 {
			fmt.Printf("Case #%d: %d\n", testCase, solution)
		} else if solution == IMPOSSIBLE {
			fmt.Printf("Case #%d: IMPOSSIBLE\n", testCase)
		} else if solution == UNKNOWN {
			fmt.Printf("Case #%d: UNKNOWN\n", testCase)
		}
	}
}

func solve(lines []*inputLine) int {
	debug("Is valid ESC: %t", isValidESC(lines))
	if !isValidESC(lines) {
		return makeValidESC(lines, 0)
	}
	return UNKNOWN
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

func isValidESC(lines []*inputLine) bool {
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
			return IMPOSSIBLE
		}
	}
	return changes
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
