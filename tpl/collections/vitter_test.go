// This is just a temporary fork of https://github.com/josharian/vitter (ISC License, https://github.com/josharian/vitter/blob/main/LICENSE)
//
// This file will be removed once https://github.com/josharian/vitter/issues/1 is resolved.

package collections

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"testing"
	"time"
)

var goldenTests = []struct {
	seed   int64
	k, max int
	want   []int
}{
	{2, 10, 100, []int{6, 20, 34, 45, 58, 59, 64, 69, 70, 72}},
	{3, 10, 100, []int{8, 11, 22, 26, 30, 40, 74, 76, 93, 95}},
	{4, 5, 1000, []int{183, 283, 443, 501, 531}},
	{5, 15, 100000, []int{12984, 17778, 20370, 23830, 27120, 33258, 45718, 50064, 57096, 58580, 80960, 84396, 84594, 95561, 97687}},
}

func TestGolden(t *testing.T) {
	for _, test := range goldenTests {
		prng := rand.New(rand.NewPCG(uint64(test.seed), 0))
		var got []int
		testD(prng, t, test.k, test.max, func(n int) {
			got = append(got, n)
		})
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("golden(%d, %d, %d) = %#v want %#v", test.seed, test.k, test.max, got, test.want)
		}
	}
}

func TestInspectCounts(t *testing.T) {
	prng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	const max = 100
	const k = 10
	const iters = 10000
	counts := make([]int, max)
	for i := 0; i < iters; i++ {
		testD(prng, t, k, max, func(n int) {
			counts[n]++
		})
	}
	for i := range counts {
		counts[i] -= (iters * k / max)
	}
	t.Log(counts)
}

func testD(prng *rand.Rand, tb testing.TB, want, max int, choose func(n int)) {
	prev := -1
	got := want
	_d(prng, want, max, func(x int) {
		if x <= prev {
			tb.Fatalf("backwards: %d then %d", prev, x)
		}
		if x < 0 || x >= max {
			tb.Fatalf("bad selection: %d", x)
		}
		prev = x
		got--
		if got < 0 {
			tb.Fatal("choose called too many times")
		}
		choose(x)
	})
	if got != 0 {
		tb.Fatal("choose not called enough times")
	}
}

func TestWantIsMax(t *testing.T) {
	// Ensure that when want == max, we get all indices.
	prng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	const n = 10000
	testD(prng, t, n, n, func(n int) {})
}

func BenchmarkD(b *testing.B) {
	prng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	// TODO: count rng calls?
	for _, want := range []int{1, 100, 10000} {
		b.Run(fmt.Sprintf("want=%d", want), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_d(prng, want, 1000000, func(int) {})
			}
		})
	}
}
