package info_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
)

// JSON vs Sprinf

func BenchmarkJsonMarshal(b *testing.B) {
	xs := make([]int, 1000000)
	for i := 0; i < len(xs); i++ {
		xs[i] = rand.Int()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(xs)
	}
}

func BenchmarkFmtSprint(b *testing.B) {
	xs := make([]int, 1000000)
	for i := 0; i < len(xs); i++ {
		xs[i] = rand.Int()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", xs)
	}
}

/*
➜  info git:(main) ✗ go test -bench .
goos: darwin
goarch: arm64
pkg: github.com/filevich/truco-ai/info
BenchmarkJsonMarshal-8   	      30	  40122399 ns/op
BenchmarkFmtSprint-8     	      10	 104271754 ns/op
PASS
ok  	github.com/filevich/truco-ai/info	3.824s
*/
