package utils

import (
	"fmt"
	"time"
)

var F = "02 Jan 2006 15:04:05"

func CurrentTimeFileFormat() string {
	currentTime := time.Now()
	return fmt.Sprintf("%d-%02d-%02d-%02d%02d",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
	)
}

func MiniCurrentTime() string {
	currentTime := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d",
		currentTime.Year()%100,
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
	)
}

func FmtDuration(d time.Duration) string {
	d = d.Round(time.Second)

	h := d / time.Hour
	d -= h * time.Hour

	m := d / time.Minute
	d -= m * time.Minute

	s := d / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func B2MiB(b uint64) uint64 {
	return b / 1024 / 1024
}
