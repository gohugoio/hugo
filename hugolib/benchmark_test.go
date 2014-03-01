package hugolib

import (
	"bytes"
	"os"
	"testing"
)

func BenchmarkParsePage(b *testing.B) {
	f, _ := os.Open("redis.cn.md")
	sample := new(bytes.Buffer)
	sample.ReadFrom(f)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadFrom(bytes.NewReader(sample.Bytes()), "bench")
	}
}
