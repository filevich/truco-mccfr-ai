package utils

import (
	"fmt"
	"runtime"
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
