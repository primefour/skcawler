package downloader

import (
	"testing"
	"time"
)

func Benchmark(b *testing.B) {
	customTimerTag := true
	if customTimerTag {
		b.StopTimer()
	}

	time.Sleep(time.Second)

	if customTimerTag {
		b.StartTimer()
	}
}
