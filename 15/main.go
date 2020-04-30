package main

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/vimeo/go-util/crc32combine"
)

const DEBUG = false

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var tcs []testCase
	jobs := make(chan testCase, 1000)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go func(i int) {
			fmt.Fprintf(os.Stderr, "Worker %d\n", i)
			for tc := range jobs {
				tc.res <- solve(tc.name, tc.additions)
			}
		}(i)
	}

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		name := line[0]
		var additions []add
		n := atoi(line[1])
		for i := 0; i < n; i++ {
			scanner.Scan()
			line := strings.Split(scanner.Text(), " ")
			additions = append(additions, add{
				pos:   int64(atoi(line[0])),
				byte:  byte(atoi(line[1])),
				order: i + 1,
			})
		}
		tc := testCase{name: name, additions: additions, res: make(chan []string, 1)}
		jobs <- tc
		tcs = append(tcs, tc)
	}
	close(jobs)

	for _, tc := range tcs {
		for _, l := range <-tc.res {
			fmt.Println(l)
		}
	}
}

type testCase struct {
	name      string
	additions []add
	res       chan []string
}

var total int64

func solve(name string, additions []add) []string {
	chunks := calculateChunks(name, additions)
	for i, c := range chunks {
		debug("Chunk %d %+v", i, c)
	}

	res := make([]string, len(additions)+1)
	for i := 0; i <= len(additions); i++ {
		debug("")
		debug("SUMADDITIONS %d FOLLOWS", i)
		res[i] = fmt.Sprintf("%s %d: %08x", name, i, sumAdditions(chunks, additions[:i]))
	}
	fmt.Fprintf(os.Stderr, "Solved %s total %d\n", name, atomic.AddInt64(&total, 1))
	return res
}

func sumAdditions(chunks []chunk, adds []add) uint32 {

	sort.Slice(adds, func(i, j int) bool {
		if adds[i].pos == adds[j].pos {
			return adds[i].order > adds[j].order
		}
		return adds[i].pos < adds[j].pos
	})

	var (
		sum    uint32
		offset int64
		chunki int
	)

	for _, a := range adds {
		debug("Proceeding to %+v", a)
		for ; offset < a.pos && chunki < len(chunks); chunki++ {
			c := chunks[chunki]
			debug("Current offset %d of chunk %d: %+v is less than addition %+v pos %d, advancing", offset, chunki, c, a, a.pos)
			sum = crc32combine.CRC32Combine(crc32.IEEE, sum, c.crc32, c.len)
			offset += c.len + 1
			debug("Offset is now %d with partial sum %08x", offset, sum)
		}
		byteSum := crc32.ChecksumIEEE([]byte{a.byte})
		debug("Adding %+v", a)
		sum = crc32combine.CRC32Combine(crc32.IEEE, sum, byteSum, 1)
		debug("Partial sum %08x", sum)
	}
	for ; chunki < len(chunks); chunki++ {
		c := chunks[chunki]
		debug("completing chunk %d %+v", chunki, c)
		sum = crc32combine.CRC32Combine(crc32.IEEE, sum, c.crc32, c.len)
	}
	debug("Sum %08x", sum)
	debug("===")
	return sum
}

type chunk struct {
	offset int64
	len    int64
	crc32  uint32
}

type add struct {
	pos   int64
	byte  byte
	order int
}

func calculateChunks(name string, additions []add) []chunk {
	f, err := os.Open("./animals/" + name)
	assertNoError(err)
	defer f.Close()

	offsets := calculateOffsets(additions)
	chunks := make([]chunk, len(offsets))

	for i, o := range offsets {
		hasher := crc32.NewIEEE()
		if i < len(offsets)-1 {
			// there's a next one
			n, err := io.CopyN(hasher, f, offsets[i+1]-o)
			if err != io.EOF {
				assertNoError(err)
			} else {
				debug("got EOF at CopyN offset %d", o)
			}
			chunks[i] = chunk{
				offset: o,
				len:    n,
				crc32:  hasher.Sum32(),
			}
		} else {
			// last
			n, err := io.Copy(hasher, f)
			if err != io.EOF {
				assertNoError(err)
			} else {
				debug("got EOF at Copy offset %d", o)
			}
			chunks[i] = chunk{
				offset: o,
				len:    n,
				crc32:  hasher.Sum32(),
			}
		}
	}
	return chunks
}

func calculateOffsets(additions []add) []int64 {
	cuts := map[int64]bool{}
	var offs []int64
	zeroCut := false
	for _, a := range additions {
		cut := a.pos
		for prev := range cuts {
			if prev < cut {
				cut--
			}
		}
		cuts[cut] = true
		if cut == 0 {
			zeroCut = true
		}
		offs = append(offs, cut)
	}
	cuts[0] = true
	if !zeroCut {
		offs = append(offs, 0)
	}
	sort.Slice(offs, func(i, j int) bool { return offs[i] < offs[j] })
	return offs
}

func debug(msg string, args ...interface{}) {
	if DEBUG {
		fmt.Printf(msg+"\n", args...)
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
