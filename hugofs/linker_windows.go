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

package hugofs

import (
	"errors"
	"os"
	"syscall"
)

func deviceID(fi os.FileInfo) uint64 {
	return 0
}

// hasMultipleLinks reports whether the file has more than one hard link.
func hasMultipleLinks(filename string, fi os.FileInfo) bool {
	p, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return true
	}
	h, err := syscall.CreateFile(p, 0, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE, nil, syscall.OPEN_EXISTING, syscall.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		return true
	}
	defer syscall.CloseHandle(h)
	var d syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(h, &d); err != nil {
		return true
	}
	return d.NumberOfLinks > 1
}

func isLinkUnsupported(err error) bool {
	const (
		errorInvalidFunction = syscall.Errno(1)
		errorNotSameDevice   = syscall.Errno(17)
	)
	return errors.Is(err, errorNotSameDevice) || errors.Is(err, errorInvalidFunction)
}

func isLinkNotPossible(err error) bool {
	const errorTooManyLinks = syscall.Errno(1142)
	return errors.Is(err, errorTooManyLinks)
}
