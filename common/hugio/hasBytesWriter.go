// Copyright 2024 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugio

import (
	"bytes"
)

// HasBytesWriter is a writer will match against a slice of patterns.
type HasBytesWriter struct {
	Patterns []*HasBytesPattern

	done bool
	// The tail of the bytes written so far, retained so we can detect a
	// pattern that straddles the boundary between two Write calls.
	buff []byte
}

type HasBytesPattern struct {
	Match   bool
	Pattern []byte
}

// maxPatternLen returns the length of the longest pattern.
func (h *HasBytesWriter) maxPatternLen() int {
	l := 0
	for _, p := range h.Patterns {
		if len(p.Pattern) > l {
			l = len(p.Pattern)
		}
	}
	return l
}

func (h *HasBytesWriter) Write(p []byte) (n int, err error) {
	if h.done {
		return len(p), nil
	}

	keep := h.maxPatternLen() - 1

	// Join the tail retained from previous Writes with the head of this chunk
	// so a pattern straddling the boundary is still detected. Only the
	// boundary window is copied; the chunk itself is scanned in place below.
	var boundary []byte
	if keep > 0 && len(h.buff) > 0 {
		head := p
		if len(head) > keep {
			head = head[:keep]
		}
		boundary = make([]byte, 0, len(h.buff)+len(head))
		boundary = append(boundary, h.buff...)
		boundary = append(boundary, head...)
	}

	// Scan each not-yet-matched pattern once per Write instead of once per byte.
	done := true
	for _, pp := range h.Patterns {
		if pp.Match {
			continue
		}
		if bytes.Contains(p, pp.Pattern) || bytes.Contains(boundary, pp.Pattern) {
			pp.Match = true
			continue
		}
		done = false
	}

	if done {
		// All patterns found; no need to look at any more data.
		h.done = true
		h.buff = nil
		return len(p), nil
	}

	// Retain the last keep bytes of (previous tail + this chunk) to detect a
	// pattern straddling into the next Write.
	switch {
	case keep <= 0:
		h.buff = h.buff[:0]
	case len(p) >= keep:
		h.buff = append(h.buff[:0], p[len(p)-keep:]...)
	default:
		// Chunk shorter than keep: slide the window over the retained tail.
		if total := len(h.buff) + len(p); total > keep {
			h.buff = h.buff[total-keep:]
		}
		h.buff = append(h.buff, p...)
	}

	return len(p), nil
}
