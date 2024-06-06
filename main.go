// Description: Ayaka is mem test tool.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"
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
	ret := make([]byte, size)

	// random write
	for i := 0; i < len(ret); i += 8 {
		r := rand.Uint64()
		b := *(*[8]byte)(unsafe.Pointer(&r))

		for j := 0; j < 8; j++ {
			ret[i+j] = b[j]
		}

	}

	return ret
}

func main() {
	blockSizeStr := flag.String("b", "100MB", "block size")
	memLimitStr := flag.String("m", "1GB", "max memory")
	// loopの待ち時間
	// waitMilliSec ごとにループを回す
	waitMilliSec := flag.Int("w", 0, "wait ms after each allocation")
	// 完了後に待機する時間
	nopTime := flag.Int("n", 0, "wait ms after all allocation")

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

	if *nopTime == 0 {
		return
	}
	fmt.Println("Nop wait...", time.Duration(*nopTime)*time.Second)
	time.Sleep(time.Duration(*nopTime) * time.Second)

	_ = buf
}
