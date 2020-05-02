package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

var debugEnabled, verboseEnabled bool

type fileExtent struct {
	DiskBytenr    int64 `json:"disk_bytenr"`
	Offset        int64 `json:"offset"`
	NumBytes      int64 `json:"num_bytes"`
	LogicalOffset int64 `json:"logical_offset"`
}

var buf = make([]byte, 4096)

func main() {
	var jsonFile string
	var filesPath string
	flag.StringVar(&jsonFile, "json", "test.json", "JSON Filename provided by the show_files.py")
	flag.StringVar(&filesPath, "path", "/mnt/test_input/", "Files path ")
	flag.BoolVar(&verboseEnabled, "verbose", false, "Verbose logging")
	flag.BoolVar(&debugEnabled, "debug", false, "Debug logging")
	flag.Parse()

	loaded := map[string][]fileExtent{}

	data, err := ioutil.ReadFile(jsonFile)
	assertNoError(err)
	err = json.Unmarshal(data, &loaded)
	assertNoError(err)

	cowsByContents := map[string][]string{}

	for filename, extents := range loaded {
		cowContents := ""

		sort.Slice(extents, func(i, j int) bool {
			return extents[i].LogicalOffset < extents[j].LogicalOffset
		})

		prefixLen := extents[0].NumBytes
		chunkSum := offsetHash(filesPath+filename, extents[1].LogicalOffset)
		cowContents = fmt.Sprintf("%d zeros + 4096b crc32=%08x", prefixLen, chunkSum)
		/*

			for _, e := range extents {
				if e.NumBytes == 4096 {
					cowContents += fmt.Sprintf("4096b:crc32:%08x...", chunkSum)
				} else {
					cowContents += fmt.Sprintf("%d:%d+%d (%d)...", e.LogicalOffset, e.DiskBytenr, e.Offset, e.NumBytes)
				}
			}
		*/
		cowsByContents[cowContents] = append(cowsByContents[cowContents], filename)
		debug("Cow %s: %s", filename, cowContents)
	}

	var solutions int
	for contents, names := range cowsByContents {
		if len(names) == 1 {
			cow := names[0]
			fmt.Println("Solution: ", cow, contents)
			solutions++
		}
	}

	debug("Total %d cows, %d solutions", len(loaded), solutions)
}

func offsetHash(filename string, offset int64) uint32 {
	file, err := os.Open(filename)
	assertNoError(err)
	defer file.Close()
	ret, err := file.Seek(offset, 0)
	assertNoError(err)
	if ret != offset {
		panic(fmt.Errorf("%s: expected to seek to %d but seeked to %d", filename, offset, ret))
	}
	n, err := file.Read(buf)
	assertNoError(err)
	if n != 4096 {
		panic(fmt.Errorf("%s: expected to read %d bytes, got %d", filename, 4096, n))
	}
	return crc32.ChecksumIEEE(buf)
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
