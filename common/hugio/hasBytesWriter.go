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

	i    int
	done bool
	buff []byte
}

type HasBytesPattern struct {
	Match   bool
	Pattern []byte
}

func (h *HasBytesWriter) patternLen() int {
	l := 0
	for _, p := range h.Patterns {
		l += len(p.Pattern)
	}
	return l
}

func (h *HasBytesWriter) Write(p []byte) (n int, err error) {
	if h.done {
		return len(p), nil
	}

	if len(h.buff) == 0 {
		h.buff = make([]byte, h.patternLen()*2)
	}

	for i := range p {
		h.buff[h.i] = p[i]
		h.i++
		if h.i == len(h.buff) {
			// Shift left.
			copy(h.buff, h.buff[len(h.buff)/2:])
			h.i = len(h.buff) / 2
		}

		for _, pp := range h.Patterns {
			if bytes.Contains(h.buff, pp.Pattern) {
				pp.Match = true
				done := true
				for _, ppp := range h.Patterns {
					if !ppp.Match {
						done = false
						break
					}
				}
				if done {
					h.done = true
				}
				return len(p), nil
			}
		}

	}

	return len(p), nil
}
