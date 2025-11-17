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

type UniqueError struct {
	Message string
	Count   int
	Indices []int
}

// For this case the ds used is bitmap.
// byte array of length 536 870 912 bytes.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}
	filename := os.Args[1]
	startTime := time.Now()
	fmt.Printf("Hi\n")
	var wg sync.WaitGroup
	// Setting number of workers to be CPU cores number for optimal solution
	workers := runtime.NumCPU()

	//filename := "ip_addresses.txt"
	lines := make(chan string, 1000)

	var mu sync.Mutex
	//Each number in IPv4 is a byte, meaning whole address is = 32 (4*8)
	//To accumulate for bit state let's assume that each state is 2,
	//then 2^32 = 4 294 967 296 <= number of possible ip addresses
	//that means that 4 294 967 296 / 8 = 536 870 912 bytes is needed for all
	bitmap := make([]byte, 536870912)
	var count uint64
	wg.Add(workers)

	for i := 1; i <= workers; i++ {
		go processor(lines, &wg, bitmap, &count, &mu)
	}
	go reader(filename, lines)
	wg.Wait()
	//count := counter(bitmap)
	finishTime := time.Since(startTime)
	fmt.Printf("Done, unique number is %v, time taken %v \n", count, finishTime)

}

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

func processor(lines <-chan string, wg *sync.WaitGroup, bitmap []byte, count *uint64, mu *sync.Mutex) {

	defer wg.Done()
	startTime := time.Now()

	for ipAddressString := range lines {
		ipAddress := stringToUint32(ipAddressString)
		// if error != nil {
		// 	// key := error.Error()
		// 	// if uniqueError, exists := errorMap[key]; exists {
		// 	// 	uniqueError.Count++
		// 	// } else {
		// 	// 	errorMap[key] = &UniqueError{
		// 	// 		Message: key,
		// 	// 		Count:   1,
		// 	// 	}
		// 	// }
		// 	// continue
		// }
		//Byte index to check in the bitmap:
		byteIndex := ipAddress / 8
		byteOffset := ipAddress % 8
		//Bit to check:
		mask := byte(1 << byteOffset)
		//println("mask ", mask)
		//If item doesn't exist in the bitmap, add there and increase count
		mu.Lock()
		if bitmap[byteIndex]&mask == 0 {
			bitmap[byteIndex] |= mask
			atomic.AddUint64(count, 1)
		}
		mu.Unlock()
		//fmt.Println("Line: ", ipAddress)
	}
	finishTime := time.Since(startTime)
	fmt.Printf("Time taken for processing %v \n", finishTime)

}

func stringToUint32(ipString string) uint32 {
	ipInt := net.ParseIP(ipString).To4()
	if ipInt == nil {
		//return 0, errors.New("invalid IPv4 address")
		panic("Invalid ipV4 address")
	}
	return binary.BigEndian.Uint32(ipInt)
}
