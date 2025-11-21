package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Please run again. Usage: go run main.go <filename>")
		return
	}
	startTime := time.Now()
	filename := os.Args[1]

	if _, error := os.Stat(filename); os.IsNotExist(error) {
		fmt.Printf("Error: file %v does not exist\n", filename)
		return
	} else if error != nil {
		fmt.Printf("Error: error checking the file: %v\n", error)
	}

	fmt.Printf("Starting the count...\n")
	fmt.Println()

	//Each number in IPv4 is a byte, meaning whole address is = 32 (4*8)
	//To accumulate for bit state let's assume that each state is 2,
	//then 2^32 = 4 294 967 296 <= number of possible ip addresses
	//that means that 4 294 967 296 / 8 = 536 870 912 bytes is needed for all
	bitmap := make([]byte, 536870912)

	var ipCount uint64
	var linesCount uint64
	var wg sync.WaitGroup
	var muUniqueIps sync.Mutex

	// Setting number of workers to be CPU cores number for optimal solution
	workers := runtime.NumCPU()
	wg.Add(workers)

	fileInfo, _ := os.Stat(filename)
	totalFileSize := fileInfo.Size()
	chunkSize := totalFileSize / int64(workers)

	// Progress bar printer
	var processed int64
	tickerDone := make(chan struct{})
	go ticker(&processed, totalFileSize, tickerDone)

	for i := 0; i < workers; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == workers-1 {
			end = totalFileSize
		}
		isFirst := i == 0
		go processor(filename, &wg, start, end, isFirst, &processed, bitmap, &ipCount, &linesCount, &muUniqueIps)
	}

	wg.Wait()
	close(tickerDone)
	finishTime := time.Since(startTime)
	fmt.Println()
	fmt.Printf("Execution complete.\n")
	fmt.Println()
	fmt.Printf("Unique number of IPs: %v\n", ipCount)
	fmt.Printf("Lines processed: %v\n", linesCount)
	fmt.Printf("Time taken to execute: %v \n", finishTime)
	fmt.Println()
}
