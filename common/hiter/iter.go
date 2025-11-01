package hiter

// Common iterator functions.
// Some of these are are based on this discsussion: https://github.com/golang/go/issues/61898

import "iter"

// Concat returns an iterator over the concatenation of the sequences.
// Any nil sequences are ignored.
func Concat[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			if seq == nil {
				continue
			}
			for e := range seq {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// Concat2 returns an iterator over the concatenation of the sequences.
// Any nil sequences are ignored.
func Concat2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, seq := range seqs {
			if seq == nil {
				continue
			}
			for k, v := range seq {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
