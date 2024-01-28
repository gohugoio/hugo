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

// Package terminal contains helper for the terminal, such as coloring output.
package terminal

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	isatty "github.com/mattn/go-isatty"
)

const (
	errorColor   = "\033[1;31m%s\033[0m"
	warningColor = "\033[0;33m%s\033[0m"
	noticeColor  = "\033[1;36m%s\033[0m"
)

// PrintANSIColors returns false if NO_COLOR env variable is set,
// else  IsTerminal(f).
func PrintANSIColors(f *os.File) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return IsTerminal(f)
}

// IsTerminal return true if the file descriptor is terminal and the TERM
// environment variable isn't a dumb one.
func IsTerminal(f *os.File) bool {
	if runtime.GOOS == "windows" {
		return false
	}

	fd := f.Fd()
	return os.Getenv("TERM") != "dumb" && (isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd))
}

// Notice colorizes the string in a noticeable color.
func Notice(s string) string {
	return colorize(s, noticeColor)
}

// Error colorizes the string in a colour that grabs attention.
func Error(s string) string {
	return colorize(s, errorColor)
}

// Warning colorizes the string in a colour that warns.
func Warning(s string) string {
	return colorize(s, warningColor)
}

// colorize s in color.
func colorize(s, color string) string {
	s = fmt.Sprintf(color, doublePercent(s))
	return singlePercent(s)
}

func doublePercent(str string) string {
	return strings.Replace(str, "%", "%%", -1)
}

func singlePercent(str string) string {
	return strings.Replace(str, "%%", "%", -1)
}
