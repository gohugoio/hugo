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
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"

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

type ProgressState int

const (
	ProgressHidden ProgressState = iota
	ProgressNormal
	ProgressError
	ProgressIntermediate
	ProgressWarning

	// ProgressDone is hidden in OSC 9;4.
	ProgressDone
)

const maxStatusMsgLen = 2048

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// StripANSI removes ANSI escape codes from s.
func StripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

// ReportProgress writes OSC 9;4 and OSC 7501 (program status) sequences to w.
// msg is only used in OSC 7501.
func ReportProgress(w io.Writer, state ProgressState, progress float64, msg string) {
	if progress < 0 {
		progress = 0.0
	}
	if progress > 1 {
		progress = 1.0
	}

	pi := int(progress * 100)

	osc94State := state
	if state == ProgressDone {
		osc94State = ProgressHidden
	}

	var status string
	switch state {
	case ProgressHidden:
		status = "state=idle"
	case ProgressNormal, ProgressWarning:
		status = fmt.Sprintf("state=working:progress=%d", pi)
	case ProgressIntermediate:
		status = "state=working"
	case ProgressError:
		status = "state=error"
	case ProgressDone:
		status = "state=done"
	}
	status += ":app=hugo"
	if msg = statusMsg(msg); msg != "" {
		status += ":msg=" + base64.StdEncoding.EncodeToString([]byte(msg))
	}

	fmt.Fprintf(w, "\033]9;4;%d;%d\007\033]7501;%s\033\\", osc94State, pi, status)
}

// statusMsg returns the first line of s without control characters, truncated to fit OSC 7501.
func statusMsg(s string) string {
	s = StripANSI(s)
	s, _, _ = strings.Cut(s, "\n")
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
	s = strings.TrimSpace(s)
	if len(s) > maxStatusMsgLen {
		s = strings.ToValidUTF8(s[:maxStatusMsgLen], "")
	}
	return s
}
