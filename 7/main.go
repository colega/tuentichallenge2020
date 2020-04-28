package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var (
	// http://ergoemacs.org/misc/qwerty_dvorak_table.html
	q2d = map[byte]byte{
		'\'': '-',
		',':  'w',
		'-':  '[',
		'.':  'v',
		'/':  'z',
		';':  's',
		'=':  ']',
		'[':  '/',
		']':  '=',
		'a':  'a',
		'b':  'x',
		'c':  'j',
		'd':  'e',
		'e':  '.',
		'f':  'u',
		'g':  'i',
		'h':  'd',
		'i':  'c',
		'j':  'h',
		'k':  't',
		'l':  'n',
		'm':  'm',
		'n':  'b',
		'o':  'r',
		'p':  'l',
		'q':  '\'',
		'r':  'p',
		's':  'o',
		't':  'y',
		'u':  'g',
		'v':  'k',
		'w':  ',',
		'x':  'q',
		'y':  'f',
		'z':  ';',
	}

	d2q = map[byte]byte{}
)

func init() {
	for q, d := range q2d {
		d2q[d] = q
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	testCases := scanInt(scanner)

	for testCase := int64(1); testCase <= testCases; testCase++ {
		mustScan(scanner)
		line := scanner.Text()
		out := ""
		for _, dvorak := range line {
			if qwerty, ok := d2q[byte(dvorak)]; ok {
				out += string(qwerty)
			} else {
				out += string(dvorak)
			}
		}

		fmt.Printf("Case #%d: %s\n", testCase, out)
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
