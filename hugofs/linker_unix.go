// Copyright 2026 The Hugo Authors. All rights reserved.
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

//go:build unix

package hugofs

import (
	"errors"
	"os"
	"syscall"
)

func deviceID(fi os.FileInfo) uint64 {
	if st, ok := fi.Sys().(*syscall.Stat_t); ok {
		return uint64(st.Dev)
	}
	return 0
}

// isLinkUnsupported reports whether err means that the source device does not support hard links to the destination.
func isLinkUnsupported(err error) bool {
	return errors.Is(err, syscall.EXDEV) || errors.Is(err, syscall.ENOTSUP)
}

// isLinkNotPossible reports whether err means that this file cannot be linked, e.g. the link limit is reached
// or fs.protected_hardlinks denies it.
func isLinkNotPossible(err error) bool {
	return errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EMLINK)
}
