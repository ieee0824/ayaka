// Description: Ayaka is mem test tool.
package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	_       = iota
	KB size = 1 << (10 * iota)
	MB
	GB
)

// parseSize parses a human-readable size string like "1GB" and returns the
func parseSize(s string) (size, error) {
	s = strings.Replace(s, " ", "", -1)

	if b, err := strconv.ParseUint(s, 10, 64); err == nil {
		return size(b), nil
	}

	// example 1GB
	// tokens[0], tokens[1] = 1, GB
	tokens := [2]string{
		s[:len(s)-2],
		s[len(s)-2:],
	}

	b, err := strconv.ParseUint(tokens[0], 10, 64)
	if err != nil {
		return 0, err
	}

	switch tokens[1] {
	case "KB":
		return size(b) * KB, nil
	case "MB":
		return size(b) * MB, nil
	case "GB":
		return size(b) * GB, nil
	default:
		return 0, fmt.Errorf("invalid size suffix: %s", tokens[1])
	}
}

type size uint64

func (s size) String() string {
	switch {
	case s >= GB:
		return fmt.Sprintf("%.2fGB", float64(s)/float64(GB))
	case s >= MB:
		return fmt.Sprintf("%.2fMB", float64(s)/float64(MB))
	case s >= KB:
		return fmt.Sprintf("%.2fKB", float64(s)/float64(KB))
	default:
		return fmt.Sprintf("%dB", s)
	}
}

func alloc(size size) []byte {
	return make([]byte, size)
}

func main() {
	blockSizeStr := flag.String("b", "100MB", "block size")
	memLimitStr := flag.String("m", "1GB", "max memory")
	waitMilliSec := flag.Int("w", 0, "wait ms after each allocation")
	flag.Parse()

	blockSize, err := parseSize(*blockSizeStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	memLimit, err := parseSize(*memLimitStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := []byte{}

	loopNum := memLimit / blockSize

	for i := size(0); i < loopNum; i++ {
		buf = append(buf, alloc(blockSize)...)
		fmt.Printf("Allocated %s, wait %s...\n", blockSize+blockSize*i, time.Duration(*waitMilliSec)*time.Millisecond)
		time.Sleep(time.Duration(*waitMilliSec) * time.Millisecond)
	}

	_ = buf
}
