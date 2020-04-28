package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	bytes, err := ioutil.ReadFile("pg17013.txt")
	assertNoError(err)
	s := strings.ToLower(string(bytes))
	tokens := tokenize(s)
	counts, words := sorted(tokens)
	ranks := make(map[string]int, len(words))
	for pos, w := range words {
		ranks[w] = pos + 1
	}

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		panic("can't scan")
	}
	testCases, err := strconv.Atoi(scanner.Text())
	assertNoError(err)

	for tc := 1; tc <= testCases; tc++ {
		if !scanner.Scan() {
			panic("can't scan")
		}
		l := scanner.Text()
		if '0' <= l[0] && l[0] <= '9' {
			r, err := strconv.Atoi(l)
			assertNoError(err)
			w := words[r-1]
			fmt.Printf("Case #%d: %s %d\n", tc, w, counts[w])
		} else {
			fmt.Printf("Case #%d: %d #%d\n", tc, counts[l], ranks[l])
		}
	}
}

func sorted(tokens []string) (map[string]int, []string) {
	counts := make(map[string]int)
	for i := range tokens {
		counts[tokens[i]]++
	}
	words := make([]string, 0, len(counts))
	for w := range counts {
		words = append(words, w)
	}
	sort.Slice(words, func(i, j int) bool {
		wi, wj := words[i], words[j]
		ci, cj := counts[wi], counts[wj]
		if ci != cj {
			return ci > cj
		}
		return wi < wj
	})
	return counts, words
}

func tokenize(s string) []string {
	runes := []rune(s)
	var tokens []string
	j := 0
	for i := 0; i < len(runes); i++ {
		if _, valid := validRunes[runes[i]]; !valid {
			if j+1 < i-1 {
				tokens = append(tokens, string(runes[j:i]))
			}
			j = i + 1
		}
	}
	return tokens
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

var validRunes = map[rune]struct{}{
	'a': {},
	'b': {},
	'c': {},
	'd': {},
	'e': {},
	'f': {},
	'g': {},
	'h': {},
	'i': {},
	'j': {},
	'k': {},
	'l': {},
	'm': {},
	'n': {},
	'ñ': {},
	'o': {},
	'p': {},
	'q': {},
	'r': {},
	's': {},
	't': {},
	'u': {},
	'v': {},
	'w': {},
	'x': {},
	'y': {},
	'z': {},
	'á': {},
	'é': {},
	'í': {},
	'ó': {},
	'ú': {},
	'ü': {},
}
