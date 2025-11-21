package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

// Function that is converting ip string value to Uint32
func stringToUint32(ipString string) (uint32, error) {
	ipInt := net.ParseIP(ipString).To4()
	if ipInt == nil {
		return 0, fmt.Errorf("%v is not a valid IPv4 address", ipString)
	}
	// fmt.Println("Address is valid")
	return binary.BigEndian.Uint32(ipInt), nil
}

func ticker(processed *int64, totalFileSize int64, tickerDone chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	start := time.Now()
	for {
		select {
		case <-ticker.C:
			p := atomic.LoadInt64(processed)
			elapced := time.Since(start).Seconds()
			percent := float64(p) / float64(totalFileSize) * 100
			speed := float64(p) / (1024 * 1024) / elapced
			fmt.Printf("\rProgress: %.2f%% (%d /%d MB), %.1f MB/s",
				percent,
				p/1024/1024,
				totalFileSize/1024/1024,
				speed,
			)
		case <-tickerDone:
			return
		}
	}
}
