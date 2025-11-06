package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func ByteToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func GetMemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("[HeapAlloc=%vMiB;TotalAlloc=%vMiB;Sys=%vMiB]",
		ByteToMb(m.HeapAlloc),
		ByteToMb(m.TotalAlloc),
		ByteToMb(m.Sys))
}

func GetMemUsageMiB() (heapAlloc uint64, totalAlloc uint64, sys uint64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return ByteToMb(m.HeapAlloc), ByteToMb(m.TotalAlloc), ByteToMb(m.Sys)
}

// GetMemUsageOSMiB returns the RSS (Resident Set Size) memory usage in MiB
// as reported by the operating system.
// This value should match what Activity Monitor (macOS) or system monitors show.
// Returns 0 if unable to get the memory usage or if running on Windows.
func GetMemUsageOSMiB() uint64 {
	// Windows is not supported - return 0
	if runtime.GOOS == "windows" {
		return 0
	}

	pid := os.Getpid()

	// Use ps command to get RSS in KB
	// -o rss= outputs only RSS value without headers
	// Works on both macOS and Linux
	cmd := exec.Command("ps", "-o", "rss=", "-p", strconv.Itoa(pid))
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	// Parse the output (RSS in KB)
	rssStr := strings.TrimSpace(string(output))
	rssKB, err := strconv.ParseUint(rssStr, 10, 64)
	if err != nil {
		return 0
	}

	// Convert KB to MiB
	return rssKB / 1024
}
