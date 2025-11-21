package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

// Function that reads lines from file and processes them directly
func processor(filename string, wg *sync.WaitGroup, start, end int64, isFirst bool, processed *int64, bitmap []byte, ipCount *uint64, linesCount *uint64, muUniqueIps *sync.Mutex) {
	defer wg.Done()
	file, _ := os.Open(filename)
	defer file.Close()

	_, _ = file.Seek(start, 0)
	scanner := bufio.NewScanner(file)

	if !isFirst {
		if scanner.Scan() {
			// skip partial line
		}
	}

	position := start

	for scanner.Scan() {
		line := scanner.Bytes()
		position += int64(len(line)) + 1

		if position > end {
			break
		}

		// Process the line
		ipAddressString := string(line)
		ipAddress, err := stringToUint32(ipAddressString)
		atomic.AddUint64(linesCount, 1)
		if err != nil {
			// Log the error but continue processing
			fmt.Println(err)
			continue
		}

		// Byte index to check in the bitmap:
		byteIndex := ipAddress / 8
		byteOffset := ipAddress % 8
		// Bit to check:
		mask := byte(1 << byteOffset)

		// If item doesn't exist in the bitmap, add there and increase count
		muUniqueIps.Lock()
		if bitmap[byteIndex]&mask == 0 {
			bitmap[byteIndex] |= mask
			*ipCount++
		}
		muUniqueIps.Unlock()

		atomic.AddInt64(processed, int64(len(line)+1))
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner error: %v", err)
	}
}
