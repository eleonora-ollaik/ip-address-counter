package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// For this case the ds used is bitmap.
// byte array of length 536 870 912 bytes.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Please run again. Usage: go run main.go <filename>")
		return
	}
	startTime := time.Now()
	filename := os.Args[1]

	// fileExtention := strings.ToLower(filepath.Ext(filename))
	// if fileExtention != ".txt" {
	// 	fmt.Printf("Error: only .txt files are allowed\n")
	// 	return
	// }

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

	// 1000 is an arbitrary number, enough for the usecase
	lines := make(chan string, 1000)

	for i := 1; i <= workers; i++ {
		go processor(lines, &wg, bitmap, &ipCount, &linesCount, &muUniqueIps)
	}
	go reader(filename, lines)
	wg.Wait()

	finishTime := time.Since(startTime)
	fmt.Printf("Execution complete.\n")
	fmt.Println()
	fmt.Printf("Unique number of IPs: %v\n", ipCount)
	fmt.Printf("Lines processed: %v\n", linesCount)
	fmt.Printf("Time taken to execute: %v \n", finishTime)
	fmt.Println()
}

// Reader that goes through file and feeds each line into the channel
func reader(filename string, lines chan<- string) {
	startTime := time.Now()
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines <- scanner.Text()
	}
	close(lines)
	finishTime := time.Since(startTime)
	fmt.Printf("Time taken for reading %v \n", finishTime)
}

// Processor function that is defining unique ip addresses from the stream
func processor(lines <-chan string, wg *sync.WaitGroup, bitmap []byte, ipCount *uint64, linesCount *uint64, muUniqueIps *sync.Mutex) {
	defer wg.Done()
	startTime := time.Now()

	for ipAddressString := range lines {
		ipAddress, err := stringToUint32(ipAddressString)
		atomic.AddUint64(linesCount, 1)
		if err != nil {
			return // Skip the non ip line to keep counting IPv4s
		}
		//Byte index to check in the bitmap:
		byteIndex := ipAddress / 8
		byteOffset := ipAddress % 8
		//Bit to check:
		mask := byte(1 << byteOffset)

		//If item doesn't exist in the bitmap, add there and increase count
		muUniqueIps.Lock()
		if bitmap[byteIndex]&mask == 0 {
			bitmap[byteIndex] |= mask
			atomic.AddUint64(ipCount, 1)
		}
		muUniqueIps.Unlock()

	}
	finishTime := time.Since(startTime)
	fmt.Printf("Time taken for processing %v \n", finishTime)

}

// Function that is converting ip string value to Uint32
func stringToUint32(ipString string) (uint32, error) {
	ipInt := net.ParseIP(ipString).To4()
	if ipInt == nil {
		return 0, fmt.Errorf("is not a valid IPv4 address")
	}
	return binary.BigEndian.Uint32(ipInt), nil
}
